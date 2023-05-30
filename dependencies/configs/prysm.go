package configs

var (
	PrysmConfigDependencies = map[string]ClientConfigDependency{
		prysmMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/prysm.yaml",
			name:     prysmMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + prysmYamlPath,
		},
		prysmTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/prysm.yaml",
			name:     prysmTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + prysmYamlPath,
		},
	}
)
