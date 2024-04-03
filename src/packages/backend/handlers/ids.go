package handlers

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostAlgorithmHandler(c *gin.Context) {
	var reqBody models.AlgorithmRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data diterima",
		"data":    reqBody,
	})
}
