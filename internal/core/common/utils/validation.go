package utils

import (
	errorcode "Eccomerce-website/internal/core/entity/error_code"
	"Eccomerce-website/internal/core/model/response"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

var (
	PhoneRe = regexp.MustCompile(`^(\+251|251|0)?[79]\d{8}$`)
)

func PasswordValidation(password string) *response.Response {
	isValid := false
	isUpper := false
	isLower := false
	isDigit := false
	isSpecialChar := false

	specialChar := "!@#$%^&*+_-?></|"

	for _, char := range password {
		if unicode.IsUpper(char) {
			isUpper = true
		}
		if unicode.IsLower(char) {
			isLower = true
		}
		if unicode.IsDigit(char) {
			isDigit = true
		}
		if strings.ContainsRune(specialChar, char) {
			isSpecialChar = true
		}
		if isUpper && isLower && isDigit && isSpecialChar {
			isValid = true
			break
		}
	}

	if !isValid {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "the password must contain at least one uppercase letter, lowercase letter, digit and special character",
		}
		return &errorResponse
	}

	if len(password) < 8 {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "password must be at least 8 characters long",
		}
		return &errorResponse
	}
	return nil
}

func PhoneValidation(phoneNumber string) (*response.Response, string) {
	if isPhoneValid := PhoneRe.MatchString(phoneNumber); !isPhoneValid {
		errorResponse := response.Response{
			Status:       http.StatusBadRequest,
			ErrorType:    errorcode.ValidationError,
			ErrorMessage: "invalid phone number format.Please enter a valid phone number",
		}
		return &errorResponse, ""
	}

	if phoneNumber[:1] == "0" {
		phoneNumber = phoneNumber[1:]
	}
	if phoneNumber[:3] == "251" {
		phoneNumber = phoneNumber[3:]
	}
	if phoneNumber[:4] != "+251" {
		phoneNumber = "+251" + phoneNumber
	}

	return nil, phoneNumber
}
