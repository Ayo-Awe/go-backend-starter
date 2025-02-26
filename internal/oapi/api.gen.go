// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package oapi

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

	externalRef0 "github.com/ayo-awe/go-backend-starter/internal/rbac"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

const (
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeInvalidAuthToken                 ErrorCode = "INVALID_AUTH_TOKEN"
	ErrorCodeInvalidCredentials               ErrorCode = "INVALID_CREDENTIALS"
	ErrorCodeInvalidRequestParameters         ErrorCode = "INVALID_REQUEST_PARAMETERS"
	ErrorCodeRateLimitExceeded                ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeSecurityRequirementsNotSatisfied ErrorCode = "SECURITY_REQUIREMENTS_NOT_SATISFIED"
	ErrorCodeUnexpectedError                  ErrorCode = "UNEXPECTED_ERROR"
	ErrorCodeUnknownEndpoint                  ErrorCode = "UNKNOWN_ENDPOINT"
	ErrorCodeUserAlreadyExists                ErrorCode = "USER_ALREADY_EXISTS"
)

// Error defines model for Error.
type Error struct {
	// ErrorCode List of all custom API error codes
	ErrorCode ErrorCode `json:"error_code"`
	Message   string    `json:"message"`
}

// ErrorCode List of all custom API error codes
type ErrorCode string

// UserProfile defines model for UserProfile.
type UserProfile struct {
	Scopes []externalRef0.Scope `json:"scopes"`
}

// BadRequestError defines model for BadRequestError.
type BadRequestError = Error

// ConflictError defines model for ConflictError.
type ConflictError = Error

// ForbiddenError defines model for ForbiddenError.
type ForbiddenError = Error

// InternalServerError defines model for InternalServerError.
type InternalServerError = Error

// NotFound defines model for NotFound.
type NotFound = Error

// UnauthorizedError defines model for UnauthorizedError.
type UnauthorizedError = Error

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Welcome to Go Starter API
	// (GET /api/v1)
	Welcome(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Welcome to Go Starter API
// (GET /api/v1)
func (_ Unimplemented) Welcome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Welcome operation middleware
func (siw *ServerInterfaceWrapper) Welcome(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, ApiKeyAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Welcome(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
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
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/v1", wrapper.Welcome)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xW0W7jthL9FYL3Pipxtn3Tm9ZmWiGO4kpyN0EQCAw1jrmRSC5JeeMGAvoP/cN+STGU",
	"nNj1ot2++E0hZ86ZmXM48SsVujVagfKOxq/UgjNaOQh/fOR1Dl86cJ5Zqy0eCa08KI+f3JhGCu6lVpPP",
	"Tis8c2INLcev/1tY0Zj+b/KOPxlu3WRA6/s+ojU4YaVBEBrTVG14I2silek87SM61WrVSHEq/hyc7qwA",
	"whsLvN4SeJHOO6zkUttHWdegTlRKIgQ4R1Y7WvLn73+QzoEltQZHlPZkzTdALHzppIWaGCs3soEnCOWm",
	"yoNVvCnAbsCeqOaBjMBwH9FM+0vdqfqEwuFcVoGzj+hS8c6vtZW/QX1iBwsLNSgveeMoRo2JiPtWibHa",
	"gPVyeGxhapXQNXwX+RQD+4i24Bx/Cjl+a4DG1Hkr1VNg3ZmDxvdvgdE+00O0y9KPn0GEN/cOH7/+rbu5",
	"dJ7oFeFNQ0TnvG5JskgHwQkCOoRXXYuEafZrMk9nVbIsf67KmyuW0YguC5ZXyTxnyeyuYrdpURY0egvN",
	"2S9LVpTVIsmTa1ayHC+XGbtdsGnJZhXL85ucRjRPSlbN0+u0rNjtlLEZm4XAq+zmU1axbLa4SbOSRrRg",
	"02WelncBOc3ZNcvKospuyqpIyrS4TEPijn6asxnLyjSZF3uTGecZ0ZczbC3jLQp2v5M66fy61M+gsAIH",
	"NhlWBxs2R7QLGxfpglveggeLV0sFLwaE37kzojn3MJet9OxFANRQh6hnpb8qpmqjpfLYFYjOSr/NB3lb",
	"dEemfcG9dCsZkkbW6Z4NH/qIFkKboOtOpEY6L9WTiy1sJHw9brsfulpYvZINHNvWIWL4kh5a92/eHQro",
	"31i4tXx75NUR9Nid+JLG5gtEHJgTI69gi0KEOtCpa+A14EQVbxHg9ixZpGdX7I6+U4csrOUjcAt2lx9K",
	"xYDHcPyesPbeDG9eqpU+fh0/afLIxTOomjjPrQeLr4NGtJEClAvDG8tZWG2sBM/tlka0s80I7+LJZMw9",
	"F7qdeLCtO+OqPhNa1RKJwnr30qMYSFkcUG3AuqGaD+cX5xcYqw0obiSN6Y/hKKKG+3UY3IQbOdl8wM8n",
	"CCsRpQ0LMa1pTD9BI3SLO+Pgd8EPFxf/aZMeOmZvY8ELb01oZGQiXpO9ni5l07qxs+/bbt+wzNGO/jYX",
	"sgR/dW2Lsvxz3J4RaXx/aMH7h/4Br/GfoQu3hwJzI8/3RKb9Q/9XAAAA///aBkTthwkAAA==",
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

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "internal/rbac/scopes.yml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
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
