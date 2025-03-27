package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Sidsha242/simple_bank/util"
	_ "github.com/lib/pq" //underscore (blank identifier) to import a package for its side-effects only (initialization of the package). if not will be removed on saving file
)

//entrypoint for running all unit tests in the db package

// It sets up the necessary environment for testing by connecting to the database and initializing the Queries object.
// It uses the testing package in Go to manage and execute tests.

//Blobal variable to ecex database queries (generated by sqlc).
var testQueries *Queries

//Global variable to store the connection to the database
var testDB *sql.DB

//Wrapper function to run all unit tests
//entry point for all testing package db  ..used by accounts_test.go
func TestMain(m *testing.M) {
		//getting configs

		config, err := util.LoadConfig("../..") //loading database config from config file
		if err != nil {
			log.Fatal("cannot load config: ", err)
		}

		testDB, err = sql.Open(config.DBDriver,  config.DBSource)  //creates a new sql.DB object  which represents the connection to the database.
		if err != nil {
			log.Fatal("cannot connect to db: ", err)
		}

		testQueries = New(testDB) //create a new Queries object, passing the database connection.
		//This object is used to execute SQL queries in the test files.

		os.Exit(m.Run())   //m.Run() to execute all tests in the package.
		//os.Exit ensures that the program exits with the appropriate status code after the tests are complete.

}