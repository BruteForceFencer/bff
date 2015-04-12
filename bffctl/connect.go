package main

import (
	"encoding/json"
	"errors"
	"net"
	"os"
)

var (
	ErrNoConfig         = errors.New("unable to open last used config file")
	ErrInvalidConfig    = errors.New("the last used config file is corrupted")
	ErrClientNotRunning = errors.New("the client is not running")
)

func connect() (net.Conn, error) {
	file, err := os.Open(lastConfig)
	if err != nil {
		return nil, ErrNoConfig
	}
	defer file.Close()

	type config struct {
		ListenAddress string `json:"listen address"`
		ListenType    string `json:"listen type"`
	}

	c := new(config)
	dec := json.NewDecoder(file)
	if err := dec.Decode(c); err != nil {
		return nil, ErrInvalidConfig
	}

	conn, err := net.Dial(c.ListenType, c.ListenAddress)
	if err != nil {
		return nil, ErrClientNotRunning
	}

	return conn, nil
}
