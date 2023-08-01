package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/apitypes"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func StartClients(ctx *cli.Context) (err error) {
	log.Info("🔎  Looking for client configuration file...")
	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  Couldn't read from config file: %v", err), 1)
	}

	selectedExecution := cfg.Execution()
	selectedConsensus := cfg.Consensus()
	if selectedExecution == "" || selectedConsensus == "" {
		return utils.Exit(errors.SelectedClientsNotFound, 1)
	}

	log.Info("🔄  Starting all clients")

	if ctx.Bool(flags.ValidatorFlag) && ctx.String(flags.TransactionFeeRecipientFlag) == "" || ctx.Bool(flags.TransactionFeeRecipientFlag) { // this means that we provided flag without value
		return utils.Exit(fmt.Sprintf("❌  %s flag is required but wasn't provided", flags.TransactionFeeRecipientFlag), 1)
	}

	executionClient, ok1 := clients.AllClients[selectedExecution]
	consensusClient, ok2 := clients.AllClients[selectedConsensus]

	if !ok1 || !ok2 {
		return utils.Exit("❌  Client found in LUKSO configuration file is not supported", 1)
	}

	execArgs, err := executionClient.PrepareStartFlags(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing %s flags: %v", executionClient.Name(), err), 1)
	}
	consArgs, err := consensusClient.PrepareStartFlags(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing %s flags: %v", consensusClient.Name(), err), 1)
	}

	if ctx.Bool(flags.CheckpointSyncFlag) && !ctx.Bool(flags.DevnetFlag) {
		log.Info("⚙️   Checkpoint sync feature enabled")

		network := "mainnet"
		if ctx.Bool(flags.TestnetFlag) {
			network = "testnet"
		}

		checkpointURL := fmt.Sprintf("https://checkpoints.%s.lukso.network", network)
		explorerURL := fmt.Sprintf("https://explorer.consensus.%s.lukso.network", network)

		var (
			root  string
			epoch int
		)

		root, epoch, err = getWeakSubjectivityCheckpoint(explorerURL)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while getting weak subjectivity checkpoint: %v", err), 1)
		}

		if root != "0x01" {
			switch consensusClient {
			case clients.Prysm:
				consArgs = append(consArgs, fmt.Sprintf("--checkpoint-sync-url=%s", checkpointURL))
				consArgs = append(consArgs, fmt.Sprintf("--genesis-beacon-api-url=%s", checkpointURL))
				consArgs = append(consArgs, fmt.Sprintf("--weak-subjectivity-checkpoint=%s:%d", root, epoch))
			case clients.Lighthouse:
				consArgs = append(consArgs, fmt.Sprintf("--checkpoint-sync-url=%s", checkpointURL))
				consArgs = append(consArgs, fmt.Sprintf("--wss-checkpoint=%s:%d", root, epoch))
			default:
				log.Warnf("️⚠️  Checkpoint sync not configured for %s: continuing without checkpoint sync", consensusClient.Name())
			}
		} else {
			log.Warn("️⚠️  Incorrect block root fetched - continuing without checkpoint sync")
		}
	}
	if ctx.Bool(flags.DevnetFlag) {
		log.Info("️️⚠️  This network doesn't have a checkpoint sync setup, starting without checkpoint sync...")
	}

	err = executionClient.Start(ctx, execArgs)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", executionClient.Name(), err), 1)
	}
	err = consensusClient.Start(ctx, consArgs)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", consensusClient.Name(), err), 1)
	}

	if ctx.Bool(flags.ValidatorFlag) {
		err = startValidator(ctx)
	}
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting validator: %v", err), 1)
	}

	log.Info("🎉  Clients have been started. Checking status:")
	log.Info("")

	_ = StatClients(ctx)

	log.Info("")
	log.Info("If execution and consensus clients are Running 🟢, your node is now 🆙.")
	log.Info("👉 Please check the logs with 'lukso logs' to make sure everything is running correctly.")

	return nil
}

func startValidator(ctx *cli.Context) (err error) {
	err = cfg.Read()
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while reading config: %v", err), 1)
	}

	var validatorClient clients.ValidatorBinaryDependency
	consensusClient, ok := clients.AllClients[cfg.Consensus()]
	if !ok {
		return utils.Exit(errors.ErrClientNotSupported.Error(), 1)
	}

	switch consensusClient {
	case clients.Prysm:
		if !utils.FileExists(fmt.Sprintf("%s/direct/accounts/all-accounts.keystore.json", ctx.String(flags.ValidatorKeysFlag))) { // path to imported keys
			return utils.Exit("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.", 1)
		}
		validatorClient = clients.PrysmValidator
	case clients.Lighthouse:
		if !utils.FileExists(fmt.Sprintf("%s/validators", ctx.String(flags.ValidatorKeysFlag))) { // path to imported keys
			return utils.Exit("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.", 1)
		}
		validatorClient = clients.LighthouseValidator
	case clients.Teku:
		if !utils.FileExists(ctx.String(flags.ValidatorKeysFlag)) { // path to imported keys
			return utils.Exit("⚠️  Validator is not initialized. Run lukso validator import to initialize your validator.", 1)
		}
		validatorClient = clients.TekuValidator
	}

	var passwordPipe *os.File = nil

	defer func() {
		if passwordPipe != nil {
			os.Remove(passwordPipe.Name())
		}
	}()

	validatorPasswordPath := ctx.String(flags.ValidatorWalletPasswordFileFlag)
	if validatorPasswordPath == "" && validatorClient == clients.PrysmValidator {
		passwordPipe, err = utils.ReadValidatorPassword(ctx)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while reading password: %v", err), 1)
		}
		err = ctx.Set(flags.ValidatorWalletPasswordFileFlag, passwordPipe.Name())
		if err != nil {
			return
		}
	}

	args, err := validatorClient.PrepareStartFlags(ctx)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while preparing %s flags: %v", validatorClient.Name(), err), 1)
	}

	err = validatorClient.Start(ctx, args)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while starting %s: %v", validatorClient.Name(), err), 1)
	}

	if validatorClient == clients.PrysmValidator {
		log.Info("⚙️  Please wait a few seconds while your password is being validated...")
		time.Sleep(time.Second * 10) // should be enough

		logFile, err := utils.GetLastFile(ctx.String(flags.LogFolderFlag), validatorClient.CommandName())
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while getting latest log file: %v", err), 1)
		}

		logs, err := os.ReadFile(ctx.String(flags.LogFolderFlag) + "/" + logFile)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while reading log file: %v", err), 1)
		}

		if strings.Contains(string(logs), errors.WrongPassword) {
			return utils.Exit("❌  Incorrect password, please restart and try again", 1)
		}
	}

	return
}

func getWeakSubjectivityCheckpoint(checkpointURL string) (finalizedRoot string, epoch int, err error) {
	checkpointURL = strings.TrimRight(checkpointURL, "/")

	var (
		res     *http.Response
		resByte []byte

		apiResp apitypes.ExplorerFinalizedSlotsResponse
	)

	res, err = http.Get(fmt.Sprintf("%s/api/v1/epoch/finalized/slots", checkpointURL))
	if err != nil {
		return
	}

	resByte, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(resByte, &apiResp)
	if err != nil {
		return
	}

	firstEpochSlot := apiResp.Data[0]
	epoch = firstEpochSlot.Epoch
	finalizedRoot = firstEpochSlot.BlockRoot

	return
}
