package config

import (
	"fmt"
	"runtime"
	"strings"
)

// Validate returns a list of validation errors, if any.
func (j *jsonConfiguration) Validate() []error {
	result := make([]error, 0)

	// Required field
	if j.ListenAddress == "" {
		result = append(result, fmt.Errorf("no listen address specified"))
	}

	// Required field
	j.ListenType = strings.ToLower(j.ListenType)
	if j.ListenType != "unix" && j.ListenType != "tcp" {
		result = append(result, fmt.Errorf("unknown listen type %s", j.ListenType))
	} else if j.ListenType == "unix" && runtime.GOOS == "windows" {
		result = append(result, fmt.Errorf("listen type unix is not available on windows"))
	}

	// Required fields
	if len(j.Directions) == 0 {
		result = append(result, fmt.Errorf("no defined directions"))
	}

	for i := range j.Directions {
		dir := &j.Directions[i]

		// Required field
		if dir.Name == "" {
			result = append(result, fmt.Errorf("direction %d has no name", i))
		}

		// Required field
		dir.Typ = strings.ToLower(dir.Typ)
		if dir.Typ != "string" && dir.Typ != "int32" {
			result = append(result, fmt.Errorf("unknown direction type %s", dir.Typ))
		}

		// Required field
		if dir.WindowSize == 0 {
			result = append(result, fmt.Errorf("direction %s has no defined window size", dir.Name))
		} else if dir.WindowSize < 0 {
			result = append(result, fmt.Errorf("direction %s has a negative window size of %f", dir.Name, dir.WindowSize))
		}

		// Required field
		if dir.MaxHits == 0 {
			result = append(result, fmt.Errorf("direction %s has no defined max hits", dir.Name))
		} else if dir.MaxHits < 0 {
			result = append(result, fmt.Errorf("direction %s has a negative max hits of %f", dir.Name, dir.MaxHits))
		}

		// Optional field
		if dir.CleanUpTime == 0 {
			dir.CleanUpTime = 5
		} else if dir.CleanUpTime < 0 {
			result = append(result, fmt.Errorf("direction %s has a negative clean up time of %f", dir.Name, dir.CleanUpTime))
		}

		// Optional field
		if dir.MaxTracked < 0 {
			result = append(result, fmt.Errorf("direction %s has a negative max tracked of %f", dir.MaxTracked))
		}
	}

	return result
}
