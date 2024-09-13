package tips_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"soul-connection.com/api/src/endpoints/tips"
)

type MockTipsDB struct {
	Tips []tips.Tip
}

func (m *MockTipsDB) FindAll() ([]tips.Tip, error) {
	if m.Tips == nil {
		return nil, errors.New("no tips found")
	}
	return m.Tips, nil
}

func (m *MockTipsDB) FindByID(id int) (*tips.Tip, error) {
	for _, tip := range m.Tips {
		if tip.Id == id {
			return &tip, nil
		}
	}
	return nil, errors.New("tip not found")
}

func (m *MockTipsDB) Add(tip *tips.AddTip) (*tips.Tip, error) {
	newTip := tips.Tip{Id: len(m.Tips) + 1, Title: tip.Title, Tip: tip.Tip}
	m.Tips = append(m.Tips, newTip)
	return &newTip, nil
}

func (m *MockTipsDB) Delete(id int) error {
	for i, tip := range m.Tips {
		if tip.Id == id {
			m.Tips = append(m.Tips[:i], m.Tips[i+1:]...)
			return nil
		}
	}
	return errors.New("tip not found")
}

func (m *MockTipsDB) Patch(id int, updates *tips.UpdateTip) (*tips.Tip, error) {
	for i, tip := range m.Tips {
		if tip.Id != id {
			continue
		}
		uReflection := reflect.ValueOf(updates).Elem()
		tReflection := reflect.ValueOf(&m.Tips[i]).Elem()

		for j := range uReflection.NumField() {
			v := uReflection.Field(j)
			if v.IsNil() {
				continue
			}

			tipField := tReflection.FieldByName(uReflection.Type().Field(j).Name)
			if !tipField.IsValid() || !tipField.CanSet() {
				continue
			}
			tipField.Set(v.Elem())
		}

		return &m.Tips[i], nil
	}
	return nil, errors.New("tip not found")
}

func setupTestModel(mockDB *MockTipsDB) *tips.TipModel {
	return &tips.TipModel{
		Tips: mockDB,
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

func testAddTip(t *testing.T) {
	model := setupTestModel(&MockTipsDB{})
	t.Run("Error JSON Decode Add Tip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/tips", bytes.NewBuffer([]byte(`{invalid json}`)))
		rr := httptest.NewRecorder()
		model.AddTip(rr, req)
		checkResponseCode(t, rr, http.StatusInternalServerError)
	})

	t.Run("Add Tip", func(t *testing.T) {
		req := createRequest(t, http.MethodPost, "/api/tips", &tips.AddTip{Title: "New Tip", Tip: "This is a new tip"})
		rr := httptest.NewRecorder()
		model.AddTip(rr, req)
		checkResponseCode(t, rr, http.StatusOK)

		var addedTip tips.Tip
		decodeResponseBody(t, rr, &addedTip)

		if addedTip.Title != "New Tip" {
			t.Errorf("Expected title 'New Tip', got %s", addedTip.Title)
		}
	})
}

func testGetAllTips(t *testing.T) {
	t.Run("Get All Tips", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		model.Tips.Add(&tips.AddTip{Title: "Tip 1", Tip: "This is tip 1"})
		model.Tips.Add(&tips.AddTip{Title: "Tip 2", Tip: "This is tip 2"})

		req := createRequest(t, http.MethodGet, "/api/tips", nil)
		rr := httptest.NewRecorder()
		model.GetAllTips(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var tips []tips.Tip
		decodeResponseBody(t, rr, &tips)

		if len(tips) != 2 {
			t.Errorf("Expected 2 tips, got %d", len(tips))
		}
	})

	t.Run("Error Fetching Tips", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodGet, "/api/tips", nil)
		rr := httptest.NewRecorder()
		model.GetAllTips(rr, req)

		checkResponseCode(t, rr, http.StatusInternalServerError)
	})
}

func testGetTipById(t *testing.T) {
	t.Run("Get Tip by ID", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		model.Tips.Add(&tips.AddTip{Title: "Tip 1", Tip: "This is tip 1"})

		req := createRequest(t, http.MethodGet, "/api/tips/1", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()
		model.GetTipById(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var tip tips.Tip
		decodeResponseBody(t, rr, &tip)

		if tip.Id != 1 {
			t.Errorf("Expected tip ID 1, got %d", tip.Id)
		}
	})

	t.Run("Tip Not Found", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodGet, "/api/tips/99", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "99"})
		rr := httptest.NewRecorder()
		model.GetTipById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

	t.Run("Error Fetching Tip", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodGet, "/api/tips/1", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()
		model.GetTipById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})
}

func testDeleteTip(t *testing.T) {
	t.Run("Delete Tip", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		model.Tips.Add(&tips.AddTip{Title: "Tip to delete", Tip: "This tip will be deleted"})

		req := createRequest(t, http.MethodDelete, "/api/tips/1", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteTip(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		if len(model.Tips.(*MockTipsDB).Tips) != 0 {
			t.Errorf("Expected no tips left, found %d", len(model.Tips.(*MockTipsDB).Tips))
		}
	})

	t.Run("Tip Not Found for Deletion", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodDelete, "/api/tips/99", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "99"})
		rr := httptest.NewRecorder()

		model.DeleteTip(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Deleting Tip", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodDelete, "/api/tips/1", nil)
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()

		model.DeleteTip(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func testPatchTips(t *testing.T) {
	newTitle := "Updated Title"
	t.Run("Patch Tip", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		model.Tips.Add(&tips.AddTip{Title: "Original Title", Tip: "Original Tip"})
		newTitle := "Updated Title"

		req := createRequest(t, http.MethodPatch, "/api/tips/1", &tips.UpdateTip{Title: &newTitle})
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchTips(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var updatedTip tips.Tip
		decodeResponseBody(t, rr, &updatedTip)

		if updatedTip.Title != "Updated Title" {
			t.Errorf("Expected updated title, got %s", updatedTip.Title)
		}
	})

	t.Run("Tip Not Found for Patch", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodPatch, "/api/tips/99", &tips.UpdateTip{Title: &newTitle})
		req = mux.SetURLVars(req, map[string]string{"tip_id": "99"})
		rr := httptest.NewRecorder()

		model.PatchTips(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Patching Tip", func(t *testing.T) {
		model := setupTestModel(&MockTipsDB{})
		req := createRequest(t, http.MethodPatch, "/api/tips/1", &tips.UpdateTip{Title: &newTitle})
		req = mux.SetURLVars(req, map[string]string{"tip_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchTips(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func TestTipsEndpoints(t *testing.T) {
	testGetAllTips(t)
	testAddTip(t)
	testGetTipById(t)
	testPatchTips(t)
	testDeleteTip(t)
}
