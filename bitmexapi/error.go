package bitmexapi

import "fmt"

type Error struct {
	Name    string
	Message string
}

func (o *Error) Error() string {
	return fmt.Sprintf("name[%s]: %s", o.Name, o.Message)
}
