package configs

var (
	PrysmValidatorConfigDependencies = map[string]ClientConfigDependency{
		validatorMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/validator.yaml",
			name:     validatorMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + ValidatorYamlPath,
		},
		validatorTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/validator.yaml",
			name:     validatorTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + ValidatorYamlPath,
		},
	}
)
