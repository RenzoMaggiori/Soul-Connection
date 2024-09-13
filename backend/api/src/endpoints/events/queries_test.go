package events

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
	CREATE TABLE event (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		date TEXT NOT NULL,
		max_participants INTEGER NOT NULL,
		location_x TEXT NOT NULL,
		location_y TEXT NOT NULL,
		type TEXT NOT NULL,
		employee_id INTEGER,
		UNIQUE (name, date, location_x, location_y)
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}

func compareEvents(expected *AddEvent, actual *Event) bool {
	return expected.Name == actual.Name &&
		expected.Date == actual.Date &&
		expected.Max_Participants == actual.Max_Participants &&
		expected.Location_X == actual.Location_X &&
		expected.Location_Y == actual.Location_Y &&
		expected.Type == actual.Type &&
		reflect.DeepEqual(expected.Employee_Id, actual.Employee_Id)
}

func TestEventQueries(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}
	defer db.Close()

	eventsDB := EventsDB{DB: db}
	employee_ID := 1

	t.Run("Add Event", func(t *testing.T) {
		newEvent := &AddEvent{
			Name:             "Test Event",
			Date:             "25-04-2023",
			Max_Participants: 100,
			Location_X:       "Location X",
			Location_Y:       "Location Y",
			Type:             "Conference",
			Employee_Id:      employee_ID,
		}

		addedEvent, err := eventsDB.Add(newEvent)
		if err != nil {
			t.Errorf("Failed to add event: %v", err)
			return
		}

		if !compareEvents(newEvent, addedEvent) {
			t.Errorf("Expected %+v, got %+v", newEvent, addedEvent)
		}
	})

	t.Run("Find All Events", func(t *testing.T) {
		_, err := eventsDB.Add(&AddEvent{
			Name:             "Another Event",
			Date:             "26-04-2023",
			Max_Participants: 50,
			Location_X:       "Another Location X",
			Location_Y:       "Another Location Y",
			Type:             "Meeting",
			Employee_Id:      employee_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add event: %v", err)
		}

		events, err := eventsDB.FindAll()
		if err != nil {
			t.Errorf("Failed to find all events: %v", err)
			return
		}

		if len(events) != 2 {
			t.Errorf("Expected 2 events, got %d", len(events))
		}
	})

	t.Run("Find Event by ID", func(t *testing.T) {
		_, err := eventsDB.Add(&AddEvent{
			Name:             "Another Event",
			Date:             "26-04-2023",
			Max_Participants: 50,
			Location_X:       "Another 2 Location X",
			Location_Y:       "Another 2 Location Y",
			Type:             "Meeting",
			Employee_Id:      employee_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add event: %v", err)
		}

		event, err := eventsDB.FindByID(1)
		if err != nil {
			t.Errorf("Failed to find event by ID: %v", err)
			return
		}

		if event == nil || event.Id != 1 {
			t.Errorf("Expected event with ID 1, got %+v", event)
		}
	})

	t.Run("Patch Event", func(t *testing.T) {
		_, err := eventsDB.Add(&AddEvent{
			Name:             "Test Event",
			Date:             "25-04-2023",
			Max_Participants: 100,
			Location_X:       "Location X",
			Location_Y:       "Location Y 3",
			Type:             "Conference",
			Employee_Id:      employee_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add event: %v", err)
		}

		newDate := "25-05-2023"
		updates := &UpdateEvent{Date: &newDate}

		updatedEvent, err := eventsDB.Patch(1, updates)
		if err != nil {
			t.Errorf("Failed to update event: %v", err)
			return
		}

		if updatedEvent.Date != newDate {
			t.Errorf("Expected date %s, got %s", newDate, updatedEvent.Date)
		}
	})

	t.Run("Delete Event", func(t *testing.T) {
		_, err := eventsDB.Add(&AddEvent{
			Name:             "Test Event",
			Date:             "25-04-2023",
			Max_Participants: 100,
			Location_X:       "Location X 4",
			Location_Y:       "Location Y",
			Type:             "Conference",
			Employee_Id:      employee_ID,
		})
		if err != nil {
			t.Fatalf("Failed to add event: %v", err)
		}

		err = eventsDB.Delete(1)
		if err != nil {
			t.Errorf("Failed to delete event: %v", err)
			return
		}

		_, err = eventsDB.FindByID(1)
		if err == nil {
			t.Errorf("Expected error when finding deleted event, got nil")
		}
	})
}
