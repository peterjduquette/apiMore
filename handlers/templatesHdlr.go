package handlers

import (
	"apiMore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTemplates (context *gin.Context) {
	templates, err := models.GetTemplates()

	models.CheckErr(err)

	if templates == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": templates})
	}
}

func GetMeritsByTemplate(context *gin.Context) {
	merits, err := models.GetMeritsByTemplate(context.Param("templateId"))

	models.CheckErr(err)

	if merits == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": merits})
	}
}