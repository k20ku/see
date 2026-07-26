package respond

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/k20ku/see/internal/validate"
)

// ---- primitives respond ----

// respond.JSON sends body as a HTTP response in JSON format
// with provided status code.
func JSON(ctx context.Context, w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("respond JSON: marshaling body: %w", err)
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		return fmt.Errorf("respond JSON: write response error: %w", err)
	}
	return nil
}

func Error(
	ctx context.Context, w http.ResponseWriter,
	status int,
	errRsp ErrorResponse,
) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bbs, err := json.Marshal(errRsp)
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

// ---- routine respond ----
func ValidationError(ctx context.Context, w http.ResponseWriter, err error) {
	fs, ok := validate.Failures(err)
	if !ok {
		log.Printf("respond validation error: unexpected error: %v", err)
		InternalServerError(ctx, w)
		return
	}
	details := make([]FieldError, 0, len(fs))
	for _, f := range fs {
		details = append(details, FieldError{Field: f.Field, Code: f.Rule, Param: f.Param})
	}
	Error(ctx, w, http.StatusBadRequest,
		ErrorResponse{Body: ErrorBody{
			Code:    CodeValidationError,
			Message: "Validation Failed",
			Details: details,
		}},
	)
}

func InternalServerError(ctx context.Context, w http.ResponseWriter) {
	Error(ctx, w, http.StatusInternalServerError,
		ErrorResponse{Body: ErrorBody{
			Code:    CodeInternalServerError,
			Message: "Internal Server Error",
		}},
	)
}

func InvalidJSONError(ctx context.Context, w http.ResponseWriter) {
	Error(ctx, w, http.StatusBadRequest,
		ErrorResponse{Body: ErrorBody{
			Code:    CodeInvalidJSONError,
			Message: "Invalid JSON",
		}},
	)
}
