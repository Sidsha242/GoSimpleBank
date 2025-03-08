package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq" //underscore (blank identifier) to import a package for its side-effects only (initialization of the package). if not will be removed on saving file
)


const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable"
)
var testQueries *Queries

//Global variable to store the connection to the database
var testDB *sql.DB

//Wrapper function to run all unit tests
//entry point for all testing package db  ..used by accounts_test.go
func TestMain(m *testing.M) {
		var err error
		testDB, err = sql.Open(dbDriver,  dbSource)  //creates a new sql.DB object
		if err != nil {
			log.Fatal("cannot connect to db: ", err)
		}

		testQueries = New(testDB) //passing connection to the queries

		os.Exit(m.Run())  //to run all unit tests
}