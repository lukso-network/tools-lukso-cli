package configs

var LighthouseConfigDependencies = map[string]ClientConfigDependency{
	lighthouseMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/lighthouse/lighthouse.toml",
		lighthouseMainnetConfigDependencyName,
		MainnetConfig+"/"+LighthouseTomlPath,
	),
	lighthouseTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/lighthouse/lighthouse.toml",
		lighthouseTestnetConfigDependencyName,
		TestnetConfig+"/"+LighthouseTomlPath,
	),
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
