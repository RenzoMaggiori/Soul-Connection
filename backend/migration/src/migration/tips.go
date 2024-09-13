package migration

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"soul-connection.com/api/src/endpoints/tips"
	"soul-connection.com/api/src/lib"
)

func migrateTips(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials) error {
	time.Sleep(2 * time.Second)
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/tips", lib.ApiBaseUri),
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
	var tip []tips.AddTip

	err = json.NewDecoder(resp.Body).Decode(&tip)
	if err != nil {
		return err
	}

	tipsDb := tips.TipsDB{DB: database}
	lib.ServerLog("PROGRESS", fmt.Sprintf("Migrating tips...:START:%d", len(tip)))
	for _, e := range tip {
		lib.ServerLog("PROGRESS", "Migrating tips...:INCREMENT")
		tipsDb.Add(&e)
	}
	lib.ServerLog("PROGRESS", "Migrating tips...:COMPLETE")
	return nil
}
