package main

import (
	"fmt"
	"log"
	"os"
)

func writefile(filepath string, sentence []byte) {
	d1 := sentence
	err := os.WriteFile(filepath, d1, 0644)
	check(err)

}

func appendfile(filepath string, sentence []byte) {
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(sentence); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
func readfile(filepath string) {
	f, err := os.ReadFile(filepath)
	check(err)
	fmt.Println(string(f))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
