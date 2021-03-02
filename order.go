package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func submitOrder(c *gin.Context) {

	c.String(http.StatusOK, "OK")
}

func queryOrder(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
