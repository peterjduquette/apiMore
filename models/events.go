package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	Id          int						`json:"id"`
	GrantedMeritTemplateId string	`json:"granted_merit_template_id"`
    Name string							`json:"name"`
    QualifyingMeritTemplateId string	`json:"qualifying_merit_template_id"`
}

// Query for all events in the events table
func GetEvents() ([]Event, error) {

	rows, err := DbConn.Query("SELECT id, granted_merit_template_id, name, qualifying_merit_template_id from events")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := make([]Event, 0)

	for rows.Next() {
		event := Event{}
		err = rows.Scan(&event.Id, &event.GrantedMeritTemplateId, &event.Name, &event.QualifyingMeritTemplateId)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return events, err
}

// Query the event with the id passed in
func GeEventById(id int) (Event, error) {

	stmt, err := DbConn.Prepare("SELECT id, granted_merit_template_id, name, qualifying_merit_template_id from events WHERE id = ?")

	if err != nil {
		return Event{}, err
	}

	event := Event{}

	sqlErr := stmt.QueryRow(id).Scan(&event.Id, &event.GrantedMeritTemplateId, &event.Name, &event.QualifyingMeritTemplateId)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Event{}, nil
		}
		return Event{}, sqlErr
	}
	return event, nil
}

// Create a new event, based on the Event struct passed in
func AddEvent(newEvent Event) (bool, error) {
	transaction, err := DbConn.Begin()
	if err != nil {
		return false, err
	}
	
	stmt, err := transaction.Prepare(`INSERT INTO events (id, granted_merit_template_id, name, qualifying_merit_template_id) 
	VALUES ((select max(id) from events) + 1, ?, ?, ?)`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newEvent.GrantedMeritTemplateId, newEvent.Name, newEvent.QualifyingMeritTemplateId)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

// Update the event record in the events table, based on the Event struct passed in
func UpdateEvent(updEvent Event, eventId int) (bool, error) {

	transaction, err := DbConn.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := transaction.Prepare("UPDATE events SET granted_merit_template_id = ?, name = ?, qualifying_merit_template_id = ? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updEvent.GrantedMeritTemplateId, updEvent.Name, updEvent.QualifyingMeritTemplateId, eventId)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

// Delete the event record in the events table, based on the Event struct passed in
func DeleteEvent(eventId int) (bool, error) {

	transaction, err := DbConn.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DbConn.Prepare("DELETE from events where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(eventId)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}