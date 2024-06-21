package main

import (
	"github.com/althaafka/alk-proj-be.git/database"
	"github.com/althaafka/alk-proj-be.git/router"
)

func main() {
	database.Connect()
	router.SetupRouter()
}