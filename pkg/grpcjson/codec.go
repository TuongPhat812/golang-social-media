package grpcjson

import (
	"encoding/json"
	"sync"

	"google.golang.org/grpc/encoding"
)

var (
	registerOnce  sync.Once
	codecInstance = &jsonCodec{}
)

func init() {
	registerOnce.Do(func() {
		encoding.RegisterCodec(codecInstance)
	})
}

func Codec() encoding.Codec {
	return codecInstance
}

type jsonCodec struct{}

func (c *jsonCodec) Name() string {
	return "json"
}

func (c *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c *jsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
