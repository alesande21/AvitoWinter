package main

import (
	testClient "AvitoWinter/test/e2e/client"
	"context"
	log2 "github.com/sirupsen/logrus"
)

//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	c := testClient.NewTestClient()
//
//	body := strings.NewReader(`{"username": "alesande123", "password": "alesande12@"}`)
//
//	authResponse, err := c.PostApiAuthWithBody(ctx, "application/json", body)
//	if err != nil {
//		log2.Errorf("PostApiAuthWithBody:%v", err)
//		return
//	}
//
//	if authResponse == nil {
//		log2.Printf("authResponse == nil")
//		return
//	}
//
//	//authResponse.JSON200.Token
//	log2.Printf("Status: %d", authResponse.StatusCode())
//	log2.Printf("Body: %s", string(authResponse.Body))
//
//	r := bytes.NewReader(authResponse.Body)
//	err = json.NewDecoder(r).Decode(&authResponse.JSON200)
//	if err != nil {
//		log2.Printf("Decode error")
//		return
//	}
//
//	token := authResponse.JSON200
//
//	if authResponse == nil {
//		log2.Printf("authResponse == nil")
//		return
//	}
//
//	if token == nil {
//		log2.Printf("token == nil")
//		return
//	}
//
//	if token.Token == nil {
//		log2.Printf("token.Token == nil")
//		return
//	}
//
//	if len(*token.Token) > 0 {
//		log2.Printf("Token: %v", *token.Token)
//	} else {
//		log2.Printf("Error get token: %v", *token.Token)
//	}
//
//}

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
	//require.NoError(t,)
	if err != nil {
		log2.Errorf("PostApiAuth:%v", err)
		return
	}

	if authResponse == nil {
		log2.Printf("authResponse == nil")
		return
	}

	//authResponse.JSON200.Token
	log2.Printf("Status: %d", authResponse.StatusCode())
	log2.Printf("Body: %s", string(authResponse.Body))

	//r := bytes.NewReader(authResponse.Body)
	//err = json.NewDecoder(r).Decode(&authResponse.JSON200)
	//if err != nil {
	//	log2.Printf("Decode error")
	//	return
	//}

	token := authResponse.JSON200

	if authResponse == nil {
		log2.Printf("authResponse == nil")
		return
	}

	if token == nil {
		log2.Printf("token == nil")
		return
	}

	if token.Token == nil {
		log2.Printf("token.Token == nil")
		return
	}

	if len(*token.Token) > 0 {
		log2.Printf("Token: %v", *token.Token)
	} else {
		log2.Printf("Error get token: %v", *token.Token)
	}

}
