package migration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"soul-connection.com/api/src/endpoints/clothes"
	"soul-connection.com/api/src/lib"
)

func migrateClothes(database *sql.DB, fileStorage *mongo.Database, credentials *ApiCredentials, ids *Ids) error {
	var cs []struct {
		Id   int
		Type string
	}
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/customers/%d/clothes", lib.ApiBaseUri, ids.old),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve clothes from user with id %d from api", ids.old)
	}
	err = json.NewDecoder(resp.Body).Decode(&cs)
	if err != nil {
		return err
	}

	bucket, err := gridfs.NewBucket(fileStorage, options.GridFSBucket().SetName("clothesBucket"))
	if err != nil {
		return err
	}
	clothesDb := clothes.ClothesDB{DB: database, Bucket: bucket}

	for _, c := range cs {
		newClothe, err := clothesDb.Add(&clothes.AddClothe{
			Soul_Connection_Id: &c.Id,
			Type:               c.Type,
			CustomerId:         ids.new,
		})
		if err != nil {
			continue
		}

		migrateClotheImage(&clothesDb, credentials, &Ids{old: c.Id, new: newClothe.Id})
	}

	return nil
}

func migrateClotheImage(db *clothes.ClothesDB, credentials *ApiCredentials, ids *Ids) error {
	resp, err := lib.Fetch(&http.Client{}, lib.FetchRequest{
		Method:  "GET",
		Url:     fmt.Sprintf("%s/api/clothes/%d/image", lib.ApiBaseUri, ids.old),
		Body:    nil,
		Headers: map[string]string{"X-Group-Authorization": credentials.XGroupAuthentication, "Authorization": fmt.Sprintf("Bearer %s", credentials.Jwt)},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("could not retrieve clothe image with id %d from api", ids.old)
	}

	_, err = db.UploadFile(ids.new, resp.Body, fmt.Sprintf("clothe_%d", ids.new))
	if err != nil {
		return err
	}
	return nil
}
