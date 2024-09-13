package payments

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type PaymentsDB struct {
	DB *sql.DB
}

type AddPayment struct {
	Soul_Connection_Id *int
	Date               string
	PaymentMethod      string
	Amount             float64
	Comment            string
	CustomerId         int
}

type UpdatePayment struct {
	Date           *string
	Payment_Method *string
	Amount         *float64
	Comment        *string
}

func (db PaymentsDB) FindAll() ([]Payment, error) {
	rows, err := db.DB.Query("SELECT * FROM payment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []Payment

	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.Id, &p.Soul_Connection_Id, &p.Date, &p.PaymentMethod, &p.Amount, &p.Comment, &p.CreatedAt, &p.CustomerId)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	err = rows.Err()

	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (db PaymentsDB) FindByID(id int) (*Payment, error) {
	query := "SELECT * FROM payment WHERE id = $1"

	row := db.DB.QueryRow(query, id)
	var p Payment

	err := row.Scan(&p.Id, &p.Soul_Connection_Id, &p.Date, &p.PaymentMethod, &p.Amount, &p.Comment, &p.CreatedAt, &p.CustomerId)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db PaymentsDB) FindByCustomerID(id int) ([]Payment, error) {
	query := "SELECT * FROM payment WHERE customer_id = $1"

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var p Payment
		err := rows.Scan(&p.Id, &p.Soul_Connection_Id, &p.Date, &p.PaymentMethod, &p.Amount, &p.Comment, &p.CreatedAt, &p.CustomerId)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (db PaymentsDB) Add(payment *AddPayment) (*Payment, error) {
	query := `
		INSERT INTO payment (soul_connection_id, date, payment_method, amount, comment, customer_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
    `
	row := db.DB.QueryRow(query, payment.Soul_Connection_Id, payment.Date, payment.PaymentMethod, payment.Amount, payment.Comment, payment.CustomerId)
	var p Payment

	err := row.Scan(&p.Id, &p.Soul_Connection_Id, &p.Date, &p.PaymentMethod, &p.Amount, &p.Comment, &p.CreatedAt, &p.CustomerId)

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (db PaymentsDB) Delete(id int) error {
	query := "DELETE FROM payment WHERE id = $1"

	_, err := db.DB.Exec(query, id)
	return err
}

func (db PaymentsDB) Patch(id int, updates *UpdatePayment) (*Payment, error) {
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
		"UPDATE payment SET %s WHERE id = $%d RETURNING *",
		strings.Join(setClauses, ", "),
		argIndex,
	)

	args = append(args, id)

	row := db.DB.QueryRow(query, args...)
	var p Payment
	err := row.Scan(&p.Id, &p.Soul_Connection_Id, &p.Date, &p.PaymentMethod, &p.Amount, &p.Comment, &p.CreatedAt, &p.CustomerId)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
