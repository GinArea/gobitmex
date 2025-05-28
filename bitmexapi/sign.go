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
	signature, err := generateSignature(o.Secret, method, url, expiresStr, data)
	if err == nil {
		h.Set("api-key", o.Key)
		h.Set("api-expires", expiresStr)
		h.Set("api-signature", signature)
	}
}

func generateSignature(secret, verb, rawURL string, expiresStr string, data string) (string, error) {

	// Parse the URL and extract the path and query
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	path := parsedURL.Path
	if parsedURL.RawQuery != "" {
		path += "?" + parsedURL.RawQuery
	}

	// Construct the message to sign
	message := verb + path + expiresStr + data

	// Create the HMAC SHA256 signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	return signature, nil
}
