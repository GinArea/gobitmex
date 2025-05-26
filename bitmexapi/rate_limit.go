package bitmexapi

type RateLimit struct {
	Limit     int `http:"X-Ratelimit-Limit"`     //  total amount
	Remaining int `http:"X-Ratelimit-Remaining"` // remaining
	// At the UNIX timestamp designated by x-ratelimit-reset, you will have enough requests left to retry your current request. If you have not exceeded your limit, this value is always the current timestamp.
	Reset int `http:"X-Ratelimit-Reset"`
}
