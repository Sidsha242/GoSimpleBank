package main

import (
	"database/sql"
	"log"

	"github.com/Sidsha242/api"
	db "github.com/Sidsha242/db/sqlc"
	"github.com/Sidsha242/util"
	_ "github.com/lib/pq"
)

//we need to connect to database and create a store

func main() {
	//entry point for server

	//getting configs from env variables and using viper

	config, err := util.LoadConfig("./") //loading config from current directory
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	//establish connection to database
	conn, err := sql.Open(config.DBDriber,config.DBSource ) //creates a new sql.DB object
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store:= db.NewStore(conn) //with connection we can create a store

	server:= api.NewServer(store) //with store creating new server

	err = server.Run(config.ServerAddress) //running server on port 8080
	if err != nil {
		log.Fatal("cannot run server: ", err)
	}

}

//need new port for lib/pq driver wihtout this we cannot talk to database