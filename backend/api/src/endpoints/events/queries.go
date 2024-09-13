package events

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type EventsDB struct {
	DB *sql.DB
}

type AddEvent struct {
	Name             string
	Date             string
	Max_Participants int
	Location_X       string
	Location_Y       string
	Type             string
	Employee_Id      int
}

type UpdateEvent struct {
	Name             *string
	Date             *string
	Max_Participants *int
	Location_X       *string
	Location_Y       *string
	Type             *string
}

func (db EventsDB) FindAll() ([]Event, error) {
	rows, err := db.DB.Query("SELECT * FROM event")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []Event

	for rows.Next() {
		var e Event
		err := rows.Scan(&e.Id, &e.Name, &e.Date, &e.Max_Participants, &e.Location_X, &e.Location_Y, &e.Type, &e.CreatedAt, &e.Employee_Id)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return events, nil
}

func (db EventsDB) FindByID(id int) (*Event, error) {
	query := "SELECT * FROM event WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var e Event

	err := row.Scan(&e.Id, &e.Name, &e.Date, &e.Max_Participants, &e.Location_X, &e.Location_Y, &e.Type, &e.CreatedAt, &e.Employee_Id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db EventsDB) Add(employee *AddEvent) (*Event, error) {
	query := `
		INSERT INTO event (name, date, max_participants, location_x, location_y, type, employee_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING *
    `
	row := db.DB.QueryRow(query, employee.Name, employee.Date, employee.Max_Participants, employee.Location_X, employee.Location_Y, employee.Type, employee.Employee_Id)
	var e Event

	err := row.Scan(&e.Id, &e.Name, &e.Date, &e.Max_Participants, &e.Location_X, &e.Location_Y, &e.Type, &e.CreatedAt, &e.Employee_Id)

	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (db EventsDB) Delete(id int) error {
	query := "DELETE FROM event WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db EventsDB) Patch(id int, updates *UpdateEvent) (*Event, error) {
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
		"UPDATE event SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var e Event
	err := row.Scan(&e.Id, &e.Name, &e.Date, &e.Max_Participants, &e.Location_X, &e.Location_Y, &e.Type, &e.CreatedAt, &e.Employee_Id)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
