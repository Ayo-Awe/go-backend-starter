package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ayo-awe/go-backend-starter/internal/apperrors"
	"github.com/ayo-awe/go-backend-starter/internal/oapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

func (s Server) ValidationMiddleware() (func(next http.Handler) http.Handler, error) {
	// disables printing details about schema validation errors
	openapi3.SchemaErrorDetailsDisabled = true

	spec, err := oapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	spec.Servers = nil

	// setup custom schema stuff here
	options := &openapi3filter.Options{
		MultiError: false,
	}

	// avoid returning detailed schema errors
	options.WithCustomSchemaErrorFunc(func(err *openapi3.SchemaError) string { return err.Reason })

	router, err := gorillamux.NewRouter(spec)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// validate request
			if err := s.validateRequest(r, router, options); err != nil {
				s.handleAppError(w, r, err)
				return
			}

			// serve
			next.ServeHTTP(w, r)
		})
	}, nil
}

// validateRequest is called from the middleware above and actually does the work
// of validating a request.
func (s Server) validateRequest(r *http.Request, router routers.Router, options *openapi3filter.Options) error {

	// Find route
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		return apperrors.New(err.Error(),
			oapi.ErrorCodeUnknownEndpoint,
			http.StatusNotFound,
		)
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
	}

	if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
		// We're not supporting multi-errors, so we just handle the first error
		me := openapi3.MultiError{}
		if errors.As(err, &me) {
			err = me[0]
		}

		switch e := err.(type) {
		case *openapi3filter.RequestError:
			return apperrors.New(e.Error(),
				oapi.ErrorCodeInvalidRequestParameters,
				http.StatusBadRequest,
			)

		case *openapi3filter.SecurityRequirementsError:
			var appErr apperrors.AppError
			if errors.As(err, &appErr) {
				return appErr
			}

			// this happens when we don't pass in an authentication function!!!
			if errors.Is(err, openapi3filter.ErrAuthenticationServiceMissing) {
				s.logger.Error(
					"openapi3filter authentication function is missing",
					slog.String("error", err.Error()),
				)

				return apperrors.New("an unexpected error occurred",
					oapi.ErrorCodeUnexpectedError,
					http.StatusForbidden,
				)
			}

			s.logger.Error(
				"unknown secuirty requirements error",
				slog.String("error", e.Error()),
			)
			return apperrors.New(e.Error(), oapi.ErrorCodeSecurityRequirementsNotSatisfied, http.StatusUnauthorized)
		default:
			// This should never happen today, but if our upstream code changes,
			// we don't want to crash the server, so handle the unexpected error.
			return fmt.Errorf("error validating route: %w", err)
		}
	}

	return nil
}
