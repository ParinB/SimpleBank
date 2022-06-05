package main

import (
	_ "github.com/lib/pq"
	"database/sql"
	"github.com/parin/simplebank/api"
	db "github.com/parin/simplebank/db/sqlc"
	"github.com/parin/simplebank/db/util"
	"log"
)

func main()  {
	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:",err)
	}
	conn , err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil {
		log.Fatalln("cannot connect to db:",err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}
}
