package main

import (
	"history/controller"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/history/log", controller.HistoryLog)

	log.Println("Server started at :8200")
	http.ListenAndServe(":8200", nil)
}
