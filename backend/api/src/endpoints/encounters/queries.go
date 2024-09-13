package encounters

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type EncountersDB struct {
	DB *sql.DB
}

type AddEncounter struct {
	Date        string
	Rating      int
	Comment     string
	Source      string
	Customer_Id int
}

type UpdateEncounter struct {
	Date    *string
	Rating  *int
	Comment *string
	Source  *string
}

func (db EncountersDB) FindAll() ([]Encounter, error) {
	rows, err := db.DB.Query("SELECT * FROM encounter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var encounter []Encounter

	for rows.Next() {
		var e Encounter
		err := rows.Scan(&e.Id, &e.Date, &e.Rating, &e.Comment, &e.Source, &e.CreatedAt, &e.Customer_Id)
		if err != nil {
			return nil, err
		}
		encounter = append(encounter, e)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return encounter, nil
}

func (db EncountersDB) FindByID(id int) (*Encounter, error) {
	query := "SELECT * FROM encounter WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var e Encounter

	err := row.Scan(&e.Id, &e.Date, &e.Rating, &e.Comment, &e.Source, &e.CreatedAt, &e.Customer_Id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db EncountersDB) FindByCustomerID(id int) ([]Encounter, error) {
	query := "SELECT * FROM encounter WHERE customer_id = $1"

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var encounter []Encounter

	for rows.Next() {
		var e Encounter
		err := rows.Scan(&e.Id, &e.Date, &e.Rating, &e.Comment, &e.Source, &e.CreatedAt, &e.Customer_Id)
		if err != nil {
			return nil, err
		}
		encounter = append(encounter, e)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return encounter, nil
}

func (db EncountersDB) Add(payment *AddEncounter) (*Encounter, error) {
	query := `
		INSERT INTO encounter (date, rating, comment, source, customer_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
    `
	row := db.DB.QueryRow(query, payment.Date, payment.Rating, payment.Comment, payment.Source, payment.Customer_Id)
	var e Encounter

	err := row.Scan(&e.Id, &e.Date, &e.Rating, &e.Comment, &e.Source, &e.CreatedAt, &e.Customer_Id)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (db EncountersDB) Delete(id int) error {
	query := "DELETE FROM encounter WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db EncountersDB) Patch(id int, updates *UpdateEncounter) (*Encounter, error) {
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
		"UPDATE encounter SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var e Encounter
	err := row.Scan(&e.Id, &e.Date, &e.Rating, &e.Comment, &e.Source, &e.CreatedAt, &e.Customer_Id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
