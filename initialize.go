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
		return data.contest, data.config, nil
	} else {
		return Contest{}, MailConfig{}, errors.New("cannot read data.json file")
	}
}
