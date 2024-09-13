package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockEventsDB struct {
	Events []Event
	Err    error
}

func (m *MockEventsDB) FindAll() ([]Event, error) {
	if m.Events == nil {
		return nil, errors.New("no events found")
	}
	return m.Events, nil
}

func (m *MockEventsDB) FindByID(id int) (*Event, error) {
	for _, event := range m.Events {
		if event.Id == id {
			return &event, nil
		}
	}
	return nil, errors.New("event not found")
}

func (m *MockEventsDB) Add(event *AddEvent) (*Event, error) {
	newEvent := Event{
		Id:               len(m.Events) + 1,
		Name:             event.Name,
		Date:             event.Date,
		Max_Participants: event.Max_Participants,
		Location_X:       event.Location_X,
		Location_Y:       event.Location_Y,
		Type:             event.Type,
		Employee_Id:      event.Employee_Id,
	}
	m.Events = append(m.Events, newEvent)
	if m.Err != nil {
		return nil, m.Err
	}
	return &newEvent, nil
}

func (m *MockEventsDB) Delete(id int) error {
	for i, event := range m.Events {
		if event.Id == id {
			m.Events = append(m.Events[:i], m.Events[i+1:]...)
			return nil
		}
	}
	return errors.New("event not found")
}

func (m *MockEventsDB) Patch(id int, updates *UpdateEvent) (*Event, error) {
	for i, event := range m.Events {
		if event.Id == id {
			if updates.Name != nil {
				m.Events[i].Name = *updates.Name
			}
			if updates.Date != nil {
				m.Events[i].Date = *updates.Date
			}
			return &m.Events[i], nil
		}
	}
	return nil, errors.New("event not found")
}

func setupTestModel(mockDB *MockEventsDB) *EventModel {
	return &EventModel{
		Events: mockDB,
	}
}

func createRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("Failed to encode request body: %v", err)
		}
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func checkResponseCode(t *testing.T, rr *httptest.ResponseRecorder, expectedCode int) {
	if rr.Code != expectedCode {
		t.Errorf("Expected status %d, got %d", expectedCode, rr.Code)
	}
}

func decodeResponseBody(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}
}

func testAddEvent(t *testing.T) {
	model := setupTestModel(&MockEventsDB{})
	t.Run("Error JSON Decode Add Event", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/events", bytes.NewBuffer([]byte(`{invalid json}`)))
		rr := httptest.NewRecorder()
		model.AddEvent(rr, req)
		checkResponseCode(t, rr, http.StatusInternalServerError)
	})

	t.Run("Error Add Method", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPost, "/api/events", &AddEvent{
			Name:             "Event 1",
			Date:             "25-04-2023",
			Max_Participants: 100,
			Location_X:       "123.45",
			Location_Y:       "678.90",
			Type:             "Conference",
			Employee_Id:      1,
		})
		rr := httptest.NewRecorder()
		model.AddEvent(rr, req)
		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Add Event", func(t *testing.T) {
		req := createRequest(t, http.MethodPost, "/api/events", &AddEvent{
			Name:             "New Event",
			Date:             "25-04-2023",
			Max_Participants: 100,
			Location_X:       "123.45",
			Location_Y:       "678.90",
			Type:             "Conference",
			Employee_Id:      1,
		})
		rr := httptest.NewRecorder()
		model.AddEvent(rr, req)
		checkResponseCode(t, rr, http.StatusOK)

		var addedEvent Event
		decodeResponseBody(t, rr, &addedEvent)

		if addedEvent.Name != "New Event" {
			t.Errorf("Expected name 'New Event', got %s", addedEvent.Name)
		}
	})
}

func testGetAllEvents(t *testing.T) {
	t.Run("Get All Events", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		model.Events.Add(&AddEvent{Name: "Event 1", Date: "25-04-2023", Max_Participants: 100, Location_X: "123.45", Location_Y: "678.90", Type: "Conference", Employee_Id: 1})
		model.Events.Add(&AddEvent{Name: "Event 2", Date: "26-04-2023", Max_Participants: 150, Location_X: "123.46", Location_Y: "678.91", Type: "Seminar", Employee_Id: 2})

		req := createRequest(t, http.MethodGet, "/api/events", nil)
		rr := httptest.NewRecorder()
		model.GetAllEvents(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var events []Event
		decodeResponseBody(t, rr, &events)

		if len(events) != 2 {
			t.Errorf("Expected 2 events, got %d", len(events))
		}
	})

	t.Run("Error Fetching Events", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/events", nil)
		rr := httptest.NewRecorder()
		model.GetAllEvents(rr, req)

		checkResponseCode(t, rr, http.StatusInternalServerError)
	})
}

func testGetEventById(t *testing.T) {
	t.Run("Get Event by ID", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		model.Events.Add(&AddEvent{Name: "Event 1", Date: "25-04-2023", Max_Participants: 100, Location_X: "123.45", Location_Y: "678.90", Type: "Conference", Employee_Id: 1})

		req := createRequest(t, http.MethodGet, "/api/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEventsById(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var event Event
		decodeResponseBody(t, rr, &event)

		if event.Id != 1 {
			t.Errorf("Expected event ID 1, got %d", event.Id)
		}
	})

	t.Run("Event Not Found", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		req := createRequest(t, http.MethodGet, "/api/events/99", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "99"})
		rr := httptest.NewRecorder()
		model.GetEventsById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

	t.Run("Error Fetching Event", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEventsById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})
}

func testDeleteEvent(t *testing.T) {
	t.Run("Delete Event", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		model.Events.Add(&AddEvent{Name: "Event to delete", Date: "25-04-2023", Max_Participants: 100, Location_X: "123.45", Location_Y: "678.90", Type: "Conference", Employee_Id: 1})

		req := createRequest(t, http.MethodDelete, "/api/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteEvent(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		if len(model.Events.(*MockEventsDB).Events) != 0 {
			t.Errorf("Expected no events left, found %d", len(model.Events.(*MockEventsDB).Events))
		}
	})

	t.Run("Event Not Found for Deletion", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		req := createRequest(t, http.MethodDelete, "/api/events/99", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "99"})
		rr := httptest.NewRecorder()

		model.DeleteEvent(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Deleting Event", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodDelete, "/api/events/1", nil)
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteEvent(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func testPatchEvents(t *testing.T) {
	newName := "Updated Event"
	t.Run("Patch Event", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		model.Events.Add(&AddEvent{Name: "Original Event", Date: "25-04-2023", Max_Participants: 100, Location_X: "123.45", Location_Y: "678.90", Type: "Conference", Employee_Id: 1})

		req := createRequest(t, http.MethodPatch, "/api/events/1", &UpdateEvent{Name: &newName})
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchEvent(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var updatedEvent Event
		decodeResponseBody(t, rr, &updatedEvent)

		if updatedEvent.Name != "Updated Event" {
			t.Errorf("Expected updated event name, got %s", updatedEvent.Name)
		}
	})

	t.Run("Event Not Found for Patch", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{})
		req := createRequest(t, http.MethodPatch, "/api/events/99", &UpdateEvent{Name: &newName})
		req = mux.SetURLVars(req, map[string]string{"event_id": "99"})
		rr := httptest.NewRecorder()

		model.PatchEvent(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Patching Event", func(t *testing.T) {
		model := setupTestModel(&MockEventsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPatch, "/api/events/1", &UpdateEvent{Name: &newName})
		req = mux.SetURLVars(req, map[string]string{"event_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchEvent(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func TestEventsEndpoints(t *testing.T) {
	testGetAllEvents(t)
	testAddEvent(t)
	testGetEventById(t)
	testPatchEvents(t)
	testDeleteEvent(t)
}
