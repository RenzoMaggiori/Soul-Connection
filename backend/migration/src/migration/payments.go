package migration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"soul-connection.com/api/src/endpoints/payments"
	"soul-connection.com/api/src/lib"
)

func migratePayments(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials, ids *Ids) error {
	var ps []struct {
		Id             int
		Date           string
		Payment_Method string
		Amount         float64
		Comment        string
	}
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/customers/%d/payments_history", lib.ApiBaseUri, ids.old),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve payments from user with id %d from api", ids.old)
	}
	err = json.NewDecoder(resp.Body).Decode(&ps)
	if err != nil {
		return err
	}

	paymentsDb := payments.PaymentsDB{DB: database}
	for _, e := range ps {
		paymentsDb.Add(&payments.AddPayment{
			Soul_Connection_Id: &e.Id,
			Date:               e.Date,
			PaymentMethod:      e.Payment_Method,
			Amount:             e.Amount,
			Comment:            e.Comment,
			CustomerId:         ids.new,
		})
	}

	return nil
}
