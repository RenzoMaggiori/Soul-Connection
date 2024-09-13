package tips

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type TipsDB struct {
	DB *sql.DB
}

type AddTip struct {
	Title string
	Tip   string
}

type UpdateTip struct {
	Title *string
	Tip   *string
}

func (db TipsDB) FindAll() ([]Tip, error) {
	rows, err := db.DB.Query("SELECT * FROM tip")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tip []Tip

	for rows.Next() {
		var t Tip
		err := rows.Scan(&t.Id, &t.Title, &t.Tip, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tip = append(tip, t)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return tip, nil
}

func (db TipsDB) FindByID(id int) (*Tip, error) {
	query := "SELECT * FROM tip WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var t Tip

	err := row.Scan(&t.Id, &t.Title, &t.Tip, &t.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (db TipsDB) Add(tip *AddTip) (*Tip, error) {
	query := `
		INSERT INTO tip (title, tip)
		VALUES ($1, $2)
		RETURNING *
    `
	row := db.DB.QueryRow(query, tip.Title, tip.Tip)
	var t Tip

	err := row.Scan(&t.Id, &t.Title, &t.Tip, &t.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (db TipsDB) Delete(id int) error {
	query := "DELETE FROM tip WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db TipsDB) Patch(id int, updates *UpdateTip) (*Tip, error) {
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
	var tip Tip
	err := row.Scan(&tip.Id, &tip.Title, &tip.Tip, &tip.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &tip, nil
}
