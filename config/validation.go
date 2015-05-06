package config

import (
	"fmt"
	"runtime"
	"strings"
)

// Validate returns a list of validation errors, if any.
func (c *Configuration) Validate() error {
	// Required field
	if c.ListenAddress == "" {
		return fmt.Errorf("no listen address specified")
	}

	// Required field
	c.ListenType = strings.ToLower(c.ListenType)
	if c.ListenType != "unix" && c.ListenType != "tcp" {
		return fmt.Errorf("unknown listen type %s", c.ListenType)
	} else if c.ListenType == "unix" && runtime.GOOS == "windows" {
		return fmt.Errorf("listen type unix is not available on windows")
	}

	// Required fields
	if len(c.Directions) == 0 {
		return fmt.Errorf("no defined directions")
	}

	for i := range c.Directions {
		dir := &c.Directions[i]

		// Required field
		if dir.Name == "" {
			return fmt.Errorf("direction %d has no name", i)
		}

		// Required field
		dir.Typ = strings.ToLower(dir.Typ)
		if dir.Typ != "string" && dir.Typ != "int32" {
			return fmt.Errorf("unknown direction type %s", dir.Typ)
		}

		// Required field
		if dir.WindowSize == 0 {
			return fmt.Errorf("direction %s has no defined window size", dir.Name)
		} else if dir.WindowSize < 0 {
			return fmt.Errorf("direction %s has a negative window size of %f", dir.Name, dir.WindowSize)
		}

		// Required field
		if dir.MaxHits == 0 {
			return fmt.Errorf("direction %s has no defined max hits", dir.Name)
		} else if dir.MaxHits < 0 {
			return fmt.Errorf("direction %s has a negative max hits of %f", dir.Name, dir.MaxHits)
		}

		// Optional field
		if dir.CleanUpTime == 0 {
			dir.CleanUpTime = 5
		} else if dir.CleanUpTime < 0 {
			return fmt.Errorf("direction %s has a negative clean up time of %f", dir.Name, dir.CleanUpTime)
		}

		// Optional field
		if dir.MaxTracked < 0 {
			return fmt.Errorf("direction %s has a negative max tracked of %f", dir.MaxTracked)
		}
	}

	return nil
}
