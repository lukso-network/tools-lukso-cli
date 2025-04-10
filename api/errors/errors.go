// package errors is responsible for holding error values used for clear error communication in API.
package errors

import "errors"

var (
	ErrClientRunning = errors.New("some client/s are already running")
	ErrCfgExists     = errors.New("config already exists")
)
