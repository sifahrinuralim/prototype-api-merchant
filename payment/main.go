package main

import (
	"log"
	"net/http"
	"payment/controller"
)

func main() {
	http.HandleFunc("/v1/payment/inquiry", controller.Inquiry)
	http.HandleFunc("/v1/payment/payment", controller.Payment)

	log.Println("Server started at :8100")
	http.ListenAndServe(":8100", nil)
}
