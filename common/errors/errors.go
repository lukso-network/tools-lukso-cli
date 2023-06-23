package errors

import "errors"

var (
	ErrNeedRoot             = errors.New("⚠️  You need root privileges to perform this action ")
	ErrFlagMissing          = errors.New("⚠️  Couldn't find given flag ")
	ErrMoreNetworksSelected = errors.New("⚠️  You can only specify 1 network ")
	ErrNotEnoughArguments   = errors.New("⚠️  Not enough arguments provided ")
	ErrProcessNotFound      = errors.New("⚠️  Process not found ")
	ErrFlagPathInvalid      = errors.New("⚠️  Invalid flag path ")
	ErrAlreadyRunning       = errors.New("⚠️  Process is already running ")
	ErrValidatorNotImported = errors.New("Validator has not been initialized - use 'lukso validator import' to initialize your validator ")
	ErrClientNotSupported   = errors.New("❌  Client found in LUKSO configuration file is not supported - if you think it is please contact LUKSO team")
	ErrGenesisNotFound      = errors.New("❌  Genesis JSON not found")
	ErrOlderFolderDetected  = errors.New(`❌  This node directory is not supported by your CLI version. To continue working with the CLI do the following:

1. 🪦  If your node clients are still running, please stop the related processes by executing 'sudo pkill geth' (Geth), 'sudo pkill prysm' (Prysm), 'sudo pkill validator' (Prysm Validator), 'sudo pkill erigon' (Erigon), or 'sudo pkill lighthouse' (Lighthouse).
2. 📁  Move into your home directory using 'cd' and create a new working directory, using 'mkdir myNewLUKSOnode && cd ./myNewLUKSOnode'
3. 🚀  Initialize the new folder using 'lukso init'
4. 🛠️  Re-Install your desired clients using 'lukso install'
5. 📦  Optional: Copy over your chaindata using 'cp -r ../myOldLUKSOnode/<network>-data .' if you dont want to resyncronize the full blockchain state
6. 🔑  Optional: If you are running a validator, import your validator keys or copy over your keystore files using 'cp -r ../myOldLUKSOnode/<network>-keystore .' Copying the keystore folder will only work in case you are using the same consensus client as before.`)
)

const (
	NoSuchFlag              = "no such flag" // no emoji here - this error should match the CLI lib error - we don't throw it to user anyway
	FolderNotInitialized    = "⚠️  Folder not initialized - please make sure that you are working in an initialized directory. You can initialize the directory with the 'lukso init' command."
	SelectedClientsNotFound = "⚠️  No selected client found in LUKSO configuration file. Please make sure that you have installed your clients. You can use the install command to install clients."
	WrongPassword           = "wrong password for wallet"
)
