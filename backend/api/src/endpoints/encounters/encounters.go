package encounters

import (
	"encoding/json"
	"net/http"
	"time"

	"soul-connection.com/api/src/lib"
)

type Encounter struct {
	Id          int
	Date        string
	Rating      int
	Comment     string
	Source      string
	CreatedAt   time.Time
	Customer_Id int
}

type EncounterModel struct {
	Encounters interface {
		FindAll() ([]Encounter, error)
		FindByID(int) (*Encounter, error)
		FindByCustomerID(int) ([]Encounter, error)
		Add(*AddEncounter) (*Encounter, error)
		Delete(int) error
		Patch(int, *UpdateEncounter) (*Encounter, error)
	}
}

func (model *EncounterModel) GetAllEncounters(res http.ResponseWriter, _ *http.Request) {
	encounters, err := model.Encounters.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(encounters); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EncounterModel) GetEncounterById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "encounter_id")
	if err != nil {
		http.Error(res, "Invalid encounter ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	encounter, err := model.Encounters.FindByID(id)
	if err != nil {
		http.Error(res, "Encounter not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*encounter); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EncounterModel) GetEncounterByCustomerId(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	encounters, err := model.Encounters.FindByCustomerID(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(encounters); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EncounterModel) AddEncounter(res http.ResponseWriter, req *http.Request) {
	var nc AddEncounter
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	encounter, err := model.Encounters.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add encounter", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*encounter); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EncounterModel) DeleteEncounter(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "encounter_id")
	if err != nil {
		http.Error(res, "Invalid encounter ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Encounters.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete encounter", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EncounterModel) PatchEncounter(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "encounter_id")
	if err != nil {
		http.Error(res, "Invalid encounter ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateEncounter
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedEncounter, err := model.Encounters.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update encounter", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedEncounter); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}
