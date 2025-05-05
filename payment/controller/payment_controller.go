package controller

import (
	"log"
	"net/http"
	"payment/service"
)

func Payment(res http.ResponseWriter, req *http.Request) {
	err := service.Payment(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}
