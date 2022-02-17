package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	db "pki-server/models"
)

func Healthz(c *gin.Context) {
	if err := db.Healthz(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "db connection exception"})
		return
	}
	c.Status(http.StatusOK)
	return
}
