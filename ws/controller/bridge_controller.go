package controller

import (
	"net/http"
	"ws/service"
)

func MainWs(res http.ResponseWriter, req *http.Request) {

	service.Bridge(res, req)

	res.WriteHeader(http.StatusOK)

}
