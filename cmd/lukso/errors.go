package main

import "errors"

var (
	errNeedRoot             = errors.New("You need root privilages to perform this action ")
	errFlagMissing          = errors.New("Couldn't find given flag ")
	errMoreNetworksSelected = errors.New("You can only specify 1 network ")
	errNotEnoughArguments   = errors.New("Not enough arguments provided ")
	errNetworkNotSupported  = errors.New("Selected network is not supported ")
)
