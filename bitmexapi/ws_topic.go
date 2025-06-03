package bitmexapi

import "encoding/json"

type Topic[T any] struct {
	Table  string
	Action string
	Data   T
}

type RawTopic Topic[json.RawMessage]

func UnmarshalRawTopic[T any](raw RawTopic) (ret Topic[T], err error) {
	ret.Table = raw.Table
	err = json.Unmarshal(raw.Data, &ret.Data)
	return
}
