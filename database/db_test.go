package database

import (
	_ "github.com/glebarez/go-sqlite"
	"os"
	"testing"
)

func TestSetup(t *testing.T) {
	// Define test parameters
	testDBName := "test.db"
	testQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER
		);
	`

	// Clean up after the test
	defer func() {
		os.Remove(testDBName)
	}()

	// Run the test
	t.Run("Setup", func(t *testing.T) {
		db, Dbclose := Setup(testQuery, testDBName)
		defer Dbclose()

		// Check if Setup returns a valid database connection
		if db == nil {
			t.Fatal("Expected non-nil database connection, got nil")
		}

		// Check if the table has been created
		rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='users';")
		if err != nil {
			t.Fatalf("Error querying table existence: %v", err)
		}
		defer rows.Close()

		if !rows.Next() {
			t.Fatal("Expected 'users' table to exist, but it doesn't")
		}
	})
}
