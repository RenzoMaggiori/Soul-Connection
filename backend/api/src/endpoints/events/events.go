package events

import (
	"encoding/json"
	"net/http"
	"time"

	"soul-connection.com/api/src/lib"
)

type Event struct {
	Id               int
	Name             string
	Date             string
	Max_Participants int
	Location_X       string
	Location_Y       string
	Type             string
	CreatedAt        time.Time
	Employee_Id      int
}

type EventModel struct {
	Events interface {
		FindAll() ([]Event, error)
		FindByID(int) (*Event, error)
		Add(*AddEvent) (*Event, error)
		Delete(int) error
		Patch(int, *UpdateEvent) (*Event, error)
	}
}

func (model *EventModel) GetAllEvents(res http.ResponseWriter, _ *http.Request) {
	events, err := model.Events.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(events); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EventModel) GetEventsById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "event_id")
	if err != nil {
		http.Error(res, "Invalid event ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	event, err := model.Events.FindByID(id)
	if err != nil {
		http.Error(res, "Event not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*event); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EventModel) AddEvent(res http.ResponseWriter, req *http.Request) {
	var nc AddEvent
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	event, err := model.Events.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add event", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*event); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EventModel) DeleteEvent(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "event_id")
	if err != nil {
		http.Error(res, "Invalid employee ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Events.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete event", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EventModel) PatchEvent(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "event_id")
	if err != nil {
		http.Error(res, "Invalid event ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateEvent
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedEvent, err := model.Events.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update event", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedEvent); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}
