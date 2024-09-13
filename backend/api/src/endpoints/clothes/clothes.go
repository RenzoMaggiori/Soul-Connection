package clothes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"soul-connection.com/api/src/lib"
)

type Clothe struct {
	Id                 int
	Soul_Connection_Id *int
	Type               string
	Image_Id           *string
	CreatedAt          time.Time
	CustomerId         int
}

type ClothesModel struct {
	Clothes interface {
		FindAll() ([]Clothe, error)
		FindByID(int) (*Clothe, error)
		FindByCustomerID(int) ([]Clothe, error)
		Add(*AddClothe) (*Clothe, error)
		Delete(int) error
		Patch(int, *UpdateClothe) (*Clothe, error)
		UploadFile(int, io.Reader, string) (*primitive.ObjectID, error)
		GetFile(primitive.ObjectID) ([]byte, error)
	}
}

func (model *ClothesModel) GetAllClothes(res http.ResponseWriter, _ *http.Request) {
	clothes, err := model.Clothes.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(clothes); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) GetClotheById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "clothe_id")
	if err != nil {
		http.Error(res, "Invalid clothe ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	clothe, err := model.Clothes.FindByID(id)
	if err != nil {
		http.Error(res, "Clothe not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*clothe); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) GetClotheByCustomerId(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	clothe, err := model.Clothes.FindByCustomerID(id)
	if err != nil {
		http.Error(res, "Clothes not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(clothe); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) AddClothe(res http.ResponseWriter, req *http.Request) {
	var nc AddClothe
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	clothe, err := model.Clothes.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add clothe", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*clothe); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) DeleteClothe(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "clothe_id")
	if err != nil {
		http.Error(res, "Invalid clothe ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Clothes.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete clothe", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) PatchClothes(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "clothe_id")
	if err != nil {
		http.Error(res, "Invalid clothe ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateClothe
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedClothe, err := model.Clothes.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update clothe", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedClothe); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *ClothesModel) GetImage(res http.ResponseWriter, req *http.Request) {
	clotheId, err := lib.GetIdFromRequest(req, "clothe_id")
	if err != nil {
		http.Error(res, "Invalid clothe ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	clothe, err := model.Clothes.FindByID(clotheId)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	if clothe.Image_Id == nil {
		http.Error(res, "Could not find image for clothe", http.StatusNotFound)
		return
	}

	fileId, err := primitive.ObjectIDFromHex(*clothe.Image_Id)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	fileContent, err := model.Clothes.GetFile(fileId)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	_, err = res.Write(fileContent)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	res.Header().Set("Content-Type", "application/octet-stream")
	res.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))
}
