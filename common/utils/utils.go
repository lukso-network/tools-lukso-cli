package utils

import (
	"math"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func Exit(message string, exitCode int) error {
	log.Error(message)

	os.Exit(exitCode)

	return nil // so we can return commands with this func
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func AcceptTermsInteractive() bool {
	message := "You are about to download the clients that are necessary to run the LUKSO Blockchain.\n" +
		"To install Prysm you are required to accept the terms of use:\n" +
		"https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md\n\n" +
		"Do you agree? [Y/n]: "

	input := RegisterInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Error("You need to type Y to continue.")
		return false
	}

	return true
}

// EthEpochToTimestamp converts given Ethereum Epoch to a timestamp of its occurence,
// assuming the chain started at given epochZeroTimestamp (in unix epoch format).
// It returns the converted time and bool, which is false when the time is invalid or too large.
func EthEpochToTimestamp(ethEpoch, epochZeroTimestamp uint64) (t time.Time, isValid bool) {
	if ethEpoch > math.MaxUint64-10 { // safety margin for any small shifts like -1 or -5
		isValid = false

		return
	}

	unixEpoch := ethEpoch * 32 * 12 // slots per epoch * seconds per slot
	t = time.Unix(int64(epochZeroTimestamp+unixEpoch), 0)
	isValid = true

	return
}
