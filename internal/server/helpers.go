package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ayo-awe/go-backend-starter/internal/apperrors"
	"github.com/ayo-awe/go-backend-starter/internal/oapi"
	validator "github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func (s Server) sendJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(bytes)
	if err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		s.logError(r, fmt.Errorf("failed to write json response: %w", err))
	}
}

func (s Server) logError(r *http.Request, err error) {
	s.logger.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
}

func (s Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)
	s.errorResponse(w, r, oapi.ErrorCodeUnexpectedError, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

func (s Server) errorResponse(w http.ResponseWriter, r *http.Request, errorCode oapi.ErrorCode, statusCode int, message string) {
	err := oapi.Error{Message: message, ErrorCode: errorCode}
	s.sendJSONResponse(w, r, statusCode, err)
}

func (s Server) handleAppError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr apperrors.AppError
	if errors.As(err, &appErr) {
		payload := oapi.Error{
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		}
		s.sendJSONResponse(w, r, appErr.StatusCode, payload)
		return
	}

	s.serverErrorResponse(w, r, err)
}

// this function uses the validate tags on the request struct to validate the request body
// Use the x-oapi-codegen-extra-tags key in the openapi spec to add validation tags to the struct
// that are not supported by oapi-codegen.
// See more here: https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#openapi-extensions
func (s Server) validateRequestStruct(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var invalidErr *validator.InvalidValidationError

	// this should only happen if the data passed into this function isn't a struct
	if errors.As(err, &invalidErr) {
		s.logger.Error("non-struct type passed into request validation function",
			"error", err.Error(),
		)

		return apperrors.New("unexpected error",
			oapi.ErrorCodeUnexpectedError,
			http.StatusInternalServerError,
		)
	}

	validationErrors := err.(validator.ValidationErrors)

	// return only the first validation error, we can return field errors in the future
	return apperrors.New(validationErrors[0].Error(),
		oapi.ErrorCodeInvalidRequestParameters,
		http.StatusBadRequest,
	)
}
