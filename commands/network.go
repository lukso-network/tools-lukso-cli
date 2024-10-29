package commands

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lukso-network/tools-lukso-cli/common/errors"
	"github.com/lukso-network/tools-lukso-cli/common/utils"
	"github.com/lukso-network/tools-lukso-cli/dependencies/configs"
	"github.com/lukso-network/tools-lukso-cli/flags"
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
	testnetDirPath       string // for lighthouse
}

// SelectNetworkFor accepts a CLI func as an argument, and adjusts all values that need to be changed depending on
// network passed as a flag. Works as a wrapper for selecting current working network
func SelectNetworkFor(f func(*cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.Command.SkipFlagParsing {
			log.Debug("Skipping flag parsing on - parsing flags manually...")
			err := parseFlags(ctx)
			if err != nil {
				return utils.Exit(fmt.Sprintf("❌  There was an error while parsing flags: %v", err), 1)
			}
		}

		testnetEnabled := ctx.Bool(flags.TestnetFlag)
		devnetEnabled := ctx.Bool(flags.DevnetFlag)

		enabledCount := utils.BoolToInt(testnetEnabled) + utils.BoolToInt(devnetEnabled)
		if enabledCount > 1 {
			return errors.ErrMoreNetworksSelected
		}

		var cfg networkConfig

		cfg = networkConfig{
			executionDatadirPath: configs.ExecutionMainnetDatadir,
			consensusDatadirPath: configs.ConsensusMainnetDatadir,
			validatorDatadirPath: configs.ValidatorMainnetDatadir,
			logPath:              configs.MainnetLogs,
			configPath:           configs.MainnetConfig,
			keysPath:             configs.MainnetKeystore,
			walletPath:           configs.MainnetKeystore,
			testnetDirPath:       configs.MainnetConfig + "/shared",
		}

		if devnetEnabled {
			return utils.Exit("❌  Network not supported", 1)
		}

		if testnetEnabled {
			cfg = networkConfig{
				executionDatadirPath: configs.ExecutionTestnetDatadir,
				consensusDatadirPath: configs.ConsensusTestnetDatadir,
				validatorDatadirPath: configs.ValidatorTestnetDatadir,
				logPath:              configs.TestnetLogs,
				configPath:           configs.TestnetConfig,
				keysPath:             configs.TestnetKeystore,
				walletPath:           configs.TestnetKeystore,
				testnetDirPath:       configs.TestnetConfig + "/shared",
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
		// genesisJson  = config.configPath + "/" + genesisJsonPath
		gethTomlPath                = config.configPath + "/" + configs.GethTomlPath
		erigonTomlPath              = config.configPath + "/" + configs.ErigonTomlPath
		nethermindCfgPath           = config.configPath + "/" + configs.NethermindCfgPath
		besuTomlPath                = config.configPath + "/" + configs.BesuTomlPath
		prysmYamlPath               = config.configPath + "/" + configs.PrysmYamlPath
		lighthouseTomlPath          = config.configPath + "/" + configs.LighthouseTomlPath
		lighthouseValidatorTomlPath = config.configPath + "/" + configs.LighthouseValidatorTomlPath
		validatorYamlPath           = config.configPath + "/" + configs.ValidatorYamlPath
		tekuYamlPath                = config.configPath + "/" + configs.TekuYamlPath
		tekuValidatorYamlPath       = config.configPath + "/" + configs.TekuValidatorYamlPath
		nimbus2TomlPath             = config.configPath + "/" + configs.Nimbus2TomlPath
		nimbus2ValidatorTomlPath    = config.configPath + "/" + configs.Nimbus2ValidatorTomlPath
		gethGenesisPath             = config.configPath + "/" + configs.GenesisJsonPath
		genesisStatePath            = config.configPath + "/" + configs.GenesisStateFilePath
		configYamlPath              = config.configPath + "/" + configs.ChainConfigYamlPath
	)

	passedArgs := make([]string, 0)

	for _, arg := range os.Args {
		if strings.Contains(arg, "--") {
			passedArgs = append(passedArgs, arg)
		}
	}

	// varyingFlags represents list of all flags that can be affected by selecting network and values that may be replaced
	varyingFlags := map[string]string{
		flags.GethDatadirFlag:                   config.executionDatadirPath,
		flags.ErigonDatadirFlag:                 config.executionDatadirPath,
		flags.PrysmDatadirFlag:                  config.consensusDatadirPath,
		flags.ValidatorDatadirFlag:              config.validatorDatadirPath,
		flags.LighthouseDatadirFlag:             config.consensusDatadirPath,
		flags.LogFolderFlag:                     config.logPath,
		flags.ValidatorKeysFlag:                 config.keysPath,
		flags.GethConfigFileFlag:                gethTomlPath,
		flags.ErigonConfigFileFlag:              erigonTomlPath,
		flags.NethermindConfigFileFlag:          nethermindCfgPath,
		flags.BesuConfigFileFlag:                besuTomlPath,
		flags.PrysmConfigFileFlag:               prysmYamlPath,
		flags.LighthouseConfigFileFlag:          lighthouseTomlPath,
		flags.LighthouseValidatorConfigFileFlag: lighthouseValidatorTomlPath,
		flags.ValidatorConfigFileFlag:           validatorYamlPath,
		flags.TekuConfigFileFlag:                tekuYamlPath,
		flags.TekuValidatorConfigFileFlag:       tekuValidatorYamlPath,
		flags.Nimbus2ConfigFileFlag:             nimbus2TomlPath,
		flags.Nimbus2ValidatorConfigFileFlag:    nimbus2ValidatorTomlPath,
		flags.GenesisJsonFlag:                   gethGenesisPath,
		flags.PrysmChainConfigFileFlag:          configYamlPath,
		flags.ValidatorChainConfigFileFlag:      configYamlPath,
		flags.GenesisStateFlag:                  genesisStatePath,
		flags.ValidatorWalletDirFlag:            config.walletPath,
		flags.TestnetDirFlag:                    config.testnetDirPath,
	}

	if len(os.Args) < 2 {
		return utils.Exit(errors.ErrNotEnoughArguments.Error(), 1)
	}

	for _, flag := range ctx.Command.Flags {
		names := flag.Names()
		if len(names) < 1 {
			return utils.Exit(errors.ErrNotEnoughArguments.Error(), 1)
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

// parseFlags takes care of parsing flags that are skipped if SkipFlagParsing is set to true - if --help or -h is found we display help and stop execution
func parseFlags(ctx *cli.Context) (err error) {
	args := ctx.Args()
	argsLen := args.Len()
	for i := 0; i < argsLen; i++ {
		arg := args.Get(i)

		if arg == "--help" || arg == "-h" {
			cli.ShowSubcommandHelpAndExit(ctx, 0)
		}

		if strings.HasPrefix(arg, "--") {
			if i+1 == argsLen {
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, "true")
				if err != nil && strings.Contains(err.Error(), errors.NoSuchFlag) {
					err = nil

					return
				}

				return
			}

			// we found a flag for our client - now we need to check if it's a value or bool flag
			nextArg := args.Get(i + 1)
			if strings.HasPrefix(nextArg, "--") { // we found a next flag, so current one is a bool
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, "true")
				if err == nil || (err != nil && strings.Contains(err.Error(), errors.NoSuchFlag)) {
					err = nil

					continue
				}

				return
			} else {
				arg = strings.TrimPrefix(arg, "--")

				err = ctx.Set(arg, nextArg)
				if err == nil || (err != nil && strings.Contains(err.Error(), errors.NoSuchFlag)) {
					err = nil

					continue
				}

				return
			}
		}
	}

	return
}
