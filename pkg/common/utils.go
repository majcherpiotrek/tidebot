package common

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func IsMobileUserAgent(userAgent string) bool {
	mobileAgents := []string{
		"Android", "iPhone", "iPad", "iPod", "Mobile",
		"BlackBerry", "Windows Phone", "Opera Mini",
	}

	userAgent = strings.ToLower(userAgent)
	for _, agent := range mobileAgents {
		if strings.Contains(userAgent, strings.ToLower(agent)) {
			return true
		}
	}
	return false
}

func FormatPhoneNumber(phoneNumber string) string {
	num, err := phonenumbers.Parse(phoneNumber, "")

	if err != nil {
		return phoneNumber
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL)
}
