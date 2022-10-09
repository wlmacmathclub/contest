package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func initialize() (Contest, MailConfig, error) {
	_, err := os.Stat("cache")
	if err != nil && os.IsNotExist(err) {
		os.Mkdir("cache", 0700)
	}
	_, err2 := os.Stat("contest.log")
	if err2 == nil {
		os.Remove("contest.log")
	}

	jsonFile, err3 := os.Open("data.json")
	if err3 == nil {
		defer jsonFile.Close()
		jsonBytes, _ := ioutil.ReadAll(jsonFile)
		data := JsonData{}
		json.Unmarshal(jsonBytes, &data)
		return data.Contest, data.Config, nil
	} else {
		os.Create("data.json")
		return Contest{}, MailConfig{}, errors.New("cannot read data.json file")
	}
}

func saveJson(contest Contest, config MailConfig) {
	data, _ := json.MarshalIndent(JsonData{
		Contest: contest,
		Config:  config,
	}, "", "\t")
	os.WriteFile("data.json", data, 0644)
}
