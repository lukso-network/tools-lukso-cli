// package errors is responsible for holding error values used for clear error communication in API.
package errors

import "errors"

var ErrAlreadyRunning = errors.New("already running")
