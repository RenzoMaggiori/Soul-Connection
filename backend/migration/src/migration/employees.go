package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"soul-connection.com/api/src/endpoints/employees"
	"soul-connection.com/api/src/lib"
)

func migrateEmployees(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials) error {
	time.Sleep(2 * time.Second)
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/employees", lib.ApiBaseUri),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("could not retrieve users from api")
	}
	var employeesResponse []struct {
		Id      int
		Email   string
		Name    string
		Surname string
	}
	err = json.NewDecoder(resp.Body).Decode(&employeesResponse)
	if err != nil {
		return err
	}

	bucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("employeesBucket"))
	if err != nil {
		return err
	}
	employeesDb := employees.EmployeesDB{DB: database, Bucket: bucket}
	lib.ServerLog("PROGRESS", fmt.Sprintf("Migrating employees...:START:%d", len(employeesResponse)))
	for _, e := range employeesResponse {
		lib.ServerLog("PROGRESS", "Migrating employees...:INCREMENT")
		migrateEmployee(employeesDb, credentials, e.Id)
	}
	lib.ServerLog("PROGRESS", "Migrating employees...:COMPLETE")
	return nil
}

func migrateEmployee(db employees.EmployeesDB, credentials *ApiCredentials, id int) error {
	var employeeResponse struct {
		Id    int
		Email string
		// Password  string
		Name       string
		Surname    string
		Birth_Date string
		Gender     string
		Work       string
	}
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/employees/%d", lib.ApiBaseUri, id),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve user with id %d from api", id)
	}
	err = json.NewDecoder(resp.Body).Decode(&employeeResponse)
	if err != nil {
		return err
	}

	employee, err := db.Add(&employees.AddEmployee{
		Soul_Connection_Id: &employeeResponse.Id,
		Email:              employeeResponse.Email,
		Name:               employeeResponse.Name,
		Surname:            employeeResponse.Surname,
		Birth_Date:         employeeResponse.Birth_Date,
		Gender:             employeeResponse.Gender,
		Work:               employeeResponse.Work,
	})
	if err != nil {
		return err
	}

	err = migrateEmployeeImage(&db, credentials, &Ids{old: id, new: employee.Id})
	if err != nil {
		return err
	}
	return nil
}

func migrateEmployeeImage(db *employees.EmployeesDB, credentials *ApiCredentials, ids *Ids) error {
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/employees/%d/image", lib.ApiBaseUri, ids.old),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve user image with id %d from api", ids.old)
	}

	_, err = db.UploadFile(ids.new, resp.Body, fmt.Sprintf("employee_%d", ids.new))
	if err != nil {
		return err
	}
	return nil
}
