package main

import (
	"AvitoWinter/internal/app"
	testClient "AvitoWinter/test/e2e/client"
	"context"
	"database/sql"
	log2 "github.com/sirupsen/logrus"
	_ "github.com/stretchr/testify"
	_ "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"testing"
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

var (
	db        *sql.DB
	serverURL = "http://localhost:8080/"
)

func TestMain(m *testing.M) {
	setupTestDB()
	runService()
	m.Run()
}

func setupTestDB() (*sql.DB, error) {
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:15", postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"), postgres.WithPassword("password"))
	if err != nil {
		panic(err)
	}

	conn, _ := container.ConnectionString(ctx)
	db, err = sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func truncateTables() {
	//db.Exec("TRUNCATE products RESTART IDENTITY CASCADE")
}

func runService() {
	go func() {
		err := app.RunApp()
		if err != nil {
			log2.Errorf("app.RunApp%v", err)
			return
		}
	}()

}

func GetToken(t *testing.T) {

	//test push
	// test push 2
}
