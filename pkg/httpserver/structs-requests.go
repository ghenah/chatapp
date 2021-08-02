package httpserver

import (
	"regexp"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RequestUserRegister struct {
	Username string   `json:"username" validate:"required" example:"johndoe" format:"string"`
	Email    string   `json:"email" validate:"required" example:"johndoe@example.com" format:"string"`
	Password Password `json:"password" validate:"required" format:"string"`
}

func (req RequestUserRegister) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z]{1}[A-Za-z0-9]{1,15}$"))),
		validation.Field(&req.Email, validation.Required, is.EmailFormat),
		validation.Field(&req.Password, validation.Required),
	)
}

type RequestUserAuthenticate struct {
	Username string   `json:"username" validate:"required" example:"johndoe" format:"string"`
	Password Password `json:"password" validate:"required" format:"string"`
}

func (req RequestUserAuthenticate) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z]{1}[A-Za-z0-9]{1,15}$"))),
		validation.Field(&req.Password, validation.Required),
	)
}

type RequestAddUserToList struct {
	UserID   uint `json:"userId" validate:"required" example:"3" format:"uint"`
	FriendID uint `json:"friendId" validate:"required" example:"5" format:"uint"`
}

func (req RequestAddUserToList) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required),
		validation.Field(&req.FriendID, validation.Required),
	)
}

type RequestUserUpdatePassword struct {
	UserID      uint     `json:"userId" validate:"required" example:"3" format:"uint"`
	Username    string   `json:"username" validate:"required" example:"johndoe" format:"string"`
	OldPassword Password `json:"oldPassword" validate:"required" format:"string"`
	NewPassword Password `json:"newPassword" validate:"required" format:"string"`
}

func (req RequestUserUpdatePassword) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required),
		validation.Field(&req.Username, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z]{1}[A-Za-z0-9]{1,15}$"))),
		validation.Field(&req.OldPassword, validation.Required),
		validation.Field(&req.NewPassword, validation.Required),
	)
}

type RequestUserUpdateUsername struct {
	UserID      uint     `json:"userId" validate:"required" example:"3" format:"uint"`
	Username    string   `json:"username" validate:"required" example:"johndoe" format:"string"`
	Password    Password `json:"password" validate:"required" format:"string"`
	NewUsername string   `json:"newUsername" validate:"required" example:"johndoe" format:"string"`
}

func (req RequestUserUpdateUsername) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.UserID, validation.Required),
		validation.Field(&req.Username, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z]{1}[A-Za-z0-9]{1,15}$"))),
		validation.Field(&req.Password, validation.Required),
		validation.Field(&req.NewUsername, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z]{1}[A-Za-z0-9]{1,15}$"))),
	)
}

const userPasswordMinLength = 8
const userPasswordMaxLength = 24
const userPasswordCutoff = 40

type Password string

func (p Password) Validate() error {
	err := ValidationErrors{}

	// If the password is way to long is probably isn't by mistake
	if len(p) > userPasswordCutoff {
		err.HasCritical = true
		err.HasErrors = true
		err.List = append(err.List, ErrorValidationPasswordTooLong)
		return err
	}

	// A list of standard errors that should be eliminated during
	// the validation
	errorsFound := map[error]struct{}{
		ErrorValidationLowercaseMissing: {},
		ErrorValidationUppercaseMissing: {},
		ErrorValidationNumberigMissing:  {},
		ErrorValidationSpecialChMissing: {},
		ErrorValidationWrongLength:      {},
	}

	// Validate password; eliminate the standard errors any add
	// new ones if conditions are met
	if userPasswordMinLength < len(p) && len(p) < userPasswordMaxLength {
		delete(errorsFound, ErrorValidationWrongLength)
	}
	for _, ch := range p {
		switch {
		case unicode.IsLower(ch):
			delete(errorsFound, ErrorValidationLowercaseMissing)
		case unicode.IsUpper(ch):
			delete(errorsFound, ErrorValidationUppercaseMissing)
		case unicode.IsNumber(ch):
			delete(errorsFound, ErrorValidationNumberigMissing)
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			delete(errorsFound, ErrorValidationSpecialChMissing)
		case ch == ' ' || ch == '\t' || ch == '\n':
			errorsFound[ErrorValidationWhitespaceNotAllowed] = struct{}{}
		}
	}

	for e := range errorsFound {
		err.HasErrors = true
		err.List = append(err.List, e)
	}

	if !err.HasErrors {
		return nil
	}

	return err
}
