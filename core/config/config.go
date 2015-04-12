// Package config reads a configuration file and generates an instance of
// Configuration.  A description of the expected fields can be found in the BFF
// wiki.
package config

import (
	"fmt"
	"github.com/BruteForceFencer/core/hitcounter"
	"github.com/BruteForceFencer/core/logger"
	"github.com/BruteForceFencer/core/store"
	"os"
)

// Configuration is a struct that represents the contents of a configuration
// file.
type Configuration struct {
	Directions       []hitcounter.Direction
	ListenAddress    string
	ListenType       string
	DashboardAddress string
	AcceptedSources  []string
	Logger           *logger.Logger
}

// ReadConfig parses a configuration file and returns an instance of
// Configuration.
func ReadConfig(filename string) (*Configuration, []error) {
	parsed, err := parseJsonFile(filename)
	if err != nil {
		return nil, []error{err}
	}
	if errs := parsed.Validate(); len(errs) != 0 {
		return nil, errs
	}

	result := new(Configuration)
	result.ListenAddress = parsed.ListenAddress
	result.ListenType = parsed.ListenType
	result.DashboardAddress = parsed.DashboardAddress
	result.Directions = make([]hitcounter.Direction, 0, len(parsed.Directions))
	result.AcceptedSources = parsed.AcceptedSources

	// Logger
	if parsed.Log != "" {
		file, err := os.OpenFile(parsed.Log, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, []error{fmt.Errorf("unable to open file %s", parsed.Log)}
		}
		result.Logger = logger.New(file)
	}

	// Directions
	for _, jsonDir := range parsed.Directions {
		// Create the direction according to its type
		dir := hitcounter.Direction{
			Store:       store.NewShardMap(int64(jsonDir.MaxTracked)),
			Name:        jsonDir.Name,
			CleanUpTime: jsonDir.CleanUpTime,
			MaxHits:     jsonDir.MaxHits,
			WindowSize:  jsonDir.WindowSize,
		}
		dir.Store.Type = jsonDir.Typ

		// Add it to the list
		result.Directions = append(result.Directions, dir)
	}

	return result, nil
}
