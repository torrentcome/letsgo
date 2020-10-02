package main

import (
	"log"
)

func main() {

	log.Println("--- Database ---")
	db := startDb()

	log.Println("--- Parse ---")
	startParse(db)

	log.Println("--- Server ---")
	startServer(db)
}
