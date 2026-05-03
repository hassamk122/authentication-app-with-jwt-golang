package main

import (
	"fmt"

	"github.com/hassamk122/authentication-app-with-jwt-golang/config"
)

func main() {
	conf, err := config.LoadConfig()
	if err == nil {
		fmt.Printf("server port: %s\n", conf.ServerPort)
		fmt.Printf("Database url: %s\n", conf.DatabaseUrl)
		fmt.Printf("environment: %s\n", conf.Environment)
	}

}
