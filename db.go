package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func startDb() *sql.DB {
	os.Remove("data/gtfs.db")
	log.Println(" => Remove then Creating ./data/gtfs.db")
	file, err := os.Create("data/gtfs.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	log.Println("./data/gtfs.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./data/gtfs.db")

	configure(sqliteDatabase)
	createTableRoute(sqliteDatabase)
	createTableStopTime(sqliteDatabase)
	return sqliteDatabase
}

func configure(db *sql.DB) {
	ddl := `
PRAGMA automatic_index = ON;
	PRAGMA cache_size = 32768;
	PRAGMA cache_spill = OFF;
	PRAGMA foreign_keys = ON;
	PRAGMA journal_size_limit = 67110000;
	PRAGMA locking_mode = NORMAL;
	PRAGMA page_size = 4096;
	PRAGMA recursive_triggers = ON;
	PRAGMA secure_delete = ON;
	PRAGMA synchronous = NORMAL;
	PRAGMA temp_store = MEMORY;
	PRAGMA journal_mode = WAL;
	PRAGMA wal_autocheckpoint = 16384;
	`
	db.Exec(ddl)
	log.Println(" => Configure DB")
}

func createTableRoute(db *sql.DB) {
	createTable := `
CREATE TABLE TABLE_ROUTE (
_ID INTEGER PRIMARY KEY,		
		COLUMN_ROUTE_ID TEXT,
		COLUMN_ROUTE_SHORT_NAME TEXT
	  );`

	log.Println(" => Creating TABLE_ROUTE...")
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("TABLE_ROUTE created")
}

func createTableStopTime(db *sql.DB) {
	createTable := `
CREATE TABLE TABLE_STOP_TIME (
_ID INTEGER PRIMARY KEY,		
		COLUMN_TRIP_ID TEXT,
		COLUMN_ROUTE_ID TEXT,
		COLUMN_ARRIVAL_TIME TEXT,
		COLUMN_DEPARTURE_TIME TEXT,
		COLUMN_STOP_ID TEXT,
		COLUMN_STOP_HEADSIGN TEXT
	);`

	log.Println(" => Creating TABLE_STOP_TIME...")
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("TABLE_STOP_TIME created")
}

func insertRoute(db *sql.DB, routeID string, routeShortName string) {
	insertSQL := `INSERT INTO TABLE_ROUTE(COLUMN_ROUTE_ID, COLUMN_ROUTE_SHORT_NAME) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(&routeID, &routeShortName)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func insertStopTime(db *sql.DB, tripID string, routeID string, arrivalTime string, departureTime string, stopID string, stopHeadsign string) {
	insertSQL := `INSERT INTO TABLE_STOP_TIME(COLUMN_TRIP_ID, COLUMN_ROUTE_ID, COLUMN_ARRIVAL_TIME, COLUMN_DEPARTURE_TIME, COLUMN_STOP_ID, COLUMN_STOP_HEADSIGN) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(&tripID, &routeID, &arrivalTime, &departureTime, &stopID, &stopHeadsign)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
