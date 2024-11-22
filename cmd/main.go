package main

import (
	"log"

	"github.com/youtube/config"
	"github.com/youtube/router"
)

func main() {
	// configure the configuration files
	configSetup, err := config.SetUpConfig()
	if err != nil {
		log.Fatalf("Failed To connect with the Database %v\n", err)
	}

	mongoClient, err := config.NewDB(configSetup)
	if err != nil {
		log.Fatalf("Failed To connect with the Database %v\n", err)
	}
	router := router.SetupRouter(mongoClient)
	router.Run(":8003")
}
