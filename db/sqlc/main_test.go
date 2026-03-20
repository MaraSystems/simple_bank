package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(t *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Error connecting to Database:", err)
	}

	testQueries = New(testDb)
	os.Exit(t.Run())
}
