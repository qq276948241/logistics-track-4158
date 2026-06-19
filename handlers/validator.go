package handlers

import (
	"regexp"
	"unicode/utf8"
)

var (
	trackingNumberRegex = regexp.MustCompile(`^[A-Za-z0-9\-]+$`)
	carrierRegex        = regexp.MustCompile(`^[\p{Han}A-Za-z0-9 ]+$`)
)

const (
	trackingNumberMinLength = 8
	trackingNumberMaxLength = 30
	carrierMinLength        = 2
	carrierMaxLength        = 50
)

func ValidateTrackingNumber(number string) (bool, string) {
	if number == "" {
		return false, "Tracking number is required"
	}
	runeLen := utf8.RuneCountInString(number)
	if runeLen < trackingNumberMinLength || runeLen > trackingNumberMaxLength {
		return false, "Tracking number length must be between 8 and 30 characters"
	}
	if !trackingNumberRegex.MatchString(number) {
		return false, "Tracking number can only contain letters, digits and hyphens"
	}
	hasDigit := false
	for _, r := range number {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false, "Tracking number must contain at least one digit"
	}
	return true, ""
}

func ValidateCarrier(carrier string) (bool, string) {
	if carrier == "" {
		return false, "Carrier is required"
	}
	runeLen := utf8.RuneCountInString(carrier)
	if runeLen < carrierMinLength || runeLen > carrierMaxLength {
		return false, "Carrier length must be between 2 and 50 characters"
	}
	if !carrierRegex.MatchString(carrier) {
		return false, "Carrier can only contain Chinese characters, letters, digits and spaces"
	}
	return true, ""
}
