package models

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Checkin struct {
	Id          int			`json:"id"`
	EventId     int			`json:"eventId"`
	BeganAt     *time.Time	`json:"beganAt"`
	CompletedAt *time.Time	`json:"completedAt"`
	MeritUserId string		`json:"mertitUserId"`
}

func GetCheckins() ([]Checkin, error) {

	rows, err := DbConn.Query("SELECT id, event_id, began_at, completed_at, merit_user_id from checkins")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	checkins := make([]Checkin, 0)

	for rows.Next() {
		checkin := Checkin{}
		err = rows.Scan(&checkin.Id, &checkin.EventId, &checkin.BeganAt, &checkin.CompletedAt, &checkin.MeritUserId)

		if err != nil {
			return nil, err
		}

		checkins = append(checkins, checkin)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return checkins, err
}

func GeCheckinById(id int) (Checkin, error) {

	stmt, err := DbConn.Prepare("SELECT id, event_id, began_at, completed_at, merit_user_id from checkins WHERE id = ?")

	if err != nil {
		return Checkin{}, err
	}

	checkin := Checkin{}

	sqlErr := stmt.QueryRow(id).Scan(&checkin.Id, &checkin.EventId, &checkin.BeganAt, &checkin.CompletedAt, &checkin.MeritUserId)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Checkin{}, nil
		}
		return Checkin{}, sqlErr
	}
	return checkin, nil
}

func AddCheckin(newCheckin Checkin) (bool, error) {
	transaction, err := DbConn.Begin()
	if err != nil {
		return false, err
	}
	
	stmt, err := transaction.Prepare(`INSERT INTO checkins (id, event_id, began_at, completed_at, merit_user_id) 
	VALUES ((select max(id) from checkins) + 1, ?, date('now'), NULL, ?)`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newCheckin.EventId, newCheckin.MeritUserId)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}

func UpdateCheckin(updCheckin Checkin, id int) (bool, error) {

	transaction, err := DbConn.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := transaction.Prepare("UPDATE checkins SET event_id = ?, began_at = ?, completed_at = ?, merit_user_id =? WHERE id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(updCheckin.EventId, updCheckin.BeganAt, updCheckin.CompletedAt, updCheckin.MeritUserId, id)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}


func DeleteCheckin(checkinId int) (bool, error) {

	transaction, err := DbConn.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DbConn.Prepare("DELETE from checkins where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(checkinId)

	if err != nil {
		return false, err
	}

	transaction.Commit()

	return true, nil
}