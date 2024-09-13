package employees

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

type EmployeesDB struct {
	DB     *sql.DB
	Bucket *gridfs.Bucket
}

type AddEmployee struct {
	Soul_Connection_Id *int
	Email              string
	// Password  string
	Name       string
	Surname    string
	Birth_Date string
	Gender     string
	Work       string
}

type UpdateEmployee struct {
	Email *string
	// Password  string
	Name       *string
	Surname    *string
	Birth_Date *string
	Gender     *string
	Work       *string
}

func (db EmployeesDB) FindAll() ([]Employee, error) {
	rows, err := db.DB.Query("SELECT * FROM employee")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var employees []Employee

	for rows.Next() {
		var e Employee
		err := rows.Scan(&e.Id, &e.Soul_Connection_Id, &e.Email, &e.Password, &e.Name, &e.Surname, &e.Birth_Date, &e.Gender, &e.Work, &e.Image_Id, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (db EmployeesDB) FindByID(id int) (*Employee, error) {
	query := "SELECT * FROM employee WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var e Employee

	err := row.Scan(&e.Id, &e.Soul_Connection_Id, &e.Email, &e.Password, &e.Name, &e.Surname, &e.Birth_Date, &e.Gender, &e.Work, &e.Image_Id, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db EmployeesDB) FindByOldID(id int) (*Employee, error) {
	query := "SELECT * FROM employee WHERE soul_connection_id = $1"

	row := db.DB.QueryRow(query, id)
	var e Employee

	err := row.Scan(&e.Id, &e.Soul_Connection_Id, &e.Email, &e.Password, &e.Name, &e.Surname, &e.Birth_Date, &e.Gender, &e.Work, &e.Image_Id, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db EmployeesDB) Add(employee *AddEmployee) (*Employee, error) {
	query := `
		INSERT INTO employee (soul_connection_id, email, password, name, surname, birth_date, gender, work)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *
    `
	// NOTE: Password is hard coded till we can retrieve it
	var e Employee
	row := db.DB.QueryRow(query, employee.Soul_Connection_Id, employee.Email, "password", employee.Name, employee.Surname, employee.Birth_Date, employee.Gender, employee.Work)
	err := row.Scan(&e.Id, &e.Soul_Connection_Id, &e.Email, &e.Password, &e.Name, &e.Surname, &e.Birth_Date, &e.Gender, &e.Work, &e.Image_Id, &e.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (db EmployeesDB) Delete(id int) error {
	query := "DELETE FROM employee WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db EmployeesDB) Patch(id int, updates *UpdateEmployee) (*Employee, error) {
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
		"UPDATE employee SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var e Employee
	err := row.Scan(&e.Id, &e.Soul_Connection_Id, &e.Email, &e.Password, &e.Name, &e.Surname, &e.Birth_Date, &e.Gender, &e.Work, &e.Image_Id, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db EmployeesDB) UploadFile(employeeId int, r io.Reader, filename string) (*primitive.ObjectID, error) {
	fileId, err := filestorage.Upload(db.Bucket, r, filename)
	if err != nil {
		return nil, err
	}

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	query := `
    UPDATE employee SET image_id = $1 WHERE id = $2
    `
	_, err = tx.Exec(query, fileId.Hex(), employeeId)
	if err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		tx.Rollback()
		lib.ServerLog("ERROR", fmt.Errorf("Failed to update employee image_id: %v", err))
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		filestorage.Delete(db.Bucket, *fileId)
		lib.ServerLog("ERROR", "Could not commit employee")
		return nil, err
	}
	return fileId, nil
}

func (db EmployeesDB) GetFile(fileId primitive.ObjectID) ([]byte, error) {
	return filestorage.DownloadById(db.Bucket, fileId)
}
