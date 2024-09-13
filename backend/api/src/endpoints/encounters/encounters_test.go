package encounters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockEncountersDB struct {
	Encounters []Encounter
	Err        error
}

func (m *MockEncountersDB) FindAll() ([]Encounter, error) {
	if m.Encounters == nil {
		return nil, errors.New("no encounters found")
	}
	return m.Encounters, nil
}

func (m *MockEncountersDB) FindByID(id int) (*Encounter, error) {
	for _, encounter := range m.Encounters {
		if encounter.Id == id {
			return &encounter, nil
		}
	}
	return nil, errors.New("encounter not found")
}

func (m *MockEncountersDB) FindByCustomerID(id int) ([]Encounter, error) {
	var e []Encounter
	for _, encounter := range m.Encounters {
		if encounter.Customer_Id == id {
			e = append(e, encounter)
		}
	}
	return e, nil
}

func (m *MockEncountersDB) Add(encounter *AddEncounter) (*Encounter, error) {
	newEncounter := Encounter{
		Id:          len(m.Encounters) + 1,
		Date:        encounter.Date,
		Rating:      encounter.Rating,
		Comment:     encounter.Comment,
		Source:      encounter.Source,
		Customer_Id: encounter.Customer_Id,
	}
	m.Encounters = append(m.Encounters, newEncounter)
	if m.Err != nil {
		return nil, m.Err
	}
	return &newEncounter, nil
}

func (m *MockEncountersDB) Delete(id int) error {
	for i, encounter := range m.Encounters {
		if encounter.Id == id {
			m.Encounters = append(m.Encounters[:i], m.Encounters[i+1:]...)
			return nil
		}
	}
	return errors.New("encounter not found")
}

func (m *MockEncountersDB) Patch(id int, updates *UpdateEncounter) (*Encounter, error) {
	for i, encounter := range m.Encounters {
		if encounter.Id == id {
			if updates.Comment != nil {
				m.Encounters[i].Comment = *updates.Comment
			}
			if updates.Date != nil {
				m.Encounters[i].Date = *updates.Date
			}
			return &m.Encounters[i], nil
		}
	}
	return nil, errors.New("encounter not found")
}

func setupTestModel(mockDB *MockEncountersDB) *EncounterModel {
	return &EncounterModel{
		Encounters: mockDB,
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

func testAddEncounter(t *testing.T) {
	model := setupTestModel(&MockEncountersDB{})
	t.Run("Error JSON Decode Add Encounter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/encounters", bytes.NewBuffer([]byte(`{invalid json}`)))
		rr := httptest.NewRecorder()
		model.AddEncounter(rr, req)
		checkResponseCode(t, rr, http.StatusInternalServerError)
	})

	t.Run("Error Add Method", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPost, "/api/encounters", &AddEncounter{
			Date:        "25-04-2023",
			Rating:      5,
			Comment:     "Great encounter",
			Source:      "Referral",
			Customer_Id: 1,
		})
		rr := httptest.NewRecorder()
		model.AddEncounter(rr, req)
		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Add Encounter", func(t *testing.T) {
		req := createRequest(t, http.MethodPost, "/api/encounters", &AddEncounter{
			Date:        "25-04-2023",
			Rating:      5,
			Comment:     "Great encounter",
			Source:      "Referral",
			Customer_Id: 1,
		})
		rr := httptest.NewRecorder()
		model.AddEncounter(rr, req)
		checkResponseCode(t, rr, http.StatusOK)

		var addedEncounter Encounter
		decodeResponseBody(t, rr, &addedEncounter)

		if addedEncounter.Comment != "Great encounter" {
			t.Errorf("Expected comment 'Great encounter', got %s", addedEncounter.Comment)
		}
	})
}

func testGetAllEncounters(t *testing.T) {
	t.Run("Get All Encounters", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		model.Encounters.Add(&AddEncounter{Date: "25-04-2023", Rating: 5, Comment: "Encounter 1", Source: "Referral", Customer_Id: 1})
		model.Encounters.Add(&AddEncounter{Date: "26-04-2023", Rating: 4, Comment: "Encounter 2", Source: "Website", Customer_Id: 2})

		req := createRequest(t, http.MethodGet, "/api/encounters", nil)
		rr := httptest.NewRecorder()
		model.GetAllEncounters(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var encounters []Encounter
		decodeResponseBody(t, rr, &encounters)

		if len(encounters) != 2 {
			t.Errorf("Expected 2 encounters, got %d", len(encounters))
		}
	})

	t.Run("Error Fetching Encounters", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/encounters", nil)
		rr := httptest.NewRecorder()
		model.GetAllEncounters(rr, req)

		checkResponseCode(t, rr, http.StatusInternalServerError)
	})
}

func testGetEncounterById(t *testing.T) {
	t.Run("Get Encounter by ID", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		model.Encounters.Add(&AddEncounter{
			Date:        "25-04-2023",
			Rating:      5,
			Comment:     "Encounter 3",
			Source:      "Referral",
			Customer_Id: 3,
		})

		req := createRequest(t, http.MethodGet, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var encounter Encounter
		decodeResponseBody(t, rr, &encounter)

		if encounter.Id != 1 {
			t.Errorf("Expected encounter ID 1, got %d", encounter.Id)
		}
	})

	t.Run("Encounter Not Found", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		req := createRequest(t, http.MethodGet, "/api/encounters/99", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "99"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

	t.Run("Error Fetching Encounter", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})
}

func testGetEncounterByCustomerId(t *testing.T) {
	t.Run("Get Encounter by customer ID", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		model.Encounters.Add(&AddEncounter{
			Date:        "20-04-2023",
			Rating:      5,
			Comment:     "Encounter 10",
			Source:      "Referral",
			Customer_Id: 10,
		})

		req := createRequest(t, http.MethodGet, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var encounter Encounter
		decodeResponseBody(t, rr, &encounter)

		if encounter.Id != 1 {
			t.Errorf("Expected encounter ID 1, got %d", encounter.Id)
		}
	})

	t.Run("Encounter Not Found", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		req := createRequest(t, http.MethodGet, "/api/encounters/99", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "99"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

	t.Run("Error Fetching Encounter", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()
		model.GetEncounterById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})
}

func testDeleteEncounter(t *testing.T) {
	t.Run("Delete Encounter", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		model.Encounters.Add(&AddEncounter{
			Date:        "26-04-2023",
			Rating:      5,
			Comment:     "Encounter 4",
			Source:      "Referral",
			Customer_Id: 4,
		})

		req := createRequest(t, http.MethodDelete, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		if len(model.Encounters.(*MockEncountersDB).Encounters) != 0 {
			t.Errorf("Expected no encounters left, found %d", len(model.Encounters.(*MockEncountersDB).Encounters))
		}
	})

	t.Run("Encounter Not Found for Deletion", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		req := createRequest(t, http.MethodDelete, "/api/encounters/99", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "99"})
		rr := httptest.NewRecorder()

		model.DeleteEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Deleting Encounter", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodDelete, "/api/encounters/1", nil)
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func testPatchEncounters(t *testing.T) {
	newComment := "Updated Encounter"
	t.Run("Patch Event", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		model.Encounters.Add(&AddEncounter{
			Date:        "27-04-2023",
			Rating:      5,
			Comment:     "Encounter 5",
			Source:      "Referral",
			Customer_Id: 5,
		})

		req := createRequest(t, http.MethodPatch, "/api/encounters/1", &UpdateEncounter{Comment: &newComment})
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var UpdateEncounter Encounter
		decodeResponseBody(t, rr, &UpdateEncounter)

		if UpdateEncounter.Comment != "Updated Encounter" {
			t.Errorf("Expected updated encounter comment, got %s", UpdateEncounter.Comment)
		}
	})

	t.Run("Encounter Not Found for Patch", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{})
		req := createRequest(t, http.MethodPatch, "/api/encounters/99", &UpdateEncounter{Comment: &newComment})
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "99"})
		rr := httptest.NewRecorder()

		model.PatchEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Patching Encounter", func(t *testing.T) {
		model := setupTestModel(&MockEncountersDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPatch, "/api/encounters/1", &UpdateEncounter{Comment: &newComment})
		req = mux.SetURLVars(req, map[string]string{"encounter_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchEncounter(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func TestEncountersEndpoints(t *testing.T) {
	testGetAllEncounters(t)
	testAddEncounter(t)
	testGetEncounterById(t)
	testGetEncounterByCustomerId(t)
	testPatchEncounters(t)
	testDeleteEncounter(t)
}
