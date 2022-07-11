package handlers

import (
	//"apiMore/main"
	"apiMore/models"
	"fmt"
	//"encoding/json"
	//"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Respond with all events
func GetEvents(context *gin.Context) {
	events, err := models.GetEvents()

	models.CheckErr(err)
	/*
	if err != nil {
		log.Fatal(err)
	}
	*/

	if events == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": events})
	}
}

// Respond with one event, base on the id passed in
func GetEventById(context *gin.Context) {
	id := context.Param("id")

	eventId, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
		}

	event, err := models.GeEventById(eventId)

	models.CheckErr(err)

	// If event id is zero, no row were founf 
	// ASSUMPTION: in the db, event.id will never be zero - starts at 1. 
	// (look into a better way to check for no rows returned)
	if event.Id == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"data": event})
	}
}

// // Create a new event, based on the request body
func AddEvent(context *gin.Context) {

	var json models.Event

	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddEvent(json)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

/*
// Add a checkin, for a give user-id and event-id, if and only if the user has the required merit for the event
// NOTE: should this be in the Events model instead of in checkins, since the api url starts with 'events'? 
func AddCheckinForUserEvent(context *gin.Context) {
	eventId := context.Param("eventId")
	userId := context.Param("userId")

	// Call Merit API to get merits for template id, user id
	// NOTE: base url for (mock) merit api hardcoede for now - should be exernalized (properties file)
	// NOTE: calls to (mock) merti api should probably be encapsulated in their own package
	requestUrl := "http://localhost:3000/templates/" + eventId + "/merits?userId=" + userId
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
*/

// Update a event, with new values based on the request body
func UpdateEvent(context *gin.Context) {

	var json models.Event

	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// DEBUG
	fmt.Println(">>>>>>>>>>>>>> NAME: " + json.Name)

	eventId, err := strconv.Atoi(context.Param("id"))

	fmt.Printf("Updating id %d", eventId)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.UpdateEvent(json, eventId)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// Delete event record based on the checkin id passed in
func DeleteEvent(context *gin.Context) {

	eventId, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteEvent(eventId)

	if success {
		context.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

// Return api options for events
func EventOptions(context *gin.Context) {

	options := "HTTP/1.1 200 OK\n" +
		"Allow: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Origin: http://locahost:8080\n" +
		"Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS\n" +
		"Access-Control-Allow-Headers: Content-Type\n"

	context.String(200, options)
}

/*
POST example request body (NOTE: id value is ignored, as endpoint always increments from max id in events table): 

{
        "id": 0,
        "granted_merit_template_id": "123",
        "name": "API Test POST",
        "qualifying_merit_template_id": "333"
}

PUT example request body

{
        "id": 2,
        "granted_merit_template_id": "321",
        "name": "API Test PUT",
        "qualifying_merit_template_id": "999"
}
*/