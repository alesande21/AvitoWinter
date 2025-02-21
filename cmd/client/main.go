package main

import (
	testClient "AvitoWinter/test/e2e/client"
	"context"
	log2 "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := testClient.NewClientWithResponses("http://localhost:8080/")
	if err != nil {
		log2.Errorf("testClient:%v", err)
		return
	}
	var authRequest testClient.PostApiAuthJSONRequestBody
	authRequest.Username = "alesande123"
	authRequest.Password = "alesande12@"

	authResponse, err := c.PostApiAuthWithResponse(ctx, authRequest)
	if err != nil {
		log2.Errorf("PostApiAuth:%v", err)
		return
	}

	if authResponse == nil {
		log2.Printf("authResponse == nil")
		return
	}

	token := *authResponse.JSON200
	if token.Token != nil {
		if len(*token.Token) > 0 {
			log2.Printf("Token: %v", *token.Token)
		} else {
			log2.Printf("Error get token: %v", *token.Token)
		}

	} else {
		log2.Printf("Error get token: nil")
	}
}
