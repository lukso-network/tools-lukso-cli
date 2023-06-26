package configs

var (
	LighthouseValidatorConfigDependencies = map[string]ClientConfigDependency{
		lighthouseValidatorMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/validator.toml",
			name:     lighthouseValidatorMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + LighthouseValidatorTomlPath,
		},
		lighthouseValidatorTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/validator.toml",
			name:     lighthouseValidatorTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + LighthouseValidatorTomlPath,
		},
	}
)
