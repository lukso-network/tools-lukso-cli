package configs

var PrysmValidatorConfigDependencies = map[string]ClientConfigDependency{
	prysmValidatorMainnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/validator.yaml",
		name:     prysmValidatorMainnetConfigDependencyName,
		filePath: MainnetConfig + "/" + PrysmValidatorYamlPath,
	},
	prysmValidatorTestnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/validator.yaml",
		name:     prysmValidatorTestnetConfigDependencyName,
		filePath: TestnetConfig + "/" + PrysmValidatorYamlPath,
	},
}
