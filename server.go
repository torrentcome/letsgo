package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	//"net/http"
)

func startServer(db *sql.DB) {
	router := gin.Default()

	router.GET("/stop_id/:stop_id", func(c *gin.Context) {
		stopId := c.Param("stop_id")
		c.JSON(200, gin.H{
			"code":    "PAGE_FOUND",
			"message": "data " + stopId})

	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found"})
	})

	router.Run(":8080")
}
