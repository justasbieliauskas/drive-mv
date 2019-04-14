package errors

import "fmt"

// PrevAwareError is an error that stores a reference to the previous error.
type prevAwareError struct {
	message string
	prev    error
}

// Nest creates an error with a reference to the previous one.
func Nest(message string, prev error) error {
	return &prevAwareError{message, prev}
}

// Error iterates all previous errors and converts each to a string,
// then joins them by a new line.
func (err *prevAwareError) Error() string {
	return err.message + "\n" + fmt.Sprint(err.prev)
}
