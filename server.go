package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func startServer(db *sql.DB) {
	router := gin.Default()

	router.GET("/stop_id/:stop_id", func(c *gin.Context) {
		stopID := c.Param("stop_id")
		rows, err := db.Query("SELECT * FROM TABLE_STOP_TIME WHERE COLUMN_STOP_ID=? INNER JOIN TABLE_ROUTE ON TABLE_ROUTE.COLUMN_ROUTE_ID = TABLE_STOP_TIME.COLUMN_ROUTE_ID", stopID)
		return500(c, err)

		var tripID string
		var routeID string
		var arrivalTime string
		var departureTime string
		var dbStopID string
		var stopHeadsign string

		if rows.Next() {
			err = rows.Scan(&tripID, &routeID, &departureTime, &arrivalTime, &dbStopID, &stopHeadsign)
			return500(c, err)

		}
		defer rows.Close()
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found"})
	})
	router.Run(":8080")
}

func return500(c *gin.Context, err error) {
	if err != nil {
		c.JSON(500, gin.H{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal Server Error"})
	}
}
