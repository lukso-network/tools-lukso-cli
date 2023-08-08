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
		tekuMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/teku.yaml",
			name:     tekuMainnetChainConfigDependencyName,
			filePath: MainnetConfig + "/" + TekuYamlPath,
		},
		tekuTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/teku.yaml",
			name:     tekuTestnetChainConfigDependencyName,
			filePath: TestnetConfig + "/" + TekuYamlPath,
		},
		tekuValidatorMainnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/validator.yaml",
			name:     tekuValidatorMainnetConfigDependencyName,
			filePath: MainnetConfig + "/" + TekuValidatorYamlPath,
		},
		tekuValidatorTestnetConfigDependencyName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/validator.yaml",
			name:     tekuValidatorTestnetConfigDependencyName,
			filePath: TestnetConfig + "/" + TekuValidatorYamlPath,
		},
	}
)
