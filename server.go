package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func startServer(db *sql.DB) {
	router := gin.Default()

	router.GET("/stop_id/:stop_id", func(c *gin.Context) {
		stopID := c.Param("stop_id")
		fmt.Println("stop_id =" + stopID)
		rows, err := db.Query(`SELECT COLUMN_TRIP_ID, TABLE_ROUTE.COLUMN_ROUTE_ID, COLUMN_ARRIVAL_TIME, COLUMN_DEPARTURE_TIME, COLUMN_STOP_ID, COLUMN_STOP_HEADSIGN, COLUMN_ROUTE_SHORT_NAME FROM TABLE_STOP_TIME INNER JOIN TABLE_ROUTE ON TABLE_ROUTE.COLUMN_ROUTE_ID = TABLE_STOP_TIME.COLUMN_ROUTE_ID WHERE COLUMN_STOP_ID=?`, stopID)
		if err != nil {
			return500(c, err)
		}
		type entry struct {
			tripID         string
			routeID        string
			arrivalTime    string
			departureTime  string
			stopID         string
			stopHeadsign   string
			routeShortName string
		}

		fmt.Println(rows)

		var array []entry

		for rows.Next() {
			e := entry{}
			err = rows.Scan(&e.tripID, &e.routeID, &e.departureTime, &e.arrivalTime, &e.stopID, &e.stopHeadsign, &e.routeShortName)
			fmt.Println(e)
			return500(c, err)
			array = append(array, e)
			fmt.Println(array)
		}
		defer rows.Close()
		if len(array) <= 0 {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "No content for stop_id = " + stopID})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": fmt.Sprint(array), "count": len(array)})
		}
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Page not found"})
	})

	router.Run(":8080")
}

func return500(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error = " + err.Error()})
	}
}
