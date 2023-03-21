package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

// networkConfig serves as a collection of variables that need to be changed when different network is selected
type networkConfig struct {
	gethDatadirPath        string
	prysmDatadirPath       string
	validatorDatadirPath   string
	gethGenesisDependency  string
	gethConfigDependency   string
	prysmGenesisDependency string
	prysmConfigDependency  string
	logPath                string
	configPath             string
	keysPath               string
}

// selectNetwork accepts a CLI func as an argument, and adjusts all values that need to be changed depending on
// network passed as a flag. Works as a wrapper for selecting current working network
func selectNetworkFor(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.Command.SkipFlagParsing {
			fmt.Println("ASDASDASD")
			log.Debug("Skipping flag parsing on - parsing flags manually...")
			err := parseFlags(ctx)
			if err != nil {
				return err
			}
		}

		fmt.Println(ctx.Args())

		mainnetEnabled := ctx.Bool(mainnetFlag)
		testnetEnabled := ctx.Bool(testnetFlag)
		devnetEnabled := ctx.Bool(devnetFlag)

		enabledCount := boolToInt(mainnetEnabled) + boolToInt(testnetEnabled) + boolToInt(devnetEnabled)
		if enabledCount > 1 {
			return errMoreNetworksSelected
		}

		if enabledCount == 0 || testnetEnabled || mainnetEnabled {
			return errNetworkNotSupported // when any other network is supported we can simply pass in the config there
		}

		var cfg networkConfig

		if devnetEnabled {
			cfg = networkConfig{
				gethDatadirPath:        gethDevnetDatadir,
				prysmDatadirPath:       prysmDevnetDatadir,
				validatorDatadirPath:   validatorDevnetDatadir,
				gethGenesisDependency:  gethDevnetGenesisDependencyName,
				gethConfigDependency:   gethDevnetConfigName,
				prysmGenesisDependency: prysmDevnetGenesisDependencyName,
				prysmConfigDependency:  prysmDevnetConfigDependencyName,
				logPath:                devnetLogs,
				configPath:             devnetConfig,
				keysPath:               devnetKeystore,
			}
		}

		err := updateValues(ctx, cfg)
		if err != nil {
			return err
		}

		return f(ctx)
	}
}

// updateValues is responsible for overriding values for data dirs, log dirs etc.
func updateValues(ctx *cli.Context, config networkConfig) (err error) {
	var (
		//genesisJson  = config.configPath + "/" + genesisJsonPath
		gethToml     = config.configPath + "/" + configTomlPath
		gethGenesis  = config.configPath + "/" + genesisJsonPath
		genesisState = config.configPath + "/" + genesisStateFilePath
		configYaml   = config.configPath + "/" + configYamlPath
		jwtSecret    = config.configPath + "/" + jwtSecretPath
	)

	passedArgs := make([]string, 0)

	for _, arg := range os.Args {
		if strings.Contains(arg, "--") {
			passedArgs = append(passedArgs, arg)
		}
	}

	// selecting dependencies for init
	gethSelectedGenesis = config.gethGenesisDependency
	gethSelectedConfig = config.gethConfigDependency
	prysmSelectedGenesis = config.prysmGenesisDependency
	prysmSelectedConfig = config.prysmConfigDependency
	jwtSelectedPath = jwtSecret

	// varyingFlags represents list of all flags that can be affected by selecting network and values that may be replaced
	varyingFlags := map[string]string{
		gethDatadirFlag:          config.gethDatadirPath,
		prysmDatadirFlag:         config.prysmDatadirPath,
		validatorDatadirFlag:     config.validatorDatadirPath,
		logFolderFlag:            config.logPath,
		validatorKeysFlag:        config.keysPath,
		jwtSecretFlag:            jwtSecret,
		gethConfigFileFlag:       gethToml,
		genesisJsonFlag:          gethGenesis,
		prysmChainConfigFileFlag: configYaml,
		prysmGenesisStateFlag:    genesisState,
	}

	if len(os.Args) < 2 {
		return errNotEnoughArguments
	}

	for _, flag := range ctx.Command.Flags {
		names := flag.Names()
		if len(names) < 1 {
			return errNotEnoughArguments
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
