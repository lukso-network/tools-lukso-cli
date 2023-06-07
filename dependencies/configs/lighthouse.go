package configs

var (
	LighthouseConfigDependencies = map[string]ClientConfigDependency{
		lighthouseMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/lighthouse.toml",
			name:     lighthouseMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + LighthouseTomlPath,
		},
		lighthouseTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/lighthouse.toml",
			name:     lighthouseTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + LighthouseTomlPath,
		},
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
