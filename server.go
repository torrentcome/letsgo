package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func startServer(db *sql.DB) {
	router := gin.Default()

	router.GET("/stop_id/:stop_id", func(c *gin.Context) {
		stopID := c.Param("stop_id")
		rows, err := db.Query("SELECT * FROM TABLE_STOP_TIME INNER JOIN TABLE_ROUTE ON TABLE_ROUTE.COLUMN_ROUTE_ID = TABLE_STOP_TIME.COLUMN_ROUTE_ID WHERE COLUMN_STOP_ID=?", stopID)
		return500(c, err)

		type entry struct {
			tripID         string
			routeID        string
			arrivalTime    string
			departureTime  string
			stopID         string
			stopHeadsign   string
			routeShortName string
		}

		var array []entry

		var tripID string
		var routeID string
		var arrivalTime string
		var departureTime string
		var dbStopID string
		var stopHeadsign string
		var routeShortName string

		if rows.Next() {
			err = rows.Scan(&tripID, &routeID, &departureTime, &arrivalTime, &dbStopID, &stopHeadsign, &routeShortName)
			return500(c, err)
			array = append(array, entry{tripID, routeID, departureTime, arrivalTime, dbStopID, stopHeadsign, routeShortName})
		}
		defer rows.Close()
		if len(array) <= 0 {
			c.JSON(http.StatusNoContent, gin.H{"code": http.StatusNoContent, "message": "No content for this stop_id"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": array, "count": len(array)})
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
			"message": "Internal Server Error"})
	}
}
