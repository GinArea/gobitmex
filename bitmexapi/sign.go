package bitmexapi

type Sign struct {
	Key    string
	Secret string
}

func NewSign(key, secret string) *Sign {
	o := new(Sign)
	o.Key = key
	o.Secret = secret
	return o
}
