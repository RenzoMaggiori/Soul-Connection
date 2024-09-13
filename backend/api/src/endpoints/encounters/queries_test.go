package encounters

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
	CREATE TABLE encounter (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		rating INTEGER NOT NULL,
		comment TEXT NOT NULL,
		source TEXT NOT NULL,
		customer_id INTEGER,
		UNIQUE (date, comment, source, customer_id)
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

func compareEncounters(expected *AddEncounter, actual *Encounter) bool {
	return expected.Date == actual.Date &&
		expected.Rating == actual.Rating &&
		expected.Comment == actual.Comment &&
		expected.Source == actual.Source &&
		reflect.DeepEqual(expected.Customer_Id, actual.Customer_Id)
}

func TestEncounterQueries(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	encountersDB := EncountersDB{DB: db}
	customer_ID := 1

	t.Run("Add Encounter", func(t *testing.T) {
		newEncounter := &AddEncounter{
			Date:        "25-04-2023",
			Rating:      5,
			Comment:     "Great encounter",
			Source:      "Referral",
			Customer_Id: customer_ID,
		}

		addedEncounter, err := encountersDB.Add(newEncounter)
		if err != nil {
			t.Errorf("Failed to add encounter: %v", err)
			return
		}

		if !compareEncounters(newEncounter, addedEncounter) {
			t.Errorf("Expected %+v, got %+v", newEncounter, addedEncounter)
		}
	})

	t.Run("Find All Encounters", func(t *testing.T) {
		_, err := encountersDB.Add(&AddEncounter{
			Date:        "26-04-2023",
			Rating:      4,
			Comment:     "Good encounter",
			Source:      "Website",
			Customer_Id: customer_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add encounter: %v", err)
		}

		encounters, err := encountersDB.FindAll()
		if err != nil {
			t.Errorf("Failed to find all encounters: %v", err)
			return
		}

		if len(encounters) != 2 {
			t.Errorf("Expected 2 encounters, got %d", len(encounters))
		}
	})

	t.Run("Find Encounter by ID", func(t *testing.T) {
		_, err := encountersDB.Add(&AddEncounter{
			Date:        "27-04-2023",
			Rating:      3,
			Comment:     "Average encounter",
			Source:      "Email",
			Customer_Id: customer_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add encounter: %v", err)
		}

		encounter, err := encountersDB.FindByID(1)
		if err != nil {
			t.Errorf("Failed to find encounter by ID: %v", err)
			return
		}

		if encounter == nil || encounter.Id != 1 {
			t.Errorf("Expected encounter with ID 1, got %+v", encounter)
		}
	})

	t.Run("Patch Encounter", func(t *testing.T) {
		_, err := encountersDB.Add(&AddEncounter{
			Date:        "28-04-2023",
			Rating:      2,
			Comment:     "Below average encounter",
			Source:      "Social Media",
			Customer_Id: customer_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add encounter: %v", err)
		}

		newComment := "Updated comment"
		updates := &UpdateEncounter{Comment: &newComment}

		updatedEncounter, err := encountersDB.Patch(1, updates)
		if err != nil {
			t.Errorf("Failed to update encounter: %v", err)
			return
		}

		if updatedEncounter.Comment != newComment {
			t.Errorf("Expected comment %s, got %s", newComment, updatedEncounter.Comment)
		}
	})

	t.Run("Delete Encounter", func(t *testing.T) {
		_, err := encountersDB.Add(&AddEncounter{
			Date:        "29-04-2023",
			Rating:      1,
			Comment:     "Bad encounter",
			Source:      "Direct",
			Customer_Id: customer_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add encounter: %v", err)
		}

		err = encountersDB.Delete(1)
		if err != nil {
			t.Errorf("Failed to delete encounter: %v", err)
			return
		}

		_, err = encountersDB.FindByID(1)
		if err == nil {
			t.Errorf("Expected error when finding deleted encounter, got nil")
		}
	})
}
