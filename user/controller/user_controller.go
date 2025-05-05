package controller

import (
	"log"
	"net/http"
	"user/service"
)

func PinValidation(res http.ResponseWriter, req *http.Request) {
	err := service.PinValidation(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	err := service.Login(res, req)
	if err != nil {
		log.Printf("Login error: %v", err)
	}
}

func Logout(res http.ResponseWriter, req *http.Request) {
	err := service.Logout(res, req)
	if err != nil {
		log.Printf("Logout error: %v", err)
	}
}
