package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hassamk122/authentication-app-with-jwt-golang/config"
)

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	db := config.ConnectToDB(configuration.DatabaseUrl)
	defer db.Close()

	fmt.Printf("server port: %s\n", configuration.ServerPort)
	fmt.Printf("Database url: %s\n", configuration.DatabaseUrl)
	fmt.Printf("environment: %s\n", configuration.Environment)

}
