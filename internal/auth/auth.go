package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("MySecretKey")

type JWTClaim struct {
	uuid string
	jwt.StandardClaims
}

func GenerateJWT(uuid string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		uuid: uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("-> token.SignedString%v", err)
	}

	return tokenString, nil
}

func ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		return fmt.Errorf("-> jwt.ParseWithClaims%v", err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return fmt.Errorf("-> token.Claims: не возможно распарсить claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return fmt.Errorf(": token истек")
	}

	return nil
}
