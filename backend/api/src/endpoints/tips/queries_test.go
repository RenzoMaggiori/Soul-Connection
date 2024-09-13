package tips_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"soul-connection.com/api/src/endpoints/tips"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema := `
    CREATE TABLE IF NOT EXISTS "tip" (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title VARCHAR(255) NOT NULL,
        tip VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT unique_tip UNIQUE (title, tip)
    );
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func testAdd(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
		defer db.Close()

		tipsDB := tips.TipsDB{DB: db}
		newTip := &tips.AddTip{
			Title: "Test Title",
			Tip:   "This is a test tip.",
		}

		addedTip, err := tipsDB.Add(newTip)
		if err != nil {
			t.Errorf("Failed to add tip: %v", err)
		}

		if addedTip.Title != newTip.Title || addedTip.Tip != newTip.Tip {
			t.Errorf("Expected %v, got %v", newTip, addedTip)
		}
	})
}

func testFindAll(t *testing.T) {
	t.Run("Find All", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
		defer db.Close()

		tipsDB := tips.TipsDB{DB: db}
		newTip := &tips.AddTip{Title: "Test Title", Tip: "This is a test tip."}

		_, err = tipsDB.Add(newTip)
		if err != nil {
			t.Errorf("Failed to add tip: %v", err)
		}

		tips, err := tipsDB.FindAll()
		if err != nil {
			t.Errorf("Failed to find all tips: %v", err)
		}

		if len(tips) != 1 {
			t.Errorf("Expected 1 tip, got %d", len(tips))
		}

		if newTip.Title != tips[0].Title || newTip.Tip != tips[0].Tip {
			t.Errorf("Expected %v, got %v", newTip, tips[0])
		}
	})
}

func testFindByID(t *testing.T) {
	t.Run("Find by ID", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
		defer db.Close()

		tipsDB := tips.TipsDB{DB: db}
		newTip := &tips.AddTip{Title: "Test Title", Tip: "This is a test tip."}

		addedTip, err := tipsDB.Add(newTip)
		if err != nil {
			t.Errorf("Failed to add tip: %v", err)
		}

		tip, err := tipsDB.FindByID(addedTip.Id)
		if err != nil {
			t.Errorf("Failed to find tip by ID: %v", err)
		}

		if tip == nil || addedTip.Id != tip.Id || addedTip.Title != tip.Title || addedTip.Tip != addedTip.Tip {
			t.Errorf("Expected %v, got %v", addedTip, tip)
		}
	})
}

func testPatch(t *testing.T) {
	t.Run("Patch", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
		defer db.Close()

		tipsDB := tips.TipsDB{DB: db}
		newTip := &tips.AddTip{Title: "Test Title", Tip: "This is a test tip."}

		addedTip, err := tipsDB.Add(newTip)
		if err != nil {
			t.Errorf("Failed to add tip: %v", err)
		}

		updatedTitle := "Updated Title"
		updatedTipDescription := "Updated Tip"
		updatedTip := &tips.UpdateTip{
			Title: &updatedTitle,
			Tip:   &updatedTipDescription,
		}

		patchedTip, err := tipsDB.Patch(addedTip.Id, updatedTip)
		if err != nil {
			t.Errorf("Failed to update tip: %v", err)
		}

		if *updatedTip.Title != patchedTip.Title || *updatedTip.Tip != patchedTip.Tip {
			t.Errorf("Expected %v, got %v", updatedTip, patchedTip)
		}
	})
}

func testDelete(t *testing.T) {
	t.Run("Delete", func(t *testing.T) {
		db, err := setupTestDB()
		if err != nil {
			t.Fatalf("Failed to set up test database: %v", err)
		}
		defer db.Close()

		tipsDB := tips.TipsDB{DB: db}
		newTip := &tips.AddTip{Title: "Test Title", Tip: "This is a test tip."}

		addedTip, err := tipsDB.Add(newTip)
		if err != nil {
			t.Errorf("Failed to add tip: %v", err)
		}

		err = tipsDB.Delete(addedTip.Id)
		if err != nil {
			t.Errorf("Failed to delete tip: %v", err)
		}

		tip, err := tipsDB.FindByID(addedTip.Id)
		if err == nil || tip != nil {
			t.Errorf("Expected error when finding deleted tip, got nil")
		}
	})
}

func TestTipQueries(t *testing.T) {
	testAdd(t)
	testFindAll(t)
	testFindByID(t)
	testPatch(t)
	testDelete(t)
}
