package payments

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE payment (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		soul_connection_id INTEGER,
		date TEXT NOT NULL,
		payment_method TEXT NOT NULL,
		amount REAL NOT NULL,
		comment TEXT NOT NULL,
		customer_id INTEGER NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

func comparePayments(expected *AddPayment, actual *Payment) bool {
	return reflect.DeepEqual(expected.Soul_Connection_Id, actual.Soul_Connection_Id) &&
		expected.Date == actual.Date &&
		expected.PaymentMethod == actual.PaymentMethod &&
		expected.Amount == actual.Amount &&
		expected.Comment == actual.Comment &&
		expected.CustomerId == actual.CustomerId
}

func TestPaymentQueries(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	paymentsDB := PaymentsDB{DB: db}
	num := 1

	t.Run("Add Payment", func(t *testing.T) {
		newPayment := &AddPayment{
			Soul_Connection_Id: &num,
			Date:               "25-04-2023",
			PaymentMethod:      "cash",
			Amount:             1.0,
			Comment:            "comment",
			CustomerId:         1,
		}

		addedPayment, err := paymentsDB.Add(newPayment)
		if err != nil {
			t.Errorf("Failed to add payment: %v", err)
			return
		}

		if !comparePayments(newPayment, addedPayment) {
			t.Errorf("Expected %+v, got %+v", newPayment, addedPayment)
		}
	})

	t.Run("Find All Payments", func(t *testing.T) {
		_, err := paymentsDB.Add(&AddPayment{
			Soul_Connection_Id: &num,
			Date:               "25-04-2023",
			PaymentMethod:      "cash",
			Amount:             1.0,
			Comment:            "comment",
			CustomerId:         1,
		})
		if err != nil {
			t.Fatalf("Failed to add payment: %v", err)
		}

		payments, err := paymentsDB.FindAll()
		if err != nil {
			t.Errorf("Failed to find all payments: %v", err)
			return
		}

		if len(payments) != 2 {
			t.Errorf("Expected 2 payments, got %d", len(payments))
		}
	})

	t.Run("Find Payment by ID", func(t *testing.T) {
		_, err := paymentsDB.Add(&AddPayment{
			Date:          "25-04-2023",
			PaymentMethod: "cash",
			Amount:        1.0,
			Comment:       "comment",
			CustomerId:    1,
		})
		if err != nil {
			t.Fatalf("Failed to add payment: %v", err)
		}

		payment, err := paymentsDB.FindByID(1)
		if err != nil {
			t.Errorf("Failed to find payment by ID: %v", err)
			return
		}

		if payment == nil || payment.Id != 1 {
			t.Errorf("Expected payment with ID 1, got %+v", payment)
		}
	})

	t.Run("Patch Payment", func(t *testing.T) {
		_, err := paymentsDB.Add(&AddPayment{
			Date:          "25-04-2023",
			PaymentMethod: "cash",
			Amount:        1.0,
			Comment:       "comment",
			CustomerId:    1,
		})
		if err != nil {
			t.Fatalf("Failed to add payment: %v", err)
		}

		newDate := "25-05-2023"
		updates := &UpdatePayment{Date: &newDate}

		updatedPayment, err := paymentsDB.Patch(1, updates)
		if err != nil {
			t.Errorf("Failed to update payment: %v", err)
			return
		}

		if updatedPayment.Date != newDate {
			t.Errorf("Expected date %s, got %s", newDate, updatedPayment.PaymentMethod)
		}
	})

	t.Run("Delete Payment", func(t *testing.T) {
		_, err := paymentsDB.Add(&AddPayment{
			Date:          "25-04-2023",
			PaymentMethod: "cash",
			Amount:        1.0,
			Comment:       "comment",
			CustomerId:    1,
		})
		if err != nil {
			t.Fatalf("Failed to add payment: %v", err)
		}

		err = paymentsDB.Delete(1)
		if err != nil {
			t.Errorf("Failed to delete payment: %v", err)
			return
		}

		_, err = paymentsDB.FindByID(1)
		if err == nil {
			t.Errorf("Expected error when finding deleted payment, got nil")
		}
	})
}
