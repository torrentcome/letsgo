package main

import (
	"log"
)

func main() {
	log.Println("1. Database")
	db := startDb()
	log.Println("2. Parse")
	parse(db)
}
