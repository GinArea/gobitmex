package bitmexapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/msw-x/moon/ufmt"
	"github.com/msw-x/moon/uhttp"
)

func GetPub[R, T any](c *Client, path string, req any, transform func(R) (T, error)) Response[T] {
	return request(c, http.MethodGet, path, req, transform, false)
}

func Get[R, T any](c *Client, path string, req any, transform func(R) (T, error)) Response[T] {
	return request(c, http.MethodGet, path, req, transform, true)
}

func Post[R, T any](c *Client, path string, req any, transform func(R) (T, error)) Response[T] {
	return request(c, http.MethodPost, path, req, transform, true)
}

func request[R, T any](c *Client, method string, path string, request any, transform func(R) (T, error), sign bool) (r Response[T]) {
	var attempt int
	for {
		r = req(c, method, path, request, transform, sign)
		if r.StatusCode != http.StatusOK && c.onTransportError != nil {
			if c.onTransportError(r.Error, method, r.StatusCode, attempt) {
				attempt++
				continue
			}
		}
		break
	}
	return
}

func req[R, T any](c *Client, method string, path string, request any, transform func(R) (T, error), sign bool) (r Response[T]) {
	var perf *uhttp.Performer
	switch method {
	case http.MethodGet:
		perf = c.c.Get(path).Params(request)
	case http.MethodPost:
		perf = c.c.Post(path).Json(request)
	default:
		r.Error = fmt.Errorf("forbidden method: %s", method)
		return
	}
	if sign && c.s != nil {
		if perf.Request.Header == nil {
			perf.Request.Header = make(http.Header)
		}
		switch method {
		case http.MethodGet:
			c.s.HeaderGet(perf.Request.Header, perf.Request.Params, path)
		case http.MethodPost:
			c.s.HeaderPost(perf.Request.Header, perf.Request.Body, path)
		}
	}
	httpResponse := perf.Do()
	if httpResponse.Error == nil {
		r.StatusCode = httpResponse.StatusCode
		// if httpResponse.BodyExists() {
		// 	fmt.Println(string(httpResponse.Body))
		// }
		if httpResponse.BodyExists() &&
			r.StatusCode != http.StatusBadGateway && // 502
			r.StatusCode != http.StatusServiceUnavailable && // 503
			r.StatusCode != http.StatusGatewayTimeout &&
			/*
				Cloudflare is a content delivery network that acts as a gateway between a user and a website server. When the 530 status code is received, it will be accompanied by a more detailed HTTP status code 1XXX error message. The errors vary in meaning and severity and are situation-dependent.

				For example, Cloudflare HTTP status code 530 / Error 1016 Origin DNS Error is a case where HTTP status code 530 was returned, with the secondary HTTP status code 1016 Origin DNS Error further specified. In this very specific case, Cloudflare was unable to resolve the origin web server’s IP address due to a DNS error.
			*/
			r.StatusCode != 530 {
			resp := new(response[R])
			r.Error = resp.parseJsonAndFillResponse(httpResponse) // If returns an error, then it is precisely a JSON decoding error.
			if r.Ok() {
				r.Error = resp.Error() // here search exchange error handling
				if r.Ok() {
					r.Data, r.Error = transform(resp.Data)
				}
			}
		} else {
			r.Error = errors.New(ufmt.Join(httpResponse.Status))
		}
		if sign {
			r.SetErrorIfNil(httpResponse.HeaderTo(&r.Limit))
		}
	} else {
		r.Error = httpResponse.Error
		r.NetError = true
	}
	return
}

func (r *response[T]) parseJsonAndFillResponse(data uhttp.Response) error {

	/*
	  If the function returns an error, then it is precisely a JSON decoding error.
	*/

	var errCheck struct {
		Error *responseError `json:"error"`
	}
	err := data.Json(&errCheck)
	if err == nil && errCheck.Error != nil {
		r.ErrorResponse = errCheck.Error
		return nil
	}

	result := new(T)
	//log.Println(string(data.Body))
	if err := data.Json(result); err == nil {
		r.Data = *result
		return nil
	}

	// Fallback: try to parse {"message": "..."} format
	var msgOnly struct {
		Message string `json:"message"`
	}
	if err := data.Json(&msgOnly); err == nil && msgOnly.Message != "" {
		r.ErrorResponse = &responseError{
			Message: msgOnly.Message,
			Name:    "MessageOnly",
		}
		return nil
	}

	// Final fallback: return last JSON error
	return fmt.Errorf("unrecognized response format or JSON parse error")
}
