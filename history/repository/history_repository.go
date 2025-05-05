package repository

import (
	"encoding/json"
	"history/model"
	"io/ioutil"
	"os"
)

func LogHistory(action, description string) {
	file := "storage/history.json"
	data, _ := ioutil.ReadFile(file)
	var history []model.History
	json.Unmarshal(data, &history)

	history = append(history, model.History{Action: action, Description: description})
	newData, _ := json.MarshalIndent(history, "", "  ")
	ioutil.WriteFile(file, newData, os.ModePerm)
}
