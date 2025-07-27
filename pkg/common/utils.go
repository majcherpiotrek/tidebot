package common

import (
	"fmt"
	"strings"
	"time"

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

var DATE_FORMATS = []string{
	"2006-01-02",      // 2023-12-25
	"01/02/2006",      // 12/25/2023
	"1/2/2006",        // 12/5/2023 (no leading zeros)
	"02/01/2006",      // 25/12/2023 (DD/MM/YYYY)
	"2/1/2006",        // 5/12/2023 (D/M/YYYY)
	"Jan 2, 2006",     // Dec 25, 2023
	"January 2, 2006", // December 25, 2023
	"2 Jan 2006",      // 25 Dec 2023
	"2 January 2006",  // 25 December 2023
	"2006/01/02",      // 2023/12/25
	"02-Jan-2006",     // 25-Dec-2023
	"2006.01.02",      // 2023.12.25
	"02.01.2006",      // 2023.12.25
	"2.1.2006",        // 2023.12.25
}

func Today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func Tomorrow() time.Time {
	return Today().Add(24 * time.Hour)
}

func ParseDate(dateStr string) (time.Time, error) {
	dateStrLower := strings.ToLower(dateStr)

	switch dateStrLower {
	case "today":
		return Today(), nil
	case "tomorrow":
		return Tomorrow(), nil
	default:
		for _, format := range DATE_FORMATS {
			if t, err := time.Parse(format, dateStr); err == nil {
				return t, nil
			}
		}

		return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
	}
}
