package http

import (
	"AvitoWinter/internal/service"
	"encoding/json"
	log2 "github.com/sirupsen/logrus"
	"net/http"
)

type UserServer struct {
	service *service.UserService
}

var _ ServerInterface = (*UserServer)(nil)

func NewUserServer(service *service.UserService) *UserServer {
	return &UserServer{service: service}
}

type Error struct {
	Code    int32
	Message string
}

func sendErrorResponse(w http.ResponseWriter, code int, resp ErrorResponse) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log2.Infof("sendErrorResponse: ошибка при формировании ответа ошибки %s: %s", resp, err.Error())
	}
}
