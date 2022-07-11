package main

import (
	"apiMore/handlers"
	"apiMore/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a db connection (db path is in dataConn.go - wants to be in properties file)
	err := models.ConnectDb()
	models.CheckErr(err)

	router := gin.Default()
	
	// Create router mappings - ideally, put these in a separate package(s)

	// Events
	router.GET("events", handlers.GetEvents)
	router.GET("events/:id", handlers.GetEventById)
	router.POST("events", handlers.AddEvent)
	router.PUT("events/:id", handlers.UpdateEvent)
	router.DELETE("events/:id", handlers.DeleteEvent)
	router.OPTIONS("events", handlers.EventOptions)

	// Checkins
	router.GET("checkins", handlers.GetCheckins)
	router.GET("checkins/:id", handlers.GetCheckinById)
	router.POST("checkins", handlers.AddCheckin)
	// Despite the URL, this is creeating a checkin, so group it with checkins
	router.POST("events/:eventId/user/:userId", handlers.AddCheckinForUserEvent) 
	router.PUT("checkins/:id", handlers.UpdateCheckin)
	router.DELETE("checkins/:id", handlers.DeleteCheckin)
	router.OPTIONS("checkins", handlers.CheckinsOptions)

	// Users 
	router.GET("users/:userId/merits", handlers.GetMeritsByUser)
	router.POST("users/:userId/merits", handlers.AddMerit)

	// Templates
	router.GET("templates", handlers.GetTemplates)
	router.GET("templates/:templateId/merits", handlers.GetMeritsByTemplate)

	router.Run("localhost:8080")
}
