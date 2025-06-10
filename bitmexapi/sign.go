package bitmexapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

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

func (o *Sign) HeaderGet(h http.Header, v url.Values, path string) {
	encodedParams := encodeSortParams(v)
	o.header(h, encodedParams, path, "GET")
}

func (o *Sign) HeaderPost(h http.Header, body []byte, path string) {
	o.header(h, string(body[:]), path, "POST")
}

func encodeSortParams(src url.Values) (s string) {
	if len(src) == 0 {
		return
	}
	keys := make([]string, len(src))
	i := 0
	for k := range src {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		s += encodeParam(k, src.Get(k)) + "&"
	}
	s = s[0 : len(s)-1]
	return
}

func encodeParam(name, value string) string {
	params := url.Values{}
	params.Add(name, value)
	return params.Encode()
}

func (o *Sign) header(h http.Header, data string, path string, method string) {

	// expiry - use Unix time + 5 seconds
	expires := time.Now().Unix() + 5
	expiresStr := strconv.FormatInt(expires, 10)
	url := "/" + ApiVersion + "/" + path
	signature := GenerateSignature(o.Secret, method, url, expiresStr, data)
	h.Set("api-key", o.Key)
	h.Set("api-expires", expiresStr)
	h.Set("api-signature", signature)
}

func (o *Sign) GetWsSignData() (signature, expires string) {
	e := time.Now().Unix() + 105 // 5 seconds not enough
	expires = strconv.FormatInt(e, 10)
	signature = GenerateSignature(o.Secret, "GET", "/realtime", expires, "")
	return
}

func GenerateSignature(secret, method, path string, expiresStr string, data string) string {
	message := method + path + expiresStr
	if data != "" {
		if method == "GET" {
			message = method + path + "?" + data + expiresStr
		} else if method == "POST" {
			message += data
		}
	}
	// Create the HMAC SHA256 signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return signature
}
