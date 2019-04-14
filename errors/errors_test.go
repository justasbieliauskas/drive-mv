package errors_test

import (
	go_errors "errors"
	"fmt"
	"strings"
	"testing"

	"github.com/justasbieliauskas/drivemv/errors"
)

func TestNesting(t *testing.T) {
	messages := []string{"foo", "bar", "baz"}
	err := errors.Nest(
		messages[0],
		errors.Nest(
			messages[1],
			go_errors.New(messages[2]),
		),
	)
	result := fmt.Sprint(err)
	expected := strings.Join(messages, "\n")
	if result != expected {
		t.Errorf(
			"Nested error message incorrect!\nExpected: %s\nActual: %s\n",
			expected,
			result,
		)
	}
}
