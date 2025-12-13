package main

import (
	"fmt"

	"github.com/brahim-driouch/envbox.git/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	err := config.ConnectDB()

	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	fmt.Println("Connected to the database successfully")

	defer config.Pool.Close()

}
