package database

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestConnectionString(t *testing.T) {
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_port := os.Getenv("POSTGRES_PORT")

	expected := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", postgres_user, postgres_password, postgres_host, postgres_port, postgres_db)

	got := ConnectionString()
	if got != expected {
		t.Errorf("ConnectionString() = %v; want %v", got, expected)
	}
}
