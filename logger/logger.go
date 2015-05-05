// Package logger implements an asynchronous logging mechanism.
package logger

import (
	"encoding/csv"
	"io"
	"time"
)

// Logger is an asynchronously writing logger.
type Logger struct {
	w *csv.Writer
	c chan []string
}

// New returns an initialized *Logger.
func New(w io.Writer) *Logger {
	result := new(Logger)
	result.w = csv.NewWriter(w)
	result.c = make(chan []string, 32)

	go func() {
		for {
			fields := <-result.c
			result.w.Write(fields)
			result.w.Flush()
		}
	}()

	return result
}

// Logs the slice of strings to the output in CSV format.
func (l *Logger) Log(s ...string) {
	const format = "2006-01-02 15:04:05"

	fields := make([]string, 1, len(s)+1)
	fields[0] = time.Now().Format(format)
	fields = append(fields, s...)
	l.c <- fields
}
