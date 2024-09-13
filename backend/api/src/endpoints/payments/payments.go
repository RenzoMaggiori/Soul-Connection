package payments

import (
	"encoding/json"
	"net/http"
	"time"

	"soul-connection.com/api/src/lib"
)

type Payment struct {
	Id                 int
	Soul_Connection_Id *int
	Date               string
	PaymentMethod      string
	Amount             float64
	Comment            string
	CreatedAt          time.Time
	CustomerId         int
}

type PaymentModel struct {
	Payments interface {
		FindAll() ([]Payment, error)
		FindByID(int) (*Payment, error)
		FindByCustomerID(int) ([]Payment, error)
		Add(*AddPayment) (*Payment, error)
		Delete(int) error
		Patch(int, *UpdatePayment) (*Payment, error)
	}
}

func (model *PaymentModel) GetAllPayments(res http.ResponseWriter, _ *http.Request) {
	payment, err := model.Payments.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(payment); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *PaymentModel) GetPaymentsById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "payment_id")
	if err != nil {
		http.Error(res, "Invalid payment ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	payment, err := model.Payments.FindByID(id)
	if err != nil {
		http.Error(res, "Payment not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*payment); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *PaymentModel) GetPaymentsByCustomerId(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	payment, err := model.Payments.FindByCustomerID(id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(payment); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *PaymentModel) AddPayment(res http.ResponseWriter, req *http.Request) {
	var nc AddPayment
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	event, err := model.Payments.Add(&nc)
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

func (model *PaymentModel) DeletePayment(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "payment_id")
	if err != nil {
		http.Error(res, "Invalid payment ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Payments.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete payment", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *PaymentModel) PatchPayment(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "payment_id")
	if err != nil {
		http.Error(res, "Invalid payment ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdatePayment
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedPayment, err := model.Payments.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update payment", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedPayment); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}
