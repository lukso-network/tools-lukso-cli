package configs

var (
	TekuConfigDependencies = map[string]ClientConfigDependency{
		tekuMainnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/config.yaml",
			name:     tekuMainnetChainConfigDependencyName,
			filePath: MainnetConfig + "/" + TekuChainConfigYamlPath,
		},
		tekuTestnetChainConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/config.yaml",
			name:     tekuTestnetChainConfigDependencyName,
			filePath: TestnetConfig + "/" + TekuChainConfigYamlPath,
		},
	}
)
