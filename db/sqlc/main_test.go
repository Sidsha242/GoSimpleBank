package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Sidsha242/util"
	_ "github.com/lib/pq" //underscore (blank identifier) to import a package for its side-effects only (initialization of the package). if not will be removed on saving file
)


var testQueries *Queries

//Global variable to store the connection to the database
var testDB *sql.DB

//Wrapper function to run all unit tests
//entry point for all testing package db  ..used by accounts_test.go
func TestMain(m *testing.M) {
		//getting configs

		config, err := util.LoadConfig("../..") //loading config from current directory
		if err != nil {
			log.Fatal("cannot load config: ", err)
		}

		testDB, err = sql.Open(config.DBDriber,  config.DBSource)  //creates a new sql.DB object
		if err != nil {
			log.Fatal("cannot connect to db: ", err)
		}

		testQueries = New(testDB) //passing connection to the queries

		os.Exit(m.Run())  //to run all unit tests
}