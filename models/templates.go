package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Template struct {
	Id*	int
	Name	string	 
}

func GetTemplates() ([]Template, error) {
	requestUrl := MeritMockApiBaseURL + "/templates"
	response, err := http.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var templates []Template
	err = json.Unmarshal(body, &templates)

	if err != nil {
		return nil, err
	}

	return templates, err
}

func GetMeritsByTemplate(templateId string) ([]Merit, error) {
	requestUrl := MeritMockApiBaseURL + "/templates/" + templateId + "/merits"
	response, err := http.Get(requestUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	var merits []Merit
	err = json.Unmarshal(body, &merits)

	if err != nil {
		return nil, err
	}

	return merits, err
}