// package errors is responsible for holding error values used for clear error communication in API.
package errors

import "errors"

var (
	// Clients
	ErrClientRunning = errors.New("some client/s are already running")

	// File system
	ErrNeedsRoot  = errors.New("operation needs root access")
	ErrFileExists = errors.New("file already exists")
	ErrCfgExists  = errors.New("config already exists")

	// HTTP
	ErrNotFound = errors.New("not found")
	ErrNotOk    = errors.New("status code is not 2xx")
)
