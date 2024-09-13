package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"soul-connection.com/api/src/endpoints/employees"
	"soul-connection.com/api/src/endpoints/events"
	"soul-connection.com/api/src/lib"
)

func migrateEvents(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials) error {
	time.Sleep(2 * time.Second)
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/events", lib.ApiBaseUri),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("could not retrieve events from api")
	}
	var eventsResponse []struct {
		Id               int
		Name             string
		Date             string
		Duration         int
		Max_Participants int
	}
	err = json.NewDecoder(resp.Body).Decode(&eventsResponse)
	if err != nil {
		return err
	}

	lib.ServerLog("PROGRESS", fmt.Sprintf("Migrating events...:START:%d", len(eventsResponse)))
	for _, e := range eventsResponse {
		lib.ServerLog("PROGRESS", "Migrating events...:INCREMENT")
		migrateEvent(database, fileStorage, credentials, e.Id)
	}
	lib.ServerLog("PROGRESS", "Migrating events...:COMPLETE")
	return nil
}

func migrateEvent(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials, id int) error {
	var event events.AddEvent
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/events/%d", lib.ApiBaseUri, id),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve event with id %d from api", id)
	}

	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return err
	}

	eventsDb := events.EventsDB{DB: database}
	employeeDb := employees.EmployeesDB{DB: database}
	employee, err := employeeDb.FindByOldID(event.Employee_Id)
	if err != nil {
		return err
	}
	event.Employee_Id = employee.Id
	_, err = eventsDb.Add(&event)
	if err != nil {
		return err
	}
	return nil
}
