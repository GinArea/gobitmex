package bitmexapi

type Response[T any] struct {
	Data       T
	Limit      RateLimit
	Error      error
	StatusCode int // filled with server response
	NetError   bool
}

type response[T any] struct {
	Data          T
	ErrorResponse *responseError `json:"error,omitempty"`
}

type responseError struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

func (o *Response[T]) Ok() bool {
	return o.Error == nil
}

func (o *Response[T]) SetErrorIfNil(err error) {
	if o.Error == nil {
		o.Error = err
	}
}

func (o *response[T]) Error() error {
	if o.ErrorResponse != nil {
		return &Error{
			Name:    o.ErrorResponse.Name,
			Message: o.ErrorResponse.Message,
		}
	}
	return nil
}
