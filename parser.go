package main

import (
	"bufio"
	"database/sql"
	"github.com/cheggaaa/pb/v3"
	"log"
	"os"
	"strings"
	"sync"
)

func startParse(db *sql.DB) {
	log.Println("Parsing --- routes file ...")
	parseRoute(db)
	updateRoute(db)
	log.Println("Parsing --- stop_time file ...")
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
		insertRoute(db, trimQuotes(phraseSplit[0]), trimQuotes(phraseSplit[2]))
		bar.Increment()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	bar.Finish()
}

func updateRoute(db *sql.DB) {
	smt, err := db.Prepare(`UPDATE TABLE_ROUTE SET COLUMN_ROUTE_ID = REPLACE(COLUMN_ROUTE_ID,'"','');`)
	if err != nil {
		log.Fatal(err)
	}
	smt.Exec()
}

func parseStopTime(db *sql.DB) {
	count, _ := lineCount("./data/stop_times.txt")
	file, err := os.Open("./data/stop_times.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bar := pb.Full.Start(count)
	scanner := bufio.NewScanner(file)

	type entry struct {
		tripID        string
		routeID       string
		arrivalTime   string
		departureTime string
		stopID        string
		stopHeadsign  string
		wg            *sync.WaitGroup
	}
	entries := make(chan entry)
	wg := sync.WaitGroup{}

	insertSQL := `INSERT INTO TABLE_STOP_TIME(COLUMN_TRIP_ID, COLUMN_ROUTE_ID, COLUMN_ARRIVAL_TIME, COLUMN_DEPARTURE_TIME, COLUMN_STOP_ID, COLUMN_STOP_HEADSIGN) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer statement.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	countEntry := 0

	go func() {
		for {
			select {
			case entry, ok := <-entries:
				if ok {
					_, err := tx.Stmt(statement).Exec(entry.tripID, entry.routeID, entry.arrivalTime, entry.departureTime, entry.stopID, entry.stopHeadsign)
					if err != nil {
						log.Fatal(err)
					}

					if countEntry%32768 == 0 {
						tx.Commit()
						db.Exec(`PRAGMA shrink_memory;`)

						tx, err = db.Begin()
						if err != nil {
							log.Fatal(err)
						}
					}

					countEntry++
					bar.Increment()
					entry.wg.Done()
				}
			}
		}
	}()

	linesChunkLen := 64 * 1024
	lines := make([]string, 0, 0)

	for scanner.Scan() {
		phrase := scanner.Text()
		lines = append(lines, phrase)

		if len(lines) == linesChunkLen {
			wg.Add(len(lines))
			process := lines
			go func() {
				for _, text := range process {
					lineSplit := strings.Split(text, ",")
					tripID := lineSplit[0]
					parseTripID := strings.Split(tripID, ".")
					routeID := ""
					if len(parseTripID) > 2 {
						routeID = parseTripID[2]
					}

					e := entry{wg: &wg}
					e.tripID = trimQuotes(tripID)
					e.routeID = trimQuotes(routeID)
					e.arrivalTime = trimQuotes(lineSplit[1])
					e.departureTime = trimQuotes(lineSplit[2])
					e.stopID = trimQuotes(lineSplit[3])
					e.stopHeadsign = trimQuotes(lineSplit[5])
					entries <- e
				}
			}()
			lines = make([]string, 0, linesChunkLen)
		}
	}
	wg.Wait()
	tx.Commit()
	close(entries)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	bar.Finish()
}
