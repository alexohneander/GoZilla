package handler

import "github.com/gin-gonic/gin"

func Scrape(c *gin.Context) {
	c.String(200, "Scrape")
}
