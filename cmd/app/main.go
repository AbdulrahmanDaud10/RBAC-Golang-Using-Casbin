package main

import (
	"fmt"
	"log"

	"github.com/AbdulrahmanDaud10/RBAC-Casbin-Golang/pkg/repository"
)

func main() {
	database, err := repository.PostgreSQLDataBaseConnection()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// app.SetUpRoutes(database)
	fmt.Println("Database connection failure", database)
}
