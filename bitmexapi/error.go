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

func (o *Error) KycNeed() bool {
	lowerCasedMessage := strings.ToLower(o.Message)
	kycNeed := false
	if strings.Contains(lowerCasedMessage, "new traders to verify") {
		/*
					{
						"error":{
			      			"message":"We require all new traders to verify theirnidentity before their first deposit. Please visit bitmex.com/verify to complete the process.",
			      			"name":"HTTPError"
			   			}
					}
		*/
		kycNeed = true
	}
	return kycNeed
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
	} else if strings.Contains(lowerCasedMessage, "invalid use of subaccount api key") {
		// subAccount keys
		invalid = true
	} else if strings.Contains(lowerCasedMessage, "signature not valid") {
		/*
			 // bad secret

						{
							"error":{
				      			"message":"Signature not valid.",
				      			"name":"HTTPError"
				   			}
						}

		*/
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

func (o *Error) Timeout() (timeout bool) {

	/*
	   not real example, AI helped
	*/

	lowerCasedMessage := strings.ToLower(o.Message)
	if strings.Contains(lowerCasedMessage, "gateway timeout") ||
		strings.Contains(lowerCasedMessage, "bad gateway") ||
		strings.Contains(lowerCasedMessage, "system is currently overloaded") ||
		strings.Contains(lowerCasedMessage, "request timed out") {
		timeout = true
	}
	return
}

func (o *Error) Restricted() (restricted bool) {
	lowerCasedMessage := strings.ToLower(o.Message)
	if strings.Contains(lowerCasedMessage, "restricted") {
		restricted = true
	}
	return
}
