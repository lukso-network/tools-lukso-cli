package configs

var TekuConfigDependencies = map[string]ClientConfigDependency{
	tekuMainnetChainConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/config.yaml",
		tekuMainnetChainConfigDependencyName,
		MainnetConfig+"/"+TekuChainConfigYamlPath,
	),
	tekuTestnetChainConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/config.yaml",
		tekuTestnetChainConfigDependencyName,
		TestnetConfig+"/"+TekuChainConfigYamlPath,
	),
	tekuMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/teku.yaml",
		tekuMainnetChainConfigDependencyName,
		MainnetConfig+"/"+TekuYamlPath,
	),
	tekuTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/teku.yaml",
		tekuTestnetChainConfigDependencyName,
		TestnetConfig+"/"+TekuYamlPath,
	),
	tekuValidatorMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/teku/validator.yaml",
		tekuValidatorMainnetConfigDependencyName,
		MainnetConfig+"/"+TekuValidatorYamlPath,
	),
	tekuValidatorTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/teku/validator.yaml",
		tekuValidatorTestnetConfigDependencyName,
		TestnetConfig+"/"+TekuValidatorYamlPath,
	),
}
