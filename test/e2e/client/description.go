package testClient

import (
	"context"
	"io"
	"net/http"
)

type TestClient struct{}

var _ ClientInterface = (*TestClient)(nil)

func NewTestClient() *TestClient {
	return &TestClient{}
}

func (t TestClient) PostApiAuthWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/api/auth", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
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
