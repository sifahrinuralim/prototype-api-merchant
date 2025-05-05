package main

import (
	"log"
	"net/http"
	"ws/controller"
)

func main() {
	http.HandleFunc("/v1/ws", controller.MainWs)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
