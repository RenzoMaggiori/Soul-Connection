package employees

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"soul-connection.com/api/src/lib"
)

type Employee struct {
	Id                 int
	Soul_Connection_Id *int
	Email              string
	Password           string
	Name               string
	Surname            string
	Birth_Date         string
	Gender             string
	Work               string
	Image_Id           *string
	CreatedAt          time.Time
}

type EmployeesModel struct {
	Employees interface {
		FindAll() ([]Employee, error)
		FindByID(int) (*Employee, error)
		FindByOldID(int) (*Employee, error)
		Add(*AddEmployee) (*Employee, error)
		Delete(int) error
		Patch(int, *UpdateEmployee) (*Employee, error)
		UploadFile(int, io.Reader, string) (*primitive.ObjectID, error)
		GetFile(primitive.ObjectID) ([]byte, error)
	}
}

func (model *EmployeesModel) GetAllEmployees(res http.ResponseWriter, _ *http.Request) {
	employees, err := model.Employees.FindAll()

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(employees); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EmployeesModel) GetEmployeeById(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "employee_id")
	if err != nil {
		http.Error(res, "Invalid employee ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	employee, err := model.Employees.FindByID(id)
	if err != nil {
		http.Error(res, "Employee not found", http.StatusNotFound)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*employee); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EmployeesModel) AddEmployee(res http.ResponseWriter, req *http.Request) {
	var nc AddEmployee
	err := json.NewDecoder(req.Body).Decode(&nc)

	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
	employee, err := model.Employees.Add(&nc)
	if err != nil {
		http.Error(res, "Unable to add employee", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(*employee); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EmployeesModel) DeleteEmployee(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "employee_id")
	if err != nil {
		http.Error(res, "Invalid employee ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	err = model.Employees.Delete(id)
	if err != nil {
		http.Error(res, "Unable to delete employee", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{"status": "ok"}); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EmployeesModel) PatchEmployee(res http.ResponseWriter, req *http.Request) {
	id, err := lib.GetIdFromRequest(req, "employee_id")
	if err != nil {
		http.Error(res, "Invalid employee ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	var updates UpdateEmployee
	err = json.NewDecoder(req.Body).Decode(&updates)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	updatedEmployee, err := model.Employees.Patch(id, &updates)
	if err != nil {
		http.Error(res, "Unable to update employee", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(updatedEmployee); err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}
}

func (model *EmployeesModel) GetImage(res http.ResponseWriter, req *http.Request) {
	employeeId, err := lib.GetIdFromRequest(req, "employee_id")
	if err != nil {
		http.Error(res, "Invalid employee ID format", http.StatusBadRequest)
		lib.ServerLog("ERROR", err)
		return
	}

	employee, err := model.Employees.FindByID(employeeId)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	if employee.Image_Id == nil {
		http.Error(res, "Could not find image for customer", http.StatusNotFound)
		return
	}

	fileId, err := primitive.ObjectIDFromHex(*employee.Image_Id)
	if err != nil {
		http.Error(res, "Internal server error", http.StatusInternalServerError)
		lib.ServerLog("ERROR", err)
		return
	}

	fileContent, err := model.Employees.GetFile(fileId)
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
