package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

// networkConfig serves as a collection of variables that need to be changed when different network is selected
type networkConfig struct {
	gethDatadirPath        string
	prysmDatadirPath       string
	validatorDatadirPath   string
	gethGenesisDependency  string
	prysmGenesisDependency string
	prysmConfigDependency  string
	logPath                string
	configPath             string
	keystorePath           string
	walletPath             string
	jwtSecretPath          string
}

// selectNetworkFor accepts a CLI func as an argument, and adjusts all values that need to be changed depending on
// network passed as a flag. Works as a wrapper for selecting current working network
func selectNetworkFor(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		mainnetEnabled := ctx.Bool(mainnetEnabledFlag)
		testnetEnabled := ctx.Bool(testnetEnabledFlag)
		devnetEnabled := ctx.Bool(devnetEnabledFlag)

		enabledCount := boolToInt(mainnetEnabled) + boolToInt(testnetEnabled) + boolToInt(devnetEnabled)
		if enabledCount > 1 {
			return errMoreNetworksSelected
		}

		if enabledCount == 0 || testnetEnabled || mainnetEnabled {
			return errNetworkNotSupported
		}

		var cfg networkConfig

		if devnetEnabled {
			cfg = networkConfig{
				gethDatadirPath:        gethDevnetDatadir,
				prysmDatadirPath:       prysmDevnetDatadir,
				validatorDatadirPath:   validatorDevnetDatadir,
				gethGenesisDependency:  gethDevnetGenesisDependencyName,
				prysmGenesisDependency: prysmDevnetGenesisDependencyName,
				prysmConfigDependency:  prysmDevnetConfigDependencyName,
				logPath:                devnetLogs,
				configPath:             devnetConfig,
				keystorePath:           devnetKeystore,
				walletPath:             devnetKeystore,
			}
		}

		err := updateValues(ctx, cfg)
		if err != nil {
			return err
		}

		return f(ctx)
	}
}

func updateValues(ctx *cli.Context, config networkConfig) (err error) {
	var (
		//genesisJson  = config.configPath + "/" + genesisJsonPath
		genesisState = config.configPath + "/" + genesisStateFilePath
		configYaml   = config.configPath + "/" + configYamlPath
		jwtSecret    = config.configPath + "/" + jwtSecretPath
	)

	if len(os.Args) < 2 {
		return errNotEnoughArguments
	}

	commandName := os.Args[1]
	switch commandName {
	case startCommand:
		// -- CONFIGS --
		gethSelectedGenesis = config.gethGenesisDependency
		if err = ctx.Set(prysmGenesisStateFlag, genesisState); err != nil {
			return
		}
		if err = ctx.Set(prysmChainConfigFileFlag, configYaml); err != nil {
			return
		}

		// -- DATADIRS --
		if err = ctx.Set(gethDatadirFlag, config.gethDatadirPath); err != nil {
			return
		}
		if err = ctx.Set(prysmDatadirFlag, config.prysmDatadirPath); err != nil {
			return
		}
		if err = ctx.Set(validatorDatadirFlag, config.validatorDatadirPath); err != nil {
			return
		}

		// -- LOGS --
		if err = ctx.Set(gethLogDirFlag, config.logPath); err != nil {
			return
		}
		if err = ctx.Set(prysmLogDirFlag, config.logPath); err != nil {
			return
		}
		if err = ctx.Set(validatorLogDirFlag, config.logPath); err != nil {
			return
		}

		// -- JWT --
		if err = ctx.Set(prysmJWTSecretFlag, jwtSecret); err != nil {
			return
		}
		if err = ctx.Set(gethAuthJWTSecretFlag, jwtSecret); err != nil {
			return
		}

		// -- KEYSTORE
		if err = ctx.Set(validatorWalletDirFlag, config.keystorePath); err != nil {
			return
		}

	case initCommand:
		gethSelectedGenesis = config.gethGenesisDependency
		prysmSelectedGenesis = config.prysmGenesisDependency
		prysmSelectedConfig = config.prysmConfigDependency
		jwtPath = jwtSecret
	}

	return
}
