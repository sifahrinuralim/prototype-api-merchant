package controller

import (
	"log"
	"net/http"
	"user/service"
)

func AccountInfo(res http.ResponseWriter, req *http.Request) {
	err := service.AccountInfo(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}

func MainTransaction(res http.ResponseWriter, req *http.Request) {
	err := service.MainTransaction(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}
