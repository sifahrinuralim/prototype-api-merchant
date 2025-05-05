package controller

import (
	"history/service"
	"log"
	"net/http"
)

func HistoryLog(res http.ResponseWriter, req *http.Request) {
	err := service.HistoryLog(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}
