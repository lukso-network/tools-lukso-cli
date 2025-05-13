package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common"
	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dep"
	"github.com/lukso-network/tools-lukso-cli/dep/clients"
	"github.com/lukso-network/tools-lukso-cli/flags"
)

func (c *commander) Install(ctx *cli.Context) (err error) {
	if clients.IsAnyRunning() {
		return
	}

	if !cfg.Exists() {
		return utils.Exit(errors.FolderNotInitialized, 1)
	}

	var (
		selectedConsensus dep.ConsensusClient
		selectedExecution dep.ExecutionClient
		consensusInput    string
		executionInput    string
		consensusTag      string
		executionTag      string
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
	case "4":
		selectedConsensus = clients.Nimbus2
		consensusTag = ctx.String(flags.Nimbus2TagFlag)

	}

	executionInput = utils.RegisterInputWithMessage(executionMessage)
	for executionInput != "1" && executionInput != "2" && executionInput != "3" && executionInput != "4" {
		executionInput = utils.RegisterInputWithMessage("Please provide a valid option\n> ")
	}

	switch executionInput {
	case "1":
		selectedExecution = clients.Geth
		executionTag = ctx.String(flags.GethTagFlag)
	case "2":
		selectedExecution = clients.Erigon
		executionTag = ctx.String(flags.ErigonTagFlag)
		// backwards compatibility - after erigon edited the repository it came with a few breaking changes
		// New tag format (_X.Y.Z_ -> _vX.Y.Z_) in GH release was introduced in v2.60.9
		tagVer := common.StringToSemver(executionTag)
		erigonVer := common.StringToSemver("2.60.8")
		if !tagVer.IsNewerThan(erigonVer) {
			log.Warn("⚠️  Erigon has introduced changes in making new release URLs.\n" +
				"This means that LUKSO CLI from now on will support versions from v2.60.9\n" +
				"If you wish to install older erigon versions (not recommended), please use <=0.22.0 version of LUKSO CLI.")
		}

	case "3":
		selectedExecution = clients.Nethermind
		executionTag = ctx.String(flags.NethermindTagFlag)
	case "4":
		selectedExecution = clients.Besu
		executionTag = ctx.String(flags.BesuTagFlag)
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

	if selectedExecution.Exists() {
		log.Infof("%s is already installed - continuing...", selectedExecution.Name())
	} else {
		log.Infof("⬇️  Downloading %s...", selectedExecution.Name())
		err = selectedExecution.Install(executionTag, false)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedExecution.Name(), err), 1)
		}
	}

	if selectedConsensus.Exists() {
		log.Infof("%s is already installed - continuing...", selectedConsensus.Name())
	} else {
		log.Infof("⬇️  Downloading %s...", selectedConsensus.Name())
		err = selectedConsensus.Install(consensusTag, false)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while downloading %s: %v", selectedConsensus.Name(), err), 1)
		}
	}

	var selectedValidator dep.ValidatorClient

	switch selectedConsensus {
	case clients.Prysm:
		selectedValidator = clients.PrysmValidator
	case clients.Lighthouse:
		selectedValidator = clients.LighthouseValidator
	case clients.Teku:
		selectedValidator = clients.TekuValidator
	case clients.Nimbus2:
		selectedValidator = clients.Nimbus2Validator
	}

	if selectedValidator.Exists() {
		log.Infof("%s is already installed - continuing...", selectedValidator.Name())
	} else {
		log.Infof("⬇️  Downloading %s...", selectedValidator.Name())
		err = selectedValidator.Install(consensusTag, false)
		if err != nil {
			return utils.Exit(fmt.Sprintf("❌  There was an error while downloading validator: %v", err), 1)
		}
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
