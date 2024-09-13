package customers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"soul-connection.com/api/src/lib"
)

type Customer struct {
	Id                 int
	Soul_Connection_Id *int
	Email              string
	Name               string
	Surname            string
	Birth_Date         string
	Gender             string
	Description        string
	Astrological_Sign  string
	Phone_Number       string
	Address            string
	Image_Id           *string
	CreatedAt          time.Time
	Employee_Id        *int
}

type CustomersModel struct {
	Customers interface {
		FindAll() ([]Customer, error)
		FindByID(int) (*Customer, error)
		FindByEmployeeID(int) ([]Customer, error)
		FindByOldID(int) (*Customer, error)
		Add(*AddCustomer) (*Customer, error)
		Delete(int) error
		Patch(int, *UpdateCustomer) (*Customer, error)
		UploadFile(int, io.Reader, string) (*primitive.ObjectID, error)
		GetFile(primitive.ObjectID) ([]byte, error)
	}
}

func (model *CustomersModel) GetAllCustomers(res http.ResponseWriter, _ *http.Request) {
	customers, err := model.Customers.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(customers); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) GetCustomerById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	customer, err := model.Customers.FindByID(id)
	if err != nil {
		http.Error(res, "Customer not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*customer); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) GetCustomerByEmployeeId(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "employee_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	customer, err := model.Customers.FindByEmployeeID(id)
	if err != nil {
		http.Error(res, "Employee not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(customer); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) AddCustomer(res http.ResponseWriter, req *http.Request) {
	var nc AddCustomer
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	customer, err := model.Customers.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add customer", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*customer); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) DeleteCustomer(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Customers.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete customer", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) PatchCustomer(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateCustomer
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedCustomer, err := model.Customers.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update customer", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedCustomer); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *CustomersModel) GetImage(res http.ResponseWriter, req *http.Request) {
	customerId, err := lib.GetIdFromRequest(req, "customer_id")
	if err != nil {
		http.Error(res, "Invalid customer ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	customer, err := model.Customers.FindByID(customerId)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	if customer.Image_Id == nil {
		http.Error(res, "Could not find image for customer", http.StatusNotFound)
		return
	}

	fileId, err := primitive.ObjectIDFromHex(*customer.Image_Id)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	fileContent, err := model.Customers.GetFile(fileId)
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
