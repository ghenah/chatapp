package httpserver

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationErrors struct {
	HasErrors   bool
	HasCritical bool
	List        []error
}

func (ve ValidationErrors) Error() string {
	errString := []string{}
	for _, e := range ve.List {
		errString = append(errString, e.Error())
	}

	return strings.Join(errString, ",")
}

var (
	// Password validator errors
	ErrorValidationPasswordTooLong      = errors.New("password too long")
	ErrorValidationLowercaseMissing     = errors.New("lowercase letter missing")
	ErrorValidationUppercaseMissing     = errors.New("uppercase letter missing")
	ErrorValidationNumberigMissing      = errors.New("numeric character missing")
	ErrorValidationSpecialChMissing     = errors.New("special character missing")
	ErrorValidationWhitespaceNotAllowed = errors.New("whitespace not allowed")
	ErrorValidationWrongLength          = fmt.Errorf("length must be between %d and %d", userPasswordMinLength, userPasswordMaxLength)
)
