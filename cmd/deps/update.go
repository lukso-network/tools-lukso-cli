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

// fetchTag fetches the newest release tag for given dependency from GitHub API
func fetchTag(githubLocation string) (string, error) {
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

// fetchTagAndCommitHash fetches both release and latest commit hash from GitHub API
func fetchTagAndCommitHash(githubLocation string) (releaseTag, commitHash string, err error) {
	latestTag, err := fetchTag(githubLocation)
	if err != nil {
		return
	}

	releaseTag = latestTag

	latestCommitUrl := fmt.Sprintf("https://api.github.com/repos/%s/git/ref/tags/%s", githubLocation, latestTag)

	response, err := http.Get(latestCommitUrl)
	if err != nil {
		return
	}

	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	var latestCommitResponse GithubApiCommitResponse

	err = json.Unmarshal(respBytes, &latestCommitResponse)
	if err != nil {
		return
	}

	commitHash = latestCommitResponse.Object.Sha

	return
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

	latestGethTag, latestGethCommitHash, err := fetchTagAndCommitHash(gethGithubLocation)
	if err != nil {
		return err
	}

	log.Infof("Fetched latest release: %s", latestGethTag)

	// since geth needs standard x.y.z semantic version to download (without "v" at the beginning) we need to strip it
	strippedTag := strings.TrimPrefix(latestGethTag, "v")
	// and we need to take only 8 first characters from commit hash
	chars := []rune(latestGethCommitHash)[:8]
	latestGethCommitHash = string(chars)

	log.WithField("dependencyTag", latestGethTag).Info("Updating Geth")

	return clientDependencies[gethDependencyName].Download(strippedTag, binDir, latestGethCommitHash, false)
}

func updatePrysm(ctx *cli.Context) error {
	log.Info("Fetching latest release for Prysm")

	latestPrysmTag, err := fetchTag(prysmaticLabsGithubLocation)
	if err != nil {
		return err
	}

	log.Infof("Fetched latest release: %s", latestPrysmTag)

	log.WithField("dependencyTag", latestPrysmTag).Info("Updating Prysm")

	return clientDependencies[prysmDependencyName].Download(latestPrysmTag, binDir, "", false)
}

func updateValidator(ctx *cli.Context) error {
	log.Info("Fetching latest release for Validator")

	latestValidatorFlag, err := fetchTag(prysmaticLabsGithubLocation)
	if err != nil {
		return err
	}

	log.Infof("Fetched latest release: %s", latestValidatorFlag)

	log.WithField("dependencyTag", latestValidatorFlag).Info("Updating Validator")

	return clientDependencies[validatorDependencyName].Download(latestValidatorFlag, binDir, "", false)
}

func updateGethToSpec(ctx *cli.Context) error {
	log.WithField("dependencyTag", gethTag).Info("Updating Geth")

	return clientDependencies[gethDependencyName].Download(gethTag, binDir, gethCommitHash, false)
}

func updatePrysmToSpec(ctx *cli.Context) error {
	log.WithField("dependencyTag", prysmTag).Info("Updating Prysm")

	return clientDependencies[prysmDependencyName].Download(prysmTag, binDir, "", false)
}

func updateValidatorToSpec(ctx *cli.Context) error {
	log.WithField("dependencyTag", validatorTag).Info("Updating Validator")

	return clientDependencies[validatorDependencyName].Download(validatorTag, binDir, "", false)
}
