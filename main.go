package main

import (
	"fmt"
	"log"
)

func main() {

	var option int
	fmt.Println("Option 1: Validate BlockChain")
	fmt.Println("Option 2: Add data")
	fmt.Println("Enter option:")
	fmt.Scanln(&option)
	if option == 1 {
		db, err := Dbinit()
		if err != nil {
			log.Fatalf("error initializing db: %w", err)
		}
		err = validateBlockchain("logs.txt", db)
		if err != nil {
			log.Fatalf("validation error: %w", err)
		}
	}

	if option == 2 {
		var d string
		fmt.Println("Enter data to be stored:")
		fmt.Scanln(&d)
		newblockchain := NewBlockchain()
		err := newblockchain.AddBlock(d)
		if err != nil {
			log.Fatalf("error adding block: %w", err)
		}
	}
}
