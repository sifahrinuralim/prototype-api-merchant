package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"history/repository"
	"history/util"
	"io"
	"net/http"
)

func HistoryLog(res http.ResponseWriter, req *http.Request) error {

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("cannot read body")
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode setelah ambil raw
	var input map[string]string
	errRequest := json.Unmarshal(bodyBytes, &input)
	if errRequest != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	action := input["action"]
	description := input["description"]

	repository.LogHistory(action, description)

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
	})

	return nil
}
