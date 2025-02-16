// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AuthRequest defines model for AuthRequest.
type AuthRequest struct {
	// Password Пароль для аутентификации.
	Password string `json:"password"`

	// Username Имя пользователя для аутентификации.
	Username string `json:"username"`
}

// AuthResponse defines model for AuthResponse.
type AuthResponse struct {
	// Token JWT-токен для доступа к защищенным ресурсам.
	Token *string `json:"token,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Errors Сообщение об ошибке, описывающее проблему.
	Errors *string `json:"errors,omitempty"`
}

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	CoinHistory *struct {
		Received *[]struct {
			// Amount Количество полученных монет.
			Amount *int `json:"amount,omitempty"`

			// FromUser Имя пользователя, который отправил монеты.
			FromUser *string `json:"fromUser,omitempty"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			// Amount Количество отправленных монет.
			Amount *int `json:"amount,omitempty"`

			// ToUser Имя пользователя, которому отправлены монеты.
			ToUser *string `json:"toUser,omitempty"`
		} `json:"sent,omitempty"`
	} `json:"coinHistory,omitempty"`

	// Coins Количество доступных монет.
	Coins     *int `json:"coins,omitempty"`
	Inventory *[]struct {
		// Quantity Количество предметов.
		Quantity *int `json:"quantity,omitempty"`

		// Type Тип предмета.
		Type *string `json:"type,omitempty"`
	} `json:"inventory,omitempty"`
}

// SendCoinRequest defines model for SendCoinRequest.
type SendCoinRequest struct {
	// Amount Количество монет, которые необходимо отправить.
	Amount int `json:"amount"`

	// ToUser Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser"`
}

// PostApiAuthJSONRequestBody defines body for PostApiAuth for application/json ContentType.
type PostApiAuthJSONRequestBody = AuthRequest

// PostApiSendCoinJSONRequestBody defines body for PostApiSendCoin for application/json ContentType.
type PostApiSendCoinJSONRequestBody = SendCoinRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	// (POST /api/auth)
	PostApiAuth(w http.ResponseWriter, r *http.Request)
	// Купить предмет за монеты.
	// (GET /api/buy/{item})
	GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string)
	// Получить информацию о монетах, инвентаре и истории транзакций.
	// (GET /api/info)
	GetApiInfo(w http.ResponseWriter, r *http.Request)
	// Отправить монеты другому пользователю.
	// (POST /api/sendCoin)
	PostApiSendCoin(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostApiAuth operation middleware
func (siw *ServerInterfaceWrapper) PostApiAuth(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiAuth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiBuyItem operation middleware
func (siw *ServerInterfaceWrapper) GetApiBuyItem(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "item" -------------
	var item string

	err = runtime.BindStyledParameterWithOptions("simple", "item", mux.Vars(r)["item"], &item, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "item", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiBuyItem(w, r, item)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiInfo operation middleware
func (siw *ServerInterfaceWrapper) GetApiInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiInfo(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostApiSendCoin operation middleware
func (siw *ServerInterfaceWrapper) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiSendCoin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{})
}

type GorillaServerOptions struct {
	BaseURL          string
	BaseRouter       *mux.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r *mux.Router) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r *mux.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, GorillaServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
//func HandlerWithOptions(si ServerInterface, options GorillaServerOptions) http.Handler {
//	r := options.BaseRouter
//
//	if r == nil {
//		r = mux.NewRouter()
//	}
//	if options.ErrorHandlerFunc == nil {
//		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
//			http.Error(w, err.Error(), http.StatusBadRequest)
//		}
//	}
//	wrapper := ServerInterfaceWrapper{
//		Handler:            si,
//		HandlerMiddlewares: options.Middlewares,
//		ErrorHandlerFunc:   options.ErrorHandlerFunc,
//	}
//
//	r.HandleFunc(options.BaseURL+"/api/auth", wrapper.PostApiAuth).Methods("POST")
//
//	r.HandleFunc(options.BaseURL+"/api/buy/{item}", wrapper.GetApiBuyItem).Methods("GET")
//
//	r.HandleFunc(options.BaseURL+"/api/info", wrapper.GetApiInfo).Methods("GET")
//
//	r.HandleFunc(options.BaseURL+"/api/sendCoin", wrapper.PostApiSendCoin).Methods("POST")
//
//	return r
//}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xXy27bRhT9FWLaJSspfQCBdnbRh7MKmhRZGF7Q0thianHomaELwRCgR9wmsGEXXRUB",
	"WqPtD9CqVNOyRP3CnT8q7iVlvW0ZsbNovBOHI94zZ86dc2afFUTZFx73tGL5faYKJV526OdKoEvf8d2A",
	"K42PvhQ+l9rl9NJ3lPpRyCL+LnJVkK6vXeGxPINTCE0NYrg0Rxa04dKcWBCapmlAB/qmAZF5BRF0ITQ/",
	"QQRRhtlsS8iyo1l+9Fmb6YrPWZ4pLV1vm1VtFiguPafM55T8DXpYZZBUhXOIoQUhVaTyy6GYqli1meS7",
	"gSt5keXXR+XtEcqNqz+JzZe8oBFmQpvyhaf4LG9a/MC92RU8efH8E9OAGLoI7wpwG2JTNw3ThAGEFnQt",
	"OIfQvIHIvMF50DeH0LNMDTqmbpqmZuoQQm/+WmaAfiWlkIuRcnyt5pD9J8QQw1kKIYKOhY8WxOY1RHCG",
	"S7BxaACRqZtD2oljmt2xYEDSOINL6EDPNJeEuuZticVIC8L1vnWVFrIy+1LyAnf3OAnV1bysZqc4ZRF4",
	"es5K36KeIDI/E78NaEE8FFkTB5MdMAcW9CCGPnRMY2xBrqf5NpeIf0uK8veKy1tL17agCzEqw9TMIVxY",
	"9IAkhtCCCC7HSpvDJdlMBxwpnQo+K56s/k7oGcd3eQuKtHhngiBGSc2BYA7fnaZ5M1B4allixnt5OUpc",
	"b497Q1Uv2JzdwPG0qytLqxcPC2hDD8silQt2g0ZmPvkXRDCY/kh4Z3w+417xS+F6C13ndlK8oneqizoW",
	"DtMZdgAxtCHCqVOtZRrm6N6V2jdN+Bf6c4svIdlxe0pR2UOOZq2JOr0QSFdXnqHLJ5SuckdyiaaFT5v0",
	"9PXQjp+8eM7sJBPgl5K3IyglrX1WrZJUtwTZm6t38M3K0zVrZc/VwlIl4TOb7XGpEpoeZXKZHPIofO45",
	"vsvy7DMaQl/VJQKVdXw366SYfJFIAYXgINdrRZZnT4XSK75LwBMmuNKrolhJHMHT6aHm+P6OW6D/ZV8q",
	"4Y1CDv76WPItlmcfZUcpKJtGoOx4/qlO0q1lwGkgsSTC/Gkud8elU7+j2lNS+9vUYQAd8xr6EF4bbcxJ",
	"Bsn+/A7RTQaHefB+hw60oGNqlFAuKLYk1m/qKZxH7xlOCK2096JhZ0KfsHzxXqn5FZveNOgERW88wd0b",
	"RafQMnUiLqEvzEz0Lcuvb9hMBeWyg67A4JfF225BNB1WMKtNJk0IMxacIikW6Qnrxrhh10TlRSfcEUKP",
	"4RzaENKpVSdhJrxDj+alJzR0MXHjyqjTN4NKdh8Nror8bvM57f4Nx25fDSprmpfpqJBOmWuOCXV9n7ke",
	"XR7oLEjuCGSYbLpp7bF9nD5PN+Y39OLOGwWyVmLiD332/+mzSWdc36hONt5bynGpUU/kIdqEKfO+UvrQ",
	"J6/RON512D16y8Rd6gZveVD4h6vw0yvrSFUeQd+8onX3Uo85tiZiNoTmwKZ5WJOsI0R05EQR5fKEtcgi",
	"2CH0aTu75CsXY22i0ovAjfFveGO4pwg4fSFZPgY+9NRDT8321B/XXvIsaJuaacI/w+vh/Jx1nEmg31SW",
	"y71hPArkTnpdy2ezO6Lg7JSE0vnHucc5Vt2o/hcAAP//D2BZ7QUWAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
