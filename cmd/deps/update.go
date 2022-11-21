package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"strings"
)

const (
	gethGithubLocation          = "ethereum/go-ethereum"
	prysmaticLabsGithubLocation = "prysmaticlabs/prysm"
)

// fetchReleaseFor fetches the newest release for given dependency
func fetchReleaseFor(githubLocation string) (newestTag string, err error) {
	latestReleaseUrl := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", githubLocation)

	response, err := http.Get(latestReleaseUrl)
	if err != nil {
		return "", err
	}

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var latestReleaseResponse GithubApiReleaseResponse

	err = json.Unmarshal(respBytes, &latestReleaseResponse)
	if err != nil {
		return "", err
	}

	return latestReleaseResponse.TagName, nil
}

func updateClients(ctx *cli.Context) error {
	err := updateGeth(ctx)
	if err != nil {
		return err
	}

	err = updatePrysm(ctx)
	if err != nil {
		return err
	}

	err = updateValidator(ctx)

	return err
}

func updateGeth(ctx *cli.Context) error {
	log.Info("Fetching latest release for Geth")
	latestGethTag, err := fetchReleaseFor(gethGithubLocation)
	if err != nil {
		return err
	}

	// since geth needs standard x.y.z semantic version to download (without "v" at the beginning) we need to strip it
	strippedTag := strings.TrimPrefix(latestGethTag, "v")

	log.WithField("dependencyTag", latestGethTag).Info("Updating Geth")

	return clientDependencies[gethDependencyName].Download(strippedTag, binDir, "e5eb32ac")
}

func updatePrysm(ctx *cli.Context) error {
	log.Info("Fetching latest release for Geth")
	latestPrysmTag, err := fetchReleaseFor(gethGithubLocation)
	if err != nil {
		return err
	}

	log.WithField("dependencyTag", latestPrysmTag).Info("Updating Prysm")

	return clientDependencies[prysmDependencyName].Download(latestPrysmTag, binDir, "")
}

func updateValidator(ctx *cli.Context) error {
	log.Info("Fetching latest release for Geth")
	latestValidatorFlag, err := fetchReleaseFor(gethGithubLocation)
	if err != nil {
		return err
	}

	log.WithField("dependencyTag", latestValidatorFlag).Info("Updating Validator")

	return clientDependencies[validatorDependencyName].Download(latestValidatorFlag, binDir, "")
}

func beforeUpdate(ctx *cli.Context) error {
	// Geth related parsing
	gethTag = ctx.String(gethTagFlag)
	gethCommitHash = ctx.String(gethCommitHashFlag)

	// Validator related parsing
	validatorTag = ctx.String(validatorTagFlag)

	// Prysm related parsing
	prysmTag = ctx.String(prysmTagFlag)

	return nil
}
