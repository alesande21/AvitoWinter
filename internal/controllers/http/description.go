package http

import (
	entity2 "AvitoWinter/internal/entity"
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
	var errorDescription string
	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		log2.Errorf("PostApiAuth-> json.NewDecoder: неверный формат для регистрационных данных пользователя: логин и пароль: %s", err.Error())
		errorDescription = "Неверный формат для регистрационных данных пользователя: логин и пароль."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	newUserCredentials, err := entity2.NewUserCredentials(authRequest.Username, authRequest.Password)
	if err != nil {
		errorDescription = "Неверный формат для регистрационных данных пользователя: логин и пароль."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	tokenString, err := u.service.AuthenticationUser(r.Context(), *newUserCredentials)
	if err != nil {
		errorDescription = "Аутификация не пройдена."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	var authResponse AuthResponse
	authResponse.Token = &tokenString

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(authResponse); err != nil {
		errorDescription = "Ошибка кодирования ответа."
		log2.Errorf("CreateUser-> json.NewEncoder: ошибка при кодирования овета: %s", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, ErrorResponse{Errors: &errorDescription})
	}
}

func (u UserServer) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	var errorDescription string

	userUUID, ok := r.Context().Value("user_value").(string)
	if !ok || userUUID == "" {
		errorDescription = "User not authenticated"
		sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
		return
	}

	purchaseInfo, err := entity2.NewPurchaseInfo(userUUID, item, 1)
	if err != nil {
		errorDescription = "Не задан предмет покупки."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	err = u.service.PurchaseItem(r.Context(), *purchaseInfo)
	if err != nil {
		errorDescription = "Какая то там ошибка"
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var sb strings.Builder
	sb.WriteString("Покупка + ")
	if err := json.NewEncoder(w).Encode(sb.String()); err != nil {
		log2.Errorf("CreateUser-> json.NewEncoder: ошибка при кодирования овета: %s", err.Error())
		errorDescription = "Ошибка кодирования ответа."
		sendErrorResponse(w, http.StatusInternalServerError, ErrorResponse{Errors: &errorDescription})
	}
}

func (u UserServer) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	var errorDescription string

	userUUID, ok := r.Context().Value("user_value").(string)
	if !ok || userUUID == "" {
		errorDescription = "User not authenticated"
		sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
		return
	}
}

func (u UserServer) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	var errorDescription string

	userUUID, ok := r.Context().Value("user_value").(string)
	if !ok || userUUID == "" {
		errorDescription = "User not authenticated"
		sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
		return
	}

	var sendRequest SendCoinRequest

	if err := json.NewDecoder(r.Body).Decode(&sendRequest); err != nil {
		log2.Errorf("PostApiSendCoin-> json.NewDecoder: неверный формат для регистрационных данных пользователя: логин и пароль: %s", err.Error())
		errorDescription = "Неверный формат для трансфера коинов."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	transferInfo, err := entity2.NewTransferInfo(userUUID, sendRequest.ToUser, sendRequest.Amount)
	if err != nil {
		log2.Errorf("PostApiSendCoin-> json.NewDecoder: неверный формат для регистрационных данных пользователя: логин и пароль: %s", err.Error())
		errorDescription = "Неверный формат для трансфера коинов."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	err = u.service.TransferCoin(r.Context(), *transferInfo)
	if err != nil {
		log2.Errorf("PostApiSendCoin-> json.NewDecoder: неверный формат для регистрационных данных пользователя: логин и пароль: %s", err.Error())
		errorDescription = "Неверный формат для трансфера коинов."
		sendErrorResponse(w, http.StatusBadRequest, ErrorResponse{Errors: &errorDescription})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	var sb strings.Builder
	sb.WriteString("Трансфер реализован.")
	if err := json.NewEncoder(w).Encode(sb.String()); err != nil {
		log2.Errorf("CreateUser-> json.NewEncoder: ошибка при кодирования овета: %s", err.Error())
		errorDescription = "Ошибка кодирования ответа."
		sendErrorResponse(w, http.StatusInternalServerError, ErrorResponse{Errors: &errorDescription})
	}
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
