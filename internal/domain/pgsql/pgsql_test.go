package pgsql_test

import (
	"os"
	"testing"
)

var (
	databaseUrl string
)

func TestMain(m *testing.M) {
	databaseUrl = os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "host=localhost port=5432 user=void password=void123 dbname=void_test sslmode=disable"
	}

	os.Exit(m.Run())
}
