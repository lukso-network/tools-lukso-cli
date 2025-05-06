package configs

var PrysmConfigDependencies = map[string]ClientConfigDependency{
	prysmMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/prysm.yaml",
		prysmMainnetConfigDependencyName,
		MainnetConfig+"/"+PrysmYamlPath,
	),
	prysmTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/prysm.yaml",
		prysmTestnetConfigDependencyName,
		TestnetConfig+"/"+PrysmYamlPath,
	),
}
