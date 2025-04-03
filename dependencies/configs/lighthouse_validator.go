package configs

var LighthouseValidatorConfigDependencies = map[string]ClientConfigDependency{
	lighthouseValidatorMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/validator.toml",
		lighthouseValidatorMainnetConfigDependencyName,
		MainnetConfig+"/"+LighthouseValidatorTomlPath,
	),
	lighthouseValidatorTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/validator.toml",
		lighthouseValidatorTestnetConfigDependencyName,
		TestnetConfig+"/"+LighthouseValidatorTomlPath,
	),
}
