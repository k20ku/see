package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    ErrCode      `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type ErrCode string

type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code"`
	Param string `json:"param,omitempty"`
}

var (
	ErrValidation  ErrCode = "validation_error"
	ErrInternal    ErrCode = "internal_server_error"
	ErrInvalidJson ErrCode = "invalid_json"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("respond json marshals body: %w", err)
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		return fmt.Errorf("write response error: %v", err)
	}
	return nil
}

func RespondError(ctx context.Context, w http.ResponseWriter, status int, er ErrResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bbs, err := json.Marshal(er)
	if err != nil {
		log.Printf("JSON marshalling response error: %v", err)
		return
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bbs); err != nil {
		log.Printf("writing error response body: %v", err)
		return
	}
}

func RespondErrValidation(ctx context.Context, w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		details := make([]FieldError, len(ve))
		for _, fe := range ve {
			details = append(details, FieldError{
				Field: fe.Field(),
				Code:  fe.Tag(),
			})
		}
		er := ErrResponse{Error: Error{
			Code:    ErrValidation,
			Message: "validation failed",
			Details: details,
		}}

		RespondError(ctx, w, http.StatusBadRequest, er)

	} else {
		log.Printf("err validation unknown error: %v", ve)
		return
	}

}

func RespondErrInternal(ctx context.Context, w http.ResponseWriter) {
	er := ErrResponse{Error: Error{
		Code:    ErrInternal,
		Message: "Internal Server Error",
	}}
	RespondError(ctx, w, http.StatusInternalServerError, er)
}

func RespondErrInvalidJson(ctx context.Context, w http.ResponseWriter) {
	er := ErrResponse{Error: Error{
		Code:    ErrInvalidJson,
		Message: "Invalid JSON",
	}}
	RespondError(ctx, w, http.StatusBadRequest, er)
}
