package main

import (
	"fmt"
	"os"

	"gitlab.com/faemproject/backend/hack/db"
	"gitlab.com/faemproject/backend/hack/logs"
	"gitlab.com/faemproject/backend/hack/models"
	"gitlab.com/faemproject/backend/hack/router"
)

func init() {
}

const version = "0.2.1"

func main() {

	eloger := logs.InitElogger("crm")

	// logs.OutputInfo("Connecting to Elasticsearch successfuly")
	fmt.Println("Connecting to Elasticsearch successfuly")

	conn, err := db.Connect()
	if err != nil {
		es := fmt.Sprintf("Error connecting to DB. %s", err)
		// logs.OutputError(es)
		fmt.Println(es)
		os.Exit(4)
	}
	defer db.CloseDbConnection(conn)
	fmt.Println("Connecting to Postgres successfuly", conn)
	// logs.OutputInfo("Connecting to Postgres successfuly")

	models.ConnectDB(conn)

	e := router.Init(eloger)
	port := fmt.Sprintf(":%v", "1324")
	e.Logger.Fatal(e.Start(port))
}
