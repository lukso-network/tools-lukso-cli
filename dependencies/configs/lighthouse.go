package configs

var (
	LighthouseConfigDependencies = map[string]ClientConfigDependency{
		lighthouseMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/lighthouse.toml",
			name:     lighthouseMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + lighthouseTomlPath,
		},
		lighthouseTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/lighthouse.toml",
			name:     lighthouseTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + lighthouseTomlPath,
		},
	}
)
