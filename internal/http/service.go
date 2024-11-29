package http

import (
	"github.com/alexohneander/GoZilla/internal/handler"
	"github.com/gin-gonic/gin"
)

func Server() {
	// Define HTTP Service
	router := gin.Default()

	router.GET("/announce", handler.Announce)
	router.GET("/scrape", handler.Scrape)

	router.Run(":4000") // listen and serve on 0.0.0.0:4000
}
