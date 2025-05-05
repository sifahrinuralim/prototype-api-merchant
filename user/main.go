package main

import (
	"log"
	"net/http"
	"user/controller"
)

func main() {
	http.HandleFunc("/v1/user/pin-validation", controller.PinValidation)
	http.HandleFunc("/v1/user/login", controller.Login)
	http.HandleFunc("/v1/user/logout", controller.Logout)
	http.HandleFunc("/v1/account/info", controller.AccountInfo)
	http.HandleFunc("/v1/account/main-transaction", controller.MainTransaction)

	log.Println("Server started at :8090")
	http.ListenAndServe(":8090", nil)
}
