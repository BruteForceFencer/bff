package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// jsonConfiguration is a struct that mirrors the data as it should be found in
// the configuration file.
type jsonConfiguration struct {
	Directions       []jsonDirection `json:"directions"`
	ListenAddress    string          `json:"listen address"`
	ListenType       string          `json:"listen type"`
	DashboardAddress string          `json:"dashboard address"`
	AcceptedSources  []string        `json:"accepted sources"`
	Log              string          `json:"log"`
}

// jsonDirection is a struct that mirrors the direction objects as they should
// be found in the configuration file.
type jsonDirection struct {
	Name        string  `json:"name"`
	Typ         string  `json:"type"`
	WindowSize  float64 `json:"window size"`
	MaxHits     float64 `json:"max hits"`
	CleanUpTime float64 `json:"clean up time"`
	MaxTracked  float64 `json:"max tracked"`
}

// parseJsonFile will read the contents of a file and return its structure as a
// jsonConfiguration.
func parseJsonFile(filename string) (*jsonConfiguration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't read configuration file %s.", filename)
	}

	parsed := new(jsonConfiguration)
	err = json.Unmarshal(data, parsed)
	if err != nil {
		return nil, parseError(err)
	}

	return parsed, nil
}

// parseError returns a more descriptive error based on the return value of
// json.Unmarshal
func parseError(err error) error {
	typeErr, ok := err.(*json.UnmarshalTypeError)
	if ok {
		return fmt.Errorf("configuration file has mismatched type; %s should be %s", typeErr.Value, typeErr.Type)
	} else {
		return fmt.Errorf("can't parse configuration file")
	}
}
