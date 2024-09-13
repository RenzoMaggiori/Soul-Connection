package customers

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

type CustomersDB struct {
	DB     *sql.DB
	Bucket *gridfs.Bucket
}

type AddCustomer struct {
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
	Employee_Id        *int
}

type UpdateCustomer struct {
	Email             *string
	Name              *string
	Surname           *string
	Birth_Date        *string
	Gender            *string
	Description       *string
	Astrological_Sign *string
	Phone_Number      *string
	Address           *string
	Employee_Id       *int
}

func (db CustomersDB) FindAll() ([]Customer, error) {
	rows, err := db.DB.Query("SELECT * FROM customer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []Customer

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (db CustomersDB) FindByID(id int) (*Customer, error) {
	query := "SELECT * FROM customer WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var c Customer

	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db CustomersDB) FindByEmployeeID(id int) ([]Customer, error) {
	query := "SELECT * FROM customer WHERE employee_id = $1"

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var customers []Customer

	for rows.Next() {
		var c Customer
	err := rows.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (db CustomersDB) FindByOldID(id int) (*Customer, error) {
	query := "SELECT * FROM customer WHERE soul_connection_id = $1"

	row := db.DB.QueryRow(query, id)
	var c Customer

	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db CustomersDB) Add(customer *AddCustomer) (*Customer, error) {
	query := `
		INSERT INTO customer (soul_connection_id, email, name, surname, birth_date, gender, description, astrological_sign, phone_number, address, employee_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING *
    `
	row := db.DB.QueryRow(query, customer.Soul_Connection_Id, customer.Email, customer.Name, customer.Surname, customer.Birth_Date, customer.Gender, customer.Description, customer.Astrological_Sign, customer.Phone_Number, customer.Address, customer.Employee_Id)
	var c Customer

	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)

	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db CustomersDB) Delete(id int) error {
	query := "DELETE FROM customer WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db CustomersDB) Patch(id int, updates *UpdateCustomer) (*Customer, error) {
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
		"UPDATE customer SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var c Customer
	err := row.Scan(&c.Id, &c.Soul_Connection_Id, &c.Email, &c.Name, &c.Surname, &c.Birth_Date, &c.Gender, &c.Description, &c.Astrological_Sign, &c.Phone_Number, &c.Address, &c.Image_Id, &c.CreatedAt, &c.Employee_Id)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db CustomersDB) UploadFile(customerId int, r io.Reader, filename string) (*primitive.ObjectID, error) {
	fileId, err := filestorage.Upload(db.Bucket, r, filename)
	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	query := `
    UPDATE customer SET image_id = $1 WHERE id = $2
    `
	_, err = tx.Exec(query, fileId.Hex(), customerId)
	if err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		tx.Rollback()
		lib.ServerLog("ERROR", fmt.Errorf("Failed to update customer image_id: %v", err))
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		lib.ServerLog("ERROR", "Could not commit customer")
		return nil, err
	}
	return fileId, nil
}

func (db CustomersDB) GetFile(fileId primitive.ObjectID) ([]byte, error) {
	return filestorage.DownloadById(db.Bucket, fileId)
}
