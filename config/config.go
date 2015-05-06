// Package config reads a configuration file and generates an instance of
// Configuration.  A description of the expected fields can be found in the BFF
// wiki.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Configuration is a struct that mirrors the data as it should be found in
// the configuration file.
type Configuration struct {
	Directions       []Direction `json:"directions"`
	ListenAddress    string      `json:"listen address"`
	ListenType       string      `json:"listen type"`
	DashboardAddress string      `json:"dashboard address"`
	AcceptedSources  []string    `json:"accepted sources"`
	Log              string      `json:"log"`
}

// Direction is a struct that mirrors the direction objects as they should
// be found in the configuration file.
type Direction struct {
	Name        string  `json:"name"`
	Typ         string  `json:"type"`
	WindowSize  float64 `json:"window size"`
	MaxHits     float64 `json:"max hits"`
	CleanUpTime float64 `json:"clean up time"`
	MaxTracked  float64 `json:"max tracked"`
}

// ReadConfig parses a configuration file and returns an instance of
// Configuration.
func ReadConfig(filename string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open configuration file %s.", filename)
	}

	parsed := new(Configuration)
	if err := json.Unmarshal(data, parsed); err != nil {

		// We want to know if the error is from bad syntax or if it's from a
		// type mismatch so that we can provide useful errors.
		typeErr, ok := err.(*json.UnmarshalTypeError)
		if ok {
			return nil, fmt.Errorf(
				"configuration file has mismatched type; %s should be %s.",
				typeErr.Value,
				typeErr.Type,
			)
		} else {
			return nil, fmt.Errorf("can't parse configuration file.")
		}
	}

	return parsed, parsed.Validate()
}
