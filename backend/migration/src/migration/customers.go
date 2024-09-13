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
	"soul-connection.com/api/src/endpoints/customers"
	"soul-connection.com/api/src/lib"
)

func migrateCustomers(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials) error {
	time.Sleep(2 * time.Second)
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/customers", lib.ApiBaseUri),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("could not retrieve customers from api")
	}
	var customersResponse []struct {
		Id      int
		Email   string
		Name    string
		Surname string
	}
	err = json.NewDecoder(resp.Body).Decode(&customersResponse)
	if err != nil {
		return err
	}

	lib.ServerLog("PROGRESS", fmt.Sprintf("Migrating customers, clothes and payments...:START:%d", len(customersResponse)))
	for _, e := range customersResponse {
		lib.ServerLog("PROGRESS", "Migrating customers, clothes and payments...:INCREMENT")
		migrateCustomer(database, fileStorage, credentials, e.Id)
	}
	lib.ServerLog("PROGRESS", "Migrating customers, clothes and payments...:COMPLETE")
	return nil
}

func migrateCustomer(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials, id int) error {
	var customer struct {
		Id                int
		Email             string
		Name              string
		Surname           string
		Birth_Date        string
		Gender            string
		Description       string
		Astrological_Sign string
		Phone_Number      string
		Address           string
	}
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/customers/%d", lib.ApiBaseUri, id),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve customer with id %d from api", id)
	}

	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		return err
	}

	bucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("customersBucket"))
	if err != nil {
		return err
	}
	customerDb := &customers.CustomersDB{DB: database, Bucket: bucket}

	newCustomer, err := customerDb.Add(&customers.AddCustomer{
		Soul_Connection_Id: &customer.Id,
		Email:              customer.Email,
		Name:               customer.Name,
		Surname:            customer.Surname,
		Birth_Date:         customer.Birth_Date,
		Gender:             customer.Gender,
		Description:        customer.Description,
		Astrological_Sign:  customer.Astrological_Sign,
		Phone_Number:       customer.Phone_Number,
		Address:            customer.Address,
	})
	if err != nil {
		return err
	}
	ids := &Ids{
		old: *newCustomer.Soul_Connection_Id,
		new: newCustomer.Id,
	}
	migrateCustomerImage(customerDb, credentials, ids)
	migrateClothes(database, fileStorage, credentials, ids)
	migratePayments(database, fileStorage, credentials, ids)
	return nil
}

func migrateCustomerImage(db *customers.CustomersDB, credentials *ApiCredentials, ids *Ids) error {
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/customers/%d/image", lib.ApiBaseUri, ids.old),
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

	_, err = db.UploadFile(ids.new, resp.Body, fmt.Sprintf("customer_%d", ids.new))
	if err != nil {
		return err
	}
	return nil
}
