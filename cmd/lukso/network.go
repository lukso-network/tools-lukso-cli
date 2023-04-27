package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// networkConfig serves as a collection of variables that need to be changed when different network is selected
type networkConfig struct {
	executionDatadirPath string
	consensusDatadirPath string
	validatorDatadirPath string
	logPath              string
	configPath           string
	keysPath             string
	walletPath           string
}

// selectNetwork accepts a CLI func as an argument, and adjusts all values that need to be changed depending on
// network passed as a flag. Works as a wrapper for selecting current working network
func selectNetworkFor(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.Command.SkipFlagParsing {
			log.Debug("Skipping flag parsing on - parsing flags manually...")
			err := parseFlags(ctx)
			if err != nil {
				return cli.Exit(fmt.Sprintf("❌  There was an error while parsing flags: %v", err), 1)
			}
		}

		testnetEnabled := ctx.Bool(testnetFlag)
		devnetEnabled := ctx.Bool(devnetFlag)

		enabledCount := boolToInt(testnetEnabled) + boolToInt(devnetEnabled)
		if enabledCount > 1 {
			return errMoreNetworksSelected
		}

		if enabledCount == 0 {
			return errNetworkNotSupported // when any other network is supported we can simply pass in the config there
		}

		var cfg networkConfig

		if devnetEnabled {
			cfg = networkConfig{
				executionDatadirPath: executionDevnetDatadir,
				consensusDatadirPath: consensusDevnetDatadir,
				validatorDatadirPath: validatorDevnetDatadir,
				logPath:              devnetLogs,
				configPath:           devnetConfig,
				keysPath:             devnetKeystore,
				walletPath:           devnetKeystore,
			}
		}

		if testnetEnabled {
			cfg = networkConfig{
				gethDatadirPath:      executionTestnetDatadir,
				prysmDatadirPath:     consensusTestnetDatadir,
				validatorDatadirPath: validatorTestnetDatadir,
				logPath:              testnetLogs,
				configPath:           testnetConfig,
				keysPath:             testnetKeystore,
				walletPath:           testnetKeystore,
			}
		}

		err := updateValues(ctx, cfg)
		if err != nil {
			return err
		}

		return f(ctx)
	}
}

// updateValues is responsible for overriding values for data and logs directories etc.
func updateValues(ctx *cli.Context, config networkConfig) (err error) {
	var (
		//genesisJson  = config.configPath + "/" + genesisJsonPath
		gethToml       = config.configPath + "/" + gethTomlPath
		erigonToml     = config.configPath + "/" + erigonTomlPath
		prysmYaml      = config.configPath + "/" + prysmYamlPath
		lighthouseYaml = config.configPath + "/" + lighthouseYamlPath
		validatorYaml  = config.configPath + "/" + validatorYamlPath
		gethGenesis    = config.configPath + "/" + genesisJsonPath
		genesisState   = config.configPath + "/" + genesisStateFilePath
		configYaml     = config.configPath + "/" + chainConfigYamlPath
	)

	passedArgs := make([]string, 0)

	for _, arg := range os.Args {
		if strings.Contains(arg, "--") {
			passedArgs = append(passedArgs, arg)
		}
	}

	// varyingFlags represents list of all flags that can be affected by selecting network and values that may be replaced
	varyingFlags := map[string]string{
		gethDatadirFlag:              config.executionDatadirPath,
		erigonDatadirFlag:            config.executionDatadirPath,
		prysmDatadirFlag:             config.consensusDatadirPath,
		validatorDatadirFlag:         config.validatorDatadirPath,
		logFolderFlag:                config.logPath,
		validatorKeysFlag:            config.keysPath,
		gethConfigFileFlag:           gethToml,
		erigonConfigFileFlag:         erigonToml,
		prysmConfigFileFlag:          prysmYaml,
		lighthouseConfigFileFlag:     lighthouseYaml,
		validatorConfigFileFlag:      validatorYaml,
		genesisJsonFlag:              gethGenesis,
		prysmChainConfigFileFlag:     configYaml,
		validatorChainConfigFileFlag: configYaml,
		prysmGenesisStateFlag:        genesisState,
		validatorWalletDirFlag:       config.walletPath,
	}

	if len(os.Args) < 2 {
		return cli.Exit(errNotEnoughArguments, 1)
	}

	for _, flag := range ctx.Command.Flags {
		names := flag.Names()
		if len(names) < 1 {
			return cli.Exit(errNotEnoughArguments, 1)
		}

		targetName := names[0]

		// search for flags that need to be changed during network selection
		for flagName, changedValue := range varyingFlags {
			if flagName == targetName {
				// BUT keep those that are passed in
				isPassed := false
				for _, passedValue := range passedArgs {
					if strings.Contains(passedValue, targetName) {
						isPassed = true
					}
				}
				if isPassed { // ignore
					continue
				}

				// at last, update value when flag that needs to be changed is found, and it isn't passed manually
				err = ctx.Set(targetName, changedValue)
				if err != nil {
					return
				}

				delete(varyingFlags, flagName)
				break
			}
		}
	}

	return nil
}
