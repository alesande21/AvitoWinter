package http

import (
	"AvitoWinter/internal/entity"
	"AvitoWinter/internal/service"
	"encoding/json"
	log2 "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type UserServer struct {
	service *service.ShopService
}

var _ ServerInterface = (*UserServer)(nil)

func NewUserServer(service *service.ShopService) *UserServer {
	return &UserServer{service: service}
}

func (u UserServer) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (u UserServer) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	purchaseInfo, err := entity.NewPurchaseInfo(item, item)
	if err != nil {
		var errorDescription = "Не задан предмет покупки."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	err = u.service.PurchaseItem(r.Context(), *purchaseInfo)
	if err != nil {
		var errorDecription string = "Какая то там ошибка"
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDecription})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var sb strings.Builder
	sb.WriteString("Покупка + ")
	if err := json.NewEncoder(w).Encode(sb.String()); err != nil {
		log2.Errorf("CreateUser-> json.NewEncoder: ошибка при кодирования овета: %s", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, ErrorResponse{Reason: "Ошибка кодирования ответа."})
	}
}

func (u UserServer) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (u UserServer) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
