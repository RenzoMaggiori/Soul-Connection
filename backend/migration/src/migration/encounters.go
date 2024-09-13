package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"soul-connection.com/api/src/endpoints/customers"
	"soul-connection.com/api/src/endpoints/encounters"
	"soul-connection.com/api/src/lib"
)

func migrateEncounters(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials) error {
	time.Sleep(2 * time.Second)
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/encounters", lib.ApiBaseUri),
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
	var encountersResponse []struct {
		Id          int
		Customer_Id int
		Date        string
		Rating      int
	}
	err = json.NewDecoder(resp.Body).Decode(&encountersResponse)
	if err != nil {
		return err
	}

	lib.ServerLog("PROGRESS", fmt.Sprintf("Migrating encounters...:START:%d", len(encountersResponse)))
	for _, e := range encountersResponse {
		lib.ServerLog("PROGRESS", "Migrating encounters...:INCREMENT")
		migrateEncounter(database, credentials, e.Id)
	}
	lib.ServerLog("PROGRESS", "Migrating encounters...:COMPLETE")
	return nil
}

func migrateEncounter(database *sql.DB, credentials *ApiCredentials, id int) error {
	var encounter encounters.AddEncounter
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/encounters/%d", lib.ApiBaseUri, id),
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
	err = json.NewDecoder(resp.Body).Decode(&encounter)
	if err != nil {
		return err
	}

	encounterDb := encounters.EncountersDB{DB: database}
	customersDb := customers.CustomersDB{DB: database}
	customer, err := customersDb.FindByOldID(encounter.Customer_Id)
	if err != nil {
		return err
	}
	encounter.Customer_Id = customer.Id
	_, err = encounterDb.Add(&encounter)
	if err != nil {
		return err
	}
	return nil
}
