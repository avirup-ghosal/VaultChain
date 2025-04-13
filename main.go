package main

import (
	"fmt"
	"log"

	"github.com/avirup-ghosal/VaultChain/core"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var option int
	fmt.Println("Option 1: Validate BlockChain")
	fmt.Println("Option 2: Add data")
	fmt.Println("Enter option:")
	fmt.Scanln(&option)
	if option == 1 {
		db, err := core.Dbinit()
		if err != nil {
			log.Fatalf("error initializing db: %w", err)
		}
		err = core.ValidateBlockchain("logs.txt", db)

		if err != nil {
			log.Fatalf("validation error: %w", err)
		}
	}

	if option == 2 {
		var d string
		fmt.Println("Enter data to be stored:")
		fmt.Scanln(&d)
		newblockchain := core.NewBlockchain()
		err := newblockchain.AddBlock(d)
		if err != nil {
			log.Fatalf("error adding block: %w", err)
		}
	}
}
