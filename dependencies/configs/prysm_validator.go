package configs

var PrysmValidatorConfigDependencies = map[string]ClientConfigDependency{
	validatorMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/validator.yaml",
		validatorMainnetConfigDependencyName,
		MainnetConfig+"/"+ValidatorYamlPath,
	),
	validatorTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/validator.yaml",
		validatorTestnetConfigDependencyName,
		TestnetConfig+"/"+ValidatorYamlPath,
	),
}
