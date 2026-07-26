package validate

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Failure is this project's canonical validation error representation.
// It isolates validator/v10's validation error terminology
// from the rest of the project
type Failure struct {
	Field string
	Rule  string
	Param string
}

// returns validation failures if err is validation error.
//   - if ok returns validation failures
//   - if not ok return nil
func Failures(err error) (failures []Failure, ok bool) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		failures = make([]Failure, 0, len(ve))
		for _, fe := range ve {
			failures = append(failures, Failure{
				Field: fe.Field(),
				Rule:  fe.Tag(),
				Param: fe.Param(),
			})
		}
		return failures, true
	}

	return nil, false
}
