package commands

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/system"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func InstallBinaries(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	isRoot, err := system.IsRoot()
	if err != nil {
		return utils.Exit(fmt.Sprintf("There was an error while checking user privileges: %v", err), 1)
	}
	if !isRoot {
		return utils.Exit(errors.ErrNeedRoot.Error(), 1)
	}

	var (
		selectedConsensus   clients.ClientBinaryDependency
		selectedExecution   clients.ClientBinaryDependency
		consensusInput      string
		executionInput      string
		consensusTag        string
		executionTag        string
		executionCommitHash string
		consensusCommitHash string
		isSetupClientsDir   bool = false
	)

	consensusMessage := "\nWhich consensus client do you want to install?\n" +
		"1: prysm\n" +
		"2: lighthouse\n" +
		"3: teku\n" +
		"4: nimbus (eth2)\n" +
		"> "

	executionMessage := "\nWhich execution client do you want to install?\n" +
		"1: geth\n" +
		"2: erigon\n" +
		"3: nethermind\n" +
		"4: besu\n" +
		"> "

	consensusInput = utils.RegisterInputWithMessage(consensusMessage)
	for consensusInput != "1" && consensusInput != "2" && consensusInput != "3" && consensusInput != "4" {
		consensusInput = utils.RegisterInputWithMessage("Please provide a valid option\n> ")
	}

	switch consensusInput {
	case "1":
		selectedConsensus = clients.Prysm
		consensusTag = ctx.String(flags.PrysmTagFlag)
	case "2":
		selectedConsensus = clients.Lighthouse
		consensusTag = ctx.String(flags.LighthouseTagFlag)
	case "3":
		selectedConsensus = clients.Teku
		consensusTag = ctx.String(flags.TekuTagFlag)
		isSetupClientsDir = true
	case "4":
		selectedConsensus = clients.Nimbus2
		consensusTag = ctx.String(flags.Nimbus2TagFlag)
		consensusCommitHash = ctx.String(flags.Nimbus2CommitHashFlag)
		isSetupClientsDir = true

	}

	executionInput = utils.RegisterInputWithMessage(executionMessage)
	for executionInput != "1" && executionInput != "2" && executionInput != "3" && executionInput != "4" {
		executionInput = utils.RegisterInputWithMessage("Please provide a valid option\n> ")
	}

	switch executionInput {
	case "1":
		selectedExecution = clients.Geth
		executionTag = ctx.String(flags.GethTagFlag)
		executionCommitHash = ctx.String(flags.GethCommitHashFlag)
	case "2":
		selectedExecution = clients.Erigon
		executionTag = ctx.String(flags.ErigonTagFlag)
	case "3":
		selectedExecution = clients.Nethermind
		executionTag = ctx.String(flags.NethermindTagFlag)
		executionCommitHash = ctx.String(flags.NethermindCommitHashFlag)
		isSetupClientsDir = true
	case "4":
		selectedExecution = clients.Besu
		executionTag = ctx.String(flags.BesuTagFlag)
		isSetupClientsDir = true
	}

	if selectedConsensus == clients.Prysm {
		termsAgreed := ctx.Bool(flags.AgreeTermsFlag)
		if !termsAgreed {
			fmt.Println("")
			accepted := utils.AcceptTermsInteractive()
			if !accepted {
				return utils.Exit("❌  Terms of use not accepted - aborting...", 1)
			}
		} else {
			log.Info("✅  You accepted Prysm's Terms of Use: https://github.com/prysmaticlabs/prysm/blob/develop/TERMS_OF_SERVICE.md")
		}
	}

	if isSetupClientsDir {
		log.Infof("⚙️   Preparing clients directory")

		err = setupClientsDir()
		if err != nil {
			log.Warnf("⚠️  There was an error while creating clients directory: %v", err)
		} else {
			log.Infof("✅  Clients directory prepared successfully")
		}
	}

	log.Infof("⬇️  Downloading %s...", selectedExecution.Name())
	err = selectedExecution.Install(selectedExecution.ParseUrl(executionTag, executionCommitHash), false)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedExecution.Name(), err), 1)
	}

	log.Infof("⬇️  Downloading %s...", selectedConsensus.Name())
	err = selectedConsensus.Install(selectedConsensus.ParseUrl(consensusTag, consensusCommitHash), false)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedConsensus.Name(), err), 1)
	}

	var selectedValidator clients.ValidatorBinaryDependency = clients.LighthouseValidator

	if selectedConsensus == clients.Prysm {
		selectedValidator = clients.PrysmValidator
		log.Infof("⬇️  Downloading %s...", selectedValidator.Name())
		err = selectedValidator.Install(selectedValidator.ParseUrl(consensusTag, ""), false)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while downloading validator: %v", err), 1)
		}
	}
	if selectedConsensus == clients.Teku {
		selectedValidator = clients.TekuValidator
	}
	if selectedConsensus == clients.Nimbus2 {
		selectedValidator = clients.Nimbus2Validator
	}

	err = cfg.WriteExecution(selectedExecution.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing execution client: %v", err), 1)
	}

	err = cfg.WriteConsensus(selectedConsensus.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing consensus client: %v", err), 1)
	}

	err = cfg.WriteValidator(selectedValidator.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while writing validator client: %v", err), 1)
	}

	log.Info("✅  Configuration files created!")
	log.Info("✅  Clients have been successfully installed.")
	log.Info("➡️  Start your node using 'lukso start'")

	return
}

func setupClientsDir() (err error) {
	var f os.FileInfo

	f, err = os.Stat(common.ClientDepsFolder)
	if err != nil {
		err = os.Mkdir(common.ClientDepsFolder, os.ModePerm)

		return
	}

	if f.Mode() != os.ModePerm {
		err = os.Chmod(common.ClientDepsFolder, os.ModePerm)
	}

	return
}
