package bitmexapi

import (
	"fmt"
	"strings"
)

type Error struct {
	Name    string
	Message string
}

func (o *Error) Error() string {
	return fmt.Sprintf("name[%s]: %s", o.Name, o.Message)
}

func (o *Error) ApiKeyInvalid() bool {
	lowerCasedMessage := strings.ToLower(o.Message)

	invalid := false
	if strings.Contains(lowerCasedMessage, "invalid api") {
		//"Invalid API Key."
		invalid = true
	} else if strings.Contains(lowerCasedMessage, "key is disabled") {
		// "This key is disabled."
		invalid = true
	}
	return invalid
}

func (o *Error) UnmatchedIp() bool {
	lowerCasedMessage := strings.ToLower(o.Message)

	unmatched := false
	//"This IP address is not allowed to use this key."
	if strings.Contains(lowerCasedMessage, "ip address is not allowed") {
		unmatched = true
	}
	return unmatched
}

func (o *Error) InsufficientBalance() (insufficient bool) {
	lowerCasedMessage := strings.ToLower(o.Message)
	if strings.Contains(lowerCasedMessage, "insufficient available balance") {
		insufficient = true
	}
	return
}
