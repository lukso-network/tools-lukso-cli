package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

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
		selectedConsensus clients.ClientBinaryDependency
		selectedExecution clients.ClientBinaryDependency
		consensusInput    string
		executionInput    string
		consensusTag      string
		executionTag      string
		commitHash        string
	)

	consensusMessage := "\nWhich consensus client do you want to install?\n" +
		"1: prysm\n2: lighthouse\n> "

	executionMessage := "\nWhich execution client do you want to install?\n" +
		"1: geth\n2: erigon\n> "

	consensusInput = utils.RegisterInputWithMessage(consensusMessage)
	for consensusInput != "1" && consensusInput != "2" {
		consensusInput = utils.RegisterInputWithMessage("Please provide a valid option\n> ")
	}

	switch consensusInput {
	case "1":
		selectedConsensus = clients.Prysm
		consensusTag = ctx.String(flags.PrysmTagFlag)
	case "2":
		selectedConsensus = clients.Lighthouse
		consensusTag = ctx.String(flags.LighthouseTagFlag)
	}

	executionInput = utils.RegisterInputWithMessage(executionMessage)
	for executionInput != "1" && executionInput != "2" {
		executionInput = utils.RegisterInputWithMessage("Please provide a valid option\n> ")
	}

	switch executionInput {
	case "1":
		selectedExecution = clients.Geth
		executionTag = ctx.String(flags.GethTagFlag)
		commitHash = ctx.String(flags.GethCommitHashFlag)
	case "2":
		selectedExecution = clients.Erigon
		executionTag = ctx.String(flags.ErigonTagFlag)
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

	log.Infof("⬇️  Downloading %s...", selectedExecution.Name())
	err = selectedExecution.Install(selectedExecution.ParseUrl(executionTag, commitHash), false)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedExecution.Name(), err), 1)
	}

	log.Infof("⬇️  Downloading %s...", selectedConsensus.Name())
	err = selectedConsensus.Install(selectedConsensus.ParseUrl(consensusTag, commitHash), false)
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedConsensus.Name(), err), 1)
	}

	var selectedValidator clients.ValidatorBinaryDependency = clients.LighthouseValidator

	if selectedConsensus == clients.Prysm {
		selectedValidator = clients.PrysmValidator
		log.Infof("⬇️  Downloading %s...", clients.PrysmValidator.Name())
		err = clients.PrysmValidator.Install(clients.PrysmValidator.ParseUrl(consensusTag, ""), false)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while downloading validator: %v", err), 1)
		}
	}

	err = cfg.Create(selectedExecution.Name(), selectedConsensus.Name(), selectedValidator.Name())
	if err != nil {
		return utils.Exit(fmt.Sprintf("❌  There was an error while creating configration file: %v", err), 1)
	}

	log.Info("✅  Configuration files created!")
	log.Info("✅  Clients have been successfully installed.")
	log.Info("➡️  Start your node using 'lukso start'")

	return
}
