package http

import (
	"net/http"
	"strings"
)

func (u UserServer) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errorDescription string
		// Извлечение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorDescription = "Требуется авторизация."
			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		// Проверка формата заголовка (Bearer <token>)
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errorDescription = "Неверный формат токена."

			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		tokenString := tokenParts[1]

		// Проверка валидности токена
		err := u.service.ValidateToken(r.Context(), tokenString)
		if err != nil {
			errorDescription = "Невалидный токен."
			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		// Если токен валиден, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	}
}
