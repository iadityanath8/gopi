package main

import (
	"github.com/iadityanath8/gopi/database"
	"github.com/iadityanath8/gopi/handlers"
)

func main() {
	err := database.InitDriver()
	if err != nil {
		panic("Failed to connect database")
	}

	router := handlers.SetupRouter()
	router.Run(":4000")
}
