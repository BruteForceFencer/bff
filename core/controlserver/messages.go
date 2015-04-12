package controlserver

import (
	"encoding/json"
	"io"
)

type RequestType int8

const (
	HitRequest RequestType = iota
	CommandRequest
)

// Request is a struct that follows the format that is expected for requests.
type Request struct {
	Type      RequestType
	Direction string
	Value     interface{}
}

// ReadRequest reads a request from a Reader.
func ReadRequest(r io.Reader) (*Request, error) {
	decoder := json.NewDecoder(r)
	result := new(Request)
	err := decoder.Decode(result)
	return result, err
}
