package testClient

import (
	"context"
	"database/sql"
	"io"
	"net/http"
)

type TestClient struct {
	db *sql.DB
}

var _ ClientInterface = (*TestClient)(nil)

func NewTestClient(db *sql.DB) *TestClient {
	return &TestClient{db: db}
}

func (t TestClient) PostApiAuthWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestClient) PostApiAuth(ctx context.Context, body PostApiAuthJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestClient) GetApiBuyItem(ctx context.Context, item string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestClient) GetApiInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestClient) PostApiSendCoinWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}

func (t TestClient) PostApiSendCoin(ctx context.Context, body PostApiSendCoinJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	//TODO implement me
	panic("implement me")
}
