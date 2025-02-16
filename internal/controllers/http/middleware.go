package http

import (
	auth2 "AvitoWinter/internal/auth"
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// исключаем аутификацию
		if r.URL.Path == "/api/auth" {
			next.ServeHTTP(w, r)
			return
		}

		var errorDescription string
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorDescription = "Требуется авторизация."
			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			errorDescription = "Неверный формат токена."
			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		tokenString := tokenParts[1]

		claims, err := auth2.ValidateToken(tokenString)
		if err != nil {
			errorDescription = "Невалидный токен."
			sendErrorResponse(w, http.StatusUnauthorized, ErrorResponse{Errors: &errorDescription})
			return
		}

		ctx := context.WithValue(r.Context(), "user_uuid", claims.UserUUID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//func middleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
//		if len(authHeader) != 2 {
//			fmt.Println("Malformed token")
//			w.WriteHeader(http.StatusUnauthorized)
//			w.Write([]byte("Malformed Token"))
//		} else {
//			jwtToken := authHeader[1]
//			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
//				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//				}
//				return []byte(SECRETKEY), nil
//			})
//			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//				ctx := context.WithValue(r.Context(), "props", claims)
//				// Access context values in handlers like this
//				// props, _ := r.Context().Value("props").(jwt.MapClaims)
//				next.ServeHTTP(w, r.WithContext(ctx))
//			} else {
//				fmt.Println(err)
//				w.WriteHeader(http.StatusUnauthorized)
//				w.Write([]byte("Unauthorized"))
//			}
//		}
//	})
//}
