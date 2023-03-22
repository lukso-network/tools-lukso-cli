package main

import "errors"

var (
	errNeedRoot                = errors.New("You need root privileges to perform this action ")
	errFlagMissing             = errors.New("Couldn't find given flag ")
	errMoreNetworksSelected    = errors.New("You can only specify 1 network ")
	errNotEnoughArguments      = errors.New("Not enough arguments provided ")
	errNetworkNotSupported     = errors.New("Selected network is not supported ")
	errDepositNotProvided      = errors.New("You need to provide a deposit data to send your deposit ")
	errKeysNotProvided         = errors.New("You need to provide a path to your keys directory ")
	errTooManyDepositsProvided = errors.New("You can only provide 1 deposit data file ")
	errAlreadyRunning          = errors.New("Process is already running ")
	errProcessNotFound         = errors.New("Process not found ")
)

const noSuchFlag = "no such flag"
