package models

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Merit struct {
	Id         int
	TemplateId *int
	UserId     *int
}


func GetMeritsByUser (userId string) ([]Merit, error) {
	requestUrl := MeritMockApiBaseURL + "/templates/" + userId + "/merits"
	response, err := http.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var merits []Merit
	err = json.Unmarshal(body, &merits)

	if err != nil {
		return nil, err
	}

	return merits, err
}

func AddMerit(body Merit) (bool, error) {
	requestUrl := MeritMockApiBaseURL + "/users/" + strconv.Itoa(*body.UserId) + "/merits"

	jsonVal, _ := json.Marshal(body)
	_, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(jsonVal))

	if err != nil {
		return false, err
	}

	return true, err
}