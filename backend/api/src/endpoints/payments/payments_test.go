package payments

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockPaymentsDB struct {
	Payments []Payment
	Err      error
}

func (m *MockPaymentsDB) FindAll() ([]Payment, error) {
	if m.Payments == nil {
		return nil, errors.New("no payments found")
	}
	return m.Payments, nil
}

func (m *MockPaymentsDB) FindByID(id int) (*Payment, error) {
	for _, payment := range m.Payments {
		if payment.Id == id {
			return &payment, nil
		}
	}
	return nil, errors.New("payment not found")
}

func (m *MockPaymentsDB) FindByCustomerID(id int) ([]Payment, error) {
	var p []Payment
	for _, payment := range m.Payments {
		if payment.CustomerId == id {
			p = append(p, payment)
		}
	}
	return p, nil
}

func (m *MockPaymentsDB) Add(payment *AddPayment) (*Payment, error) {
	newPayment := Payment{
		Id:                 len(m.Payments) + 1,
		Soul_Connection_Id: payment.Soul_Connection_Id,
		Date:               payment.Date,
		PaymentMethod:      payment.PaymentMethod,
		Amount:             payment.Amount,
		Comment:            payment.Comment,
		CustomerId:         payment.CustomerId,
	}
	m.Payments = append(m.Payments, newPayment)
	if m.Err != nil {
		return nil, m.Err
	}
	return &newPayment, nil
}

func (m *MockPaymentsDB) Delete(id int) error {
	for i, payment := range m.Payments {
		if payment.Id == id {
			m.Payments = append(m.Payments[:i], m.Payments[i+1:]...)
			return nil
		}
	}
	return errors.New("payment not found")
}

func (m *MockPaymentsDB) Patch(id int, updates *UpdatePayment) (*Payment, error) {
	for i, payment := range m.Payments {
		if payment.Id == id {
			if updates.Payment_Method != nil {
				m.Payments[i].PaymentMethod = *updates.Payment_Method
			}
			if updates.Comment != nil {
				m.Payments[i].Comment = *updates.Comment
			}
			return &m.Payments[i], nil
		}
	}
	return nil, errors.New("payment not found")
}

func setupTestModel(mockDB *MockPaymentsDB) *PaymentModel {
	return &PaymentModel{
		Payments: mockDB,
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

func testAddPayment(t *testing.T) {
	model := setupTestModel(&MockPaymentsDB{})
	t.Run("Error JSON Decode Add Payment", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/payments", bytes.NewBuffer([]byte(`{invalid json}`)))
		rr := httptest.NewRecorder()
		model.AddPayment(rr, req)
		checkResponseCode(t, rr, http.StatusInternalServerError)
	})

	t.Run("Error Add Method", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPost, "/api/payments", &AddPayment{
			Date:          "25-04-2023",
			PaymentMethod: "cash",
			Amount:        1.0,
			Comment:       "This should cause an error",
			CustomerId:    1,
		})
		rr := httptest.NewRecorder()
		model.AddPayment(rr, req)
		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Add Payment", func(t *testing.T) {
		req := createRequest(t, http.MethodPost, "/api/payments", &AddPayment{
			Date:          "25-04-2023",
			PaymentMethod: "cash",
			Amount:        1.0,
			Comment:       "This is a new payment",
			CustomerId:    1,
		})
		rr := httptest.NewRecorder()
		model.AddPayment(rr, req)
		checkResponseCode(t, rr, http.StatusOK)

		var addedPayment Payment
		decodeResponseBody(t, rr, &addedPayment)

		if addedPayment.Comment != "This is a new payment" {
			t.Errorf("Expected comment 'This is a new payment', got %s", addedPayment.Comment)
		}
	})
}

func testGetAllPayments(t *testing.T) {
	t.Run("Get All Payments", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		model.Payments.Add(&AddPayment{Date: "25-04-2023", PaymentMethod: "cash", Amount: 1.0, Comment: "Payment 1", CustomerId: 1})
		model.Payments.Add(&AddPayment{Date: "26-04-2023", PaymentMethod: "card", Amount: 2.0, Comment: "Payment 2", CustomerId: 1})

		req := createRequest(t, http.MethodGet, "/api/payments", nil)
		rr := httptest.NewRecorder()
		model.GetAllPayments(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var payments []Payment
		decodeResponseBody(t, rr, &payments)

		if len(payments) != 2 {
			t.Errorf("Expected 2 payments, got %d", len(payments))
		}
	})

	t.Run("Error Fetching Payments", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/payments", nil)
		rr := httptest.NewRecorder()
		model.GetAllPayments(rr, req)

		checkResponseCode(t, rr, http.StatusInternalServerError)
	})
}

func testGetPaymentById(t *testing.T) {
	t.Run("Get Payment by ID", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		model.Payments.Add(&AddPayment{Date: "25-04-2023", PaymentMethod: "cash", Amount: 1.0, Comment: "Payment 1", CustomerId: 1})

		req := createRequest(t, http.MethodGet, "/api/payments/1", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()
		model.GetPaymentsById(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var payment Payment
		decodeResponseBody(t, rr, &payment)

		if payment.Id != 1 {
			t.Errorf("Expected payment ID 1, got %d", payment.Id)
		}
	})

	t.Run("Payment Not Found", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		req := createRequest(t, http.MethodGet, "/api/payments/99", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "99"})
		rr := httptest.NewRecorder()
		model.GetPaymentsById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

	t.Run("Error Fetching Payment", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodGet, "/api/payments/1", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()
		model.GetPaymentsById(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})
}

func testGetPaymentByCustomerId(t *testing.T) {

	t.Run("Payment Not Found", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		req := createRequest(t, http.MethodGet, "/api/payments/99", nil)
		req = mux.SetURLVars(req, map[string]string{"customer_id": "99"})
		rr := httptest.NewRecorder()

		model.GetPaymentsByCustomerId(rr, req)

		checkResponseCode(t, rr, http.StatusNotFound)
	})

}

func testDeletePayment(t *testing.T) {
	t.Run("Delete Payment", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		model.Payments.Add(&AddPayment{Date: "25-04-2023", PaymentMethod: "cash", Amount: 1.0, Comment: "Payment to delete", CustomerId: 1})

		req := createRequest(t, http.MethodDelete, "/api/payments/1", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()

		model.DeletePayment(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		if len(model.Payments.(*MockPaymentsDB).Payments) != 0 {
			t.Errorf("Expected no payments left, found %d", len(model.Payments.(*MockPaymentsDB).Payments))
		}
	})

	t.Run("Payment Not Found for Deletion", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		req := createRequest(t, http.MethodDelete, "/api/payments/99", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "99"})
		rr := httptest.NewRecorder()

		model.DeletePayment(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Deleting Payment", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodDelete, "/api/payments/1", nil)
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()

		model.DeletePayment(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func testPatchPayments(t *testing.T) {
	newMethod := "credit"
	t.Run("Patch Payment", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		model.Payments.Add(&AddPayment{Date: "25-04-2023", PaymentMethod: "cash", Amount: 1.0, Comment: "Original Payment", CustomerId: 1})

		req := createRequest(t, http.MethodPatch, "/api/payments/1", &UpdatePayment{Payment_Method: &newMethod})
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchPayment(rr, req)

		checkResponseCode(t, rr, http.StatusOK)

		var updatedPayment Payment
		decodeResponseBody(t, rr, &updatedPayment)

		if updatedPayment.PaymentMethod != newMethod {
			t.Errorf("Expected updated payment method, got %s", updatedPayment.PaymentMethod)
		}
	})

	t.Run("Payment Not Found for Patch", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{})
		req := createRequest(t, http.MethodPatch, "/api/payments/99", &UpdatePayment{Payment_Method: &newMethod})
		req = mux.SetURLVars(req, map[string]string{"payment_id": "99"})
		rr := httptest.NewRecorder()

		model.PatchPayment(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})

	t.Run("Error Patching Payment", func(t *testing.T) {
		model := setupTestModel(&MockPaymentsDB{Err: errors.New("database error")})
		req := createRequest(t, http.MethodPatch, "/api/payments/1", &UpdatePayment{Payment_Method: &newMethod})
		req = mux.SetURLVars(req, map[string]string{"payment_id": "1"})
		rr := httptest.NewRecorder()

		model.PatchPayment(rr, req)

		checkResponseCode(t, rr, http.StatusBadRequest)
	})
}

func TestPaymentsEndpoints(t *testing.T) {
	testGetAllPayments(t)
	testAddPayment(t)
	testGetPaymentById(t)
	testGetPaymentByCustomerId(t)
	testPatchPayments(t)
	testDeletePayment(t)
}
