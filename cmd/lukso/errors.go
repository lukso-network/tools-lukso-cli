package main

import "errors"

var (
	errNeedRoot             = errors.New("⚠️  You need root privileges to perform this action ")
	errFlagMissing          = errors.New("⚠️  Couldn't find given flag ")
	errMoreNetworksSelected = errors.New("⚠️  You can only specify 1 network ")
	errNotEnoughArguments   = errors.New("⚠️  Not enough arguments provided ")
	errProcessNotFound      = errors.New("⚠️  Process not found ")
	errFlagPathInvalid      = errors.New("⚠️  Invalid flag path ")
	errAlreadyRunning       = errors.New("⚠️  Process is already running ")
	errValidatorNotImported = errors.New("Validator has not been initialized - use 'lukso validator import' to initialize your validator ")
)

const (
	noSuchFlag              = "no such flag" // no emoji here - this error should match the CLI lib error - we don't throw it to user anyway
	folderNotInitialized    = "⚠️  Folder not initialized - please make sure that you are working in an initialized directory. You can initialize the directory with the 'lukso init' command."
	selectedClientsNotFound = "⚠️  No selected client found in LUKSO configuration file. Please make sure that you have installed your clients. You can use the install command to install clients."
)
