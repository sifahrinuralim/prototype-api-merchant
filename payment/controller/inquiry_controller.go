package controller

import (
	"log"
	"net/http"
	"payment/service"
)

func Inquiry(res http.ResponseWriter, req *http.Request) {
	err := service.Inquiry(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}
