package main

import (
	"google.golang.org/grpc/encoding"
)

type codec struct{}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	b := v.([]byte)
	return b, nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	b := v.(*[]byte)
	*b = data
	return nil
}

func (c codec) Name() string {
	return "custom"
}

func init() {
	encoding.RegisterCodec(codec{})
}

func main() {
	run()
}
