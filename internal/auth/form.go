package auth

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type authForm struct {
	email    string
	password string
	Errors   map[string]string
}

func newAuthForm(email string, password string) *authForm {
	return &authForm{
		email:    strings.ToLower(strings.TrimSpace(email)),
		password: password,
		Errors:   make(map[string]string),
	}
}

func (f *authForm) addError(field, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}
	f.Errors[field] = message
}

func (f *authForm) validate() bool {
	if f.email == "" || !emailRegex.MatchString(f.email) {
		f.addError("email", "Please enter a valid email address.")
	}

	if f.password == "" {
		f.addError("password", "password is required")
	}

	return len(f.Errors) == 0
}
