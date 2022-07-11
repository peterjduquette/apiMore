package handlers

import (
	"apiMore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMeritsByUser (context *gin.Context) {
	userId := context.Param("userId")

	merits, err := models.GetMeritsByUser(userId)

	models.CheckErr(err)

	if merits == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": merits})
	}
}

func AddMerit(context *gin.Context) {
	var json models.Merit

	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddMerit(json)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}