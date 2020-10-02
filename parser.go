package main

import (
	"bufio"
	"database/sql"
	"github.com/cheggaaa/pb/v3"
	"log"
	"os"
	"strings"
)

func startParse(db *sql.DB) {
	parseRoute(db)
	//	parseStopTime(db)
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

		routeId := phraseSplit[0]
		routeShortName := phraseSplit[2]

		insertRoute(db, routeId, routeShortName)
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

	bar := pb.Simple.Start(count)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		phrase := scanner.Text()
		lineSplit := strings.Split(phrase, ",")

		tripId := lineSplit[0]
		parseTripId := strings.Split(tripId, ".")
		routeId := ""
		if len(parseTripId) > 2 {
			routeId = parseTripId[2]
		}

		arrivalTime := lineSplit[1]
		departureTime := lineSplit[2]
		stopId := lineSplit[3]
		stopHeadsign := lineSplit[5]

		insertStopTime(db, tripId, routeId, arrivalTime, departureTime, stopId, stopHeadsign)
		bar.Increment()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	bar.Finish()
}
