package tips

import (
	"encoding/json"
	"net/http"
	"time"

	"soul-connection.com/api/src/lib"
)

type Tip struct {
	Id        int
	Title     string
	Tip       string
	CreatedAt time.Time
}

type TipModel struct {
	Tips interface {
		FindAll() ([]Tip, error)
		FindByID(int) (*Tip, error)
		Add(*AddTip) (*Tip, error)
		Delete(int) error
		Patch(int, *UpdateTip) (*Tip, error)
	}
}

func (model *TipModel) GetAllTips(res http.ResponseWriter, _ *http.Request) {
	tips, err := model.Tips.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(tips); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *TipModel) GetTipById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "tip_id")
	if err != nil {
		http.Error(res, "Invalid tip ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	tip, err := model.Tips.FindByID(id)
	if err != nil {
		http.Error(res, "Tips not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*tip); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *TipModel) AddTip(res http.ResponseWriter, req *http.Request) {
	var nc AddTip
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	tip, err := model.Tips.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add tip", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*tip); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *TipModel) DeleteTip(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "tip_id")
	if err != nil {
		http.Error(res, "Invalid tip ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Tips.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete tip", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *TipModel) PatchTips(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "tip_id")
	if err != nil {
		http.Error(res, "Invalid tip ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateTip
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedTip, err := model.Tips.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update tip", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedTip); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}
