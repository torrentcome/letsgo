package main

import (
	"bufio"
	"database/sql"
	"github.com/cheggaaa/pb/v3"
	"log"
	"os"
	"strings"
	// "sync"
)

func startParse(db *sql.DB) {
	parseRoute(db)
	parseStopTime(db)
}

func parseRoute(db *sql.DB) {
	file, err := os.Open("./data/routes.txt")
	count, _ := lineCount("./data/routes.txt")

	bar := pb.Simple.Start(count)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		phrase := scanner.Text()
		phraseSplit := strings.Split(phrase, ",")
		insertRoute(db, phraseSplit[0], phraseSplit[2])
		bar.Increment()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	bar.Finish()
}

func parseStopTime(db *sql.DB) {
	file, err := os.Open("./data/stop_times.txt")
	count, _ := lineCount("./data/stop_times.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bar := pb.Simple.Start(count)
	scanner := bufio.NewScanner(file)

	// phraseC := make(chan string)
	// wg := sync.WaitGroup{}

	for scanner.Scan() {
		phrase := scanner.Text()
		lineSplit := strings.Split(phrase, ",")

		tripID := lineSplit[0]
		parseTripID := strings.Split(tripID, ".")
		routeID := ""
		if len(parseTripID) > 2 {
			routeID = parseTripID[2]
		}

		insertStopTime(db, lineSplit[0], routeID, lineSplit[1], lineSplit[2], lineSplit[3], lineSplit[5])
		bar.Increment()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	bar.Finish()
}
