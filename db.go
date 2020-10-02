package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func startDb() *sql.DB {
	os.Remove("data/gtfs.db")
	log.Println("Remove then Creating ./data/gtfs.db")
	file, err := os.Create("data/gtfs.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("./data/gtfs.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./data/gtfs.db")
	// defer sqliteDatabase.Close()

	createTableRoute(sqliteDatabase)
	createTableStopTime(sqliteDatabase)
	return sqliteDatabase
}

func createTableRoute(db *sql.DB) {
	createTable := `CREATE TABLE TABLE_ROUTE (
		_ID INTEGER PRIMARY KEY,		
		COLUMN_ROUTE_ID TEXT,
		COLUMN_ROUTE_SHORT_NAME TEXT
	  );`

	log.Println("creating TABLE_ROUTE...")
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("TABLE_ROUTE created")
}

func createTableStopTime(db *sql.DB) {
	createTable := `CREATE TABLE TABLE_STOP_TIME (
		_ID INTEGER PRIMARY KEY,		
		COLUMN_TRIP_ID TEXT,
                COLUMN_ROUTE_ID TEXT,
                COLUMN_ARRIVAL_TIME TEXT,
                COLUMN_DEPARTURE_TIME TEXT,
                COLUMN_STOP_ID TEXT,
                COLUMN_STOP_HEADSIGN TEXT
	);`

	log.Println("Creating TABLE_STOP_TIME...")
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("TABLE_STOP_TIME created")
}

func insertRoute(db *sql.DB, routeId string, routeShortName string) {
	//	log.Println("Inserting route...")
	insertSQL := `INSERT INTO TABLE_ROUTE(COLUMN_ROUTE_ID, COLUMN_ROUTE_SHORT_NAME) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(&routeId, &routeShortName)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func insertStopTime(db *sql.DB, tripId string, routeId string, arrivalTime string, departureTime string, stopId string, stopHeadsign string) {
	// log.Println("Inserting registry...")
	insertSQL := `INSERT INTO TABLE_STOP_TIME(COLUMN_TRIP_ID, COLUMN_ROUTE_ID, COLUMN_ARRIVAL_TIME, COLUMN_DEPARTURE_TIME, COLUMN_STOP_ID, COLUMN_STOP_HEADSIGN) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(tripId, routeId, arrivalTime, departureTime, stopId, stopHeadsign)
	if err != nil {
		log.Fatalln(err.Error())
	}
}