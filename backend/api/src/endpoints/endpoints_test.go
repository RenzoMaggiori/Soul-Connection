package endpoints

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"
//
// 	"github.com/joho/godotenv"
// 	_ "github.com/mattn/go-sqlite3"
// )
//
// func loadEnv(t *testing.T) {
// 	err := godotenv.Load("../../../.env")
// 	if err != nil {
// 		t.Fatalf("Error loading .env file: %v", err)
// 	}
// }
//
// func setupTestDB() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", ":memory:")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	schema := `
// 	CREATE TABLE IF NOT EXISTS employee (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		soul_connection_id INTEGER,
// 		email TEXT NOT NULL,
// 		password TEXT NOT NULL,
// 		surname TEXT NOT NULL,
// 		ad TEXT,
// 		gender TEXT,
// 		work TEXT
// 	);
// 	`
//
// 	_, err = db.Exec(schema)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return db, nil
// }
//
// func TestCreateRouter(t *testing.T) {
// 	loadEnv(t)
//
// 	email := os.Getenv("API_EMAIL")
// 	password := os.Getenv("API_PASSWORD")
//
// 	if email == "" || password == "" {
// 		t.Fatalf("Email or password not set in .env file")
// 	}
//
// 	db, err := setupTestDB()
// 	if err != nil {
// 		t.Fatalf("Failed to set up test database: %v", err)
// 	}
// 	defer db.Close()
//
// 	router := CreateRouter(db)
//
// 	testCases := []struct {
// 		name           string
// 		method         string
// 		url            string
// 		expectedStatus int
// 		body           interface{}
// 	}{
// 		{"Login Route", http.MethodPost, "/api/auth/login", http.StatusOK, map[string]string{"email": email, "password": password}},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var req *http.Request
// 			if tc.body != nil {
// 				jsonBody, err := json.Marshal(tc.body)
// 				if err != nil {
// 					t.Fatalf("Failed to encode JSON body: %v", err)
// 				}
// 				req = httptest.NewRequest(tc.method, tc.url, bytes.NewBuffer(jsonBody))
// 				req.Header.Set("Content-Type", "application/json")
// 			} else {
// 				req, err = http.NewRequest(tc.method, tc.url, nil)
// 				if err != nil {
// 					t.Fatalf("Failed to create request: %v", err)
// 				}
// 			}
//
// 			rr := httptest.NewRecorder()
// 			router.ServeHTTP(rr, req)
//
// 			t.Logf("Response Body: %s", rr.Body.String())
//
// 			if rr.Code != tc.expectedStatus {
// 				t.Errorf("Expected status %d, got %d. Response: %s", tc.expectedStatus, rr.Code, rr.Body.String())
// 			}
// 		})
// 	}
// }
