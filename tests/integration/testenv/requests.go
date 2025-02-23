package testenv

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type RequestOption func(r *http.Request)

func WithAuthToken(token string) RequestOption {
	return func(r *http.Request) {
		r.Header.Set("Authorization", "Bearer "+token)
	}
}

func DoGet(r http.Handler, url string, options ...RequestOption) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, url, nil)

	for _, option := range options {
		option(req)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func DoPost(r http.Handler, url string, body any, options ...RequestOption) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	req := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	for _, option := range options {
		option(req)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func DoPut(r http.Handler, url string, body any, options ...RequestOption) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	for _, option := range options {
		option(req)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func DoPatch(r http.Handler, url string, body any, options ...RequestOption) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil
	}

	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	for _, option := range options {
		option(req)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func DoDelete(r http.Handler, url string, options ...RequestOption) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, url, nil)

	for _, option := range options {
		option(req)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}
