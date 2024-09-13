package clothes

import (
	"database/sql"
	"fmt"
	"io"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	filestorage "soul-connection.com/api/src/file-storage"
	"soul-connection.com/api/src/lib"
)

type ClothesDB struct {
	DB     *sql.DB
	Bucket *gridfs.Bucket
}

type AddClothe struct {
	Soul_Connection_Id *int
	Type               string
	CustomerId         int
}

type UpdateClothe struct {
	Type *string
}

func (db ClothesDB) FindAll() ([]Clothe, error) {
	rows, err := db.DB.Query("SELECT * FROM clothe")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clothe []Clothe

	for rows.Next() {
		var c Clothe
		err := rows.Scan(&c.Id, &c.Soul_Connection_Id, &c.Type, &c.Image_Id, &c.CreatedAt, &c.CustomerId)
		if err != nil {
			return nil, err
		}
		clothe = append(clothe, c)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return clothe, nil
}

func (db ClothesDB) FindByID(id int) (*Clothe, error) {
	query := "SELECT * FROM clothe WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var c Clothe

	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Type, &c.Image_Id, &c.CreatedAt, &c.CustomerId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db ClothesDB) FindByCustomerID(id int) ([]Clothe, error) {
	query := "SELECT * FROM clothe WHERE customer_id = $1"

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clothes []Clothe

	for rows.Next() {
		var c Clothe
		err := rows.Scan(&c.Id, &c.Soul_Connection_Id, &c.Type, &c.Image_Id, &c.CreatedAt, &c.CustomerId)
		if err != nil {
			return nil, err
		}
		clothes = append(clothes, c)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return clothes, nil
}

func (db ClothesDB) Add(clothe *AddClothe) (*Clothe, error) {
	query := `
		INSERT INTO clothe (soul_connection_id, type, customer_id)
		VALUES ($1, $2, $3)
		RETURNING *
    `
	row := db.DB.QueryRow(query, clothe.Soul_Connection_Id, clothe.Type, clothe.CustomerId)
	var c Clothe

	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Type, &c.Image_Id, &c.CreatedAt, &c.CustomerId)

	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db ClothesDB) Delete(id int) error {
	query := "DELETE FROM clothe WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db ClothesDB) Patch(id int, updates *UpdateClothe) (*Clothe, error) {
	v := reflect.ValueOf(updates).Elem()
	t := reflect.TypeOf(*updates)

	var setClauses []string
	var args []interface{}
	argIndex := 1

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.IsNil() {
			fieldName := t.Field(i).Tag.Get("db")
			if fieldName == "" {
				fieldName = strings.ToLower(strings.Replace(t.Field(i).Name, "_", "", -1))
			}
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", fieldName, argIndex))
			args = append(args, field.Interface())
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(
		"UPDATE tip SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var c Clothe
	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Type, &c.Image_Id, &c.CreatedAt, &c.CustomerId)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db ClothesDB) UploadFile(clotheId int, r io.Reader, filename string) (*primitive.ObjectID, error) {
	fileId, err := filestorage.Upload(db.Bucket, r, filename)
	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	query := `
    UPDATE clothe SET image_id = $1 WHERE id = $2
    `
	_, err = tx.Exec(query, fileId.Hex(), clotheId)
	if err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		tx.Rollback()
		lib.ServerLog("ERROR", fmt.Errorf("Failed to update clothe image_id: %v", err))
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		lib.ServerLog("ERROR", "Could not commit clothe")
		return nil, err
	}
	return fileId, nil
}

func (db ClothesDB) GetFile(fileId primitive.ObjectID) ([]byte, error) {
	return filestorage.DownloadById(db.Bucket, fileId)
}
