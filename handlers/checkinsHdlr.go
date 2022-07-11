package handlers

import (
	"apiMore/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Respond with all checkins
func GetCheckins(context *gin.Context) {
	checkins, err := models.GetCheckins()

	models.CheckErr(err)

	if checkins == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": checkins})
	}
}

// Respond with one checkin, base on the id passed in
func GetCheckinById(context *gin.Context) {
	id := context.Param("id")

	checkinId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
		}

	checkin, err := models.GeCheckinById(checkinId)

	models.CheckErr(err)

	// If checkin id is zero, no rows were found 
	// ASSUMPTION: in the db, checkin.id will never be zero - starts at 1. 
	// (look into a better way to check for no rows returned)
	if checkin.Id == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": checkin})
	}
}

// Create a new checkin, based on the request body
func AddCheckin(context *gin.Context) {

	var json models.Checkin

	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddCheckin(json)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// Add a checkin, for a given user-id and event-id, if and only if the user has the required merit for the event
// NOTE: should this be in the Events model instead of in checkins, since the api url starts with 'events'? 
func AddCheckinForUserEvent(context *gin.Context) {
	eventId := context.Param("eventId")
	userId := context.Param("userId")

	// Call Merit API to get merits for template id, user id
	requestUrl := models.MeritMockApiBaseURL + "/templates/" + eventId + "/merits?userId=" + userId
	response, err := http.Get(requestUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var merits []models.Merit
	err = json.Unmarshal(body, &merits)
	if err != nil {
		log.Fatal(err)
	}

	// If user has the merit for this event, add a new checkin
	// ASSUMPTION: Only one merit returned for each event id (or we would need to loop through merits, 
	// and add checkins for each)
	if len(merits) > 0 {
		intEventId, err := strconv.Atoi(eventId)
		if err != nil {
			log.Fatal(err)
		}

		checkinToAdd := models.Checkin{
			EventId:     intEventId,
			MeritUserId: userId,
		}

		success, err := models.AddCheckin(checkinToAdd)

		if success {
			context.JSON(http.StatusOK, gin.H{"message": "Success"})
		} else {
			context.JSON(http.StatusBadRequest, gin.H{"error": err})
		}
	}
}

// Update a checkin, with new values based on the request body
func UpdateCheckin(context *gin.Context) {

	var json models.Checkin

	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkinId, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateCheckin(json, checkinId)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// Delete checkin record based on the checkin id passed in
func DeleteCheckin(context *gin.Context) {

	eventId, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteCheckin(eventId)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// Return api options for checkins
func CheckinsOptions(context *gin.Context) {

	options := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	context.String(200, options)
}

/*

// Example checkins POST request body
// NOTE: id value is ignored, as endpoint always increments id from max id value in checkins table
{
    "id": 0,
    "eventId": 123,
    "beganAt": "2022-07-05T00:00:00Z",
    "completedAt": null,
    "mertitUserId": "555"
}

// Example checkins PUT request body
// NOTE id value in body is ignored - id to target update record in checkins table is nabbed from request URL. 
{
    "id": 0,
    "eventId": 321,
    "beganAt": "2022-07-07T00:00:00Z",
    "completedAt": null,
    "mertitUserId": "updated"
}
*/