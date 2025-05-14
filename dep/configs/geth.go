package configs

var GethConfigDependencies = map[string]ClientConfigDependency{
	gethMainnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
		gethMainnetConfigName,
		MainnetConfig+"/"+GethTomlPath,
	),
	gethTestnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
		gethTestnetConfigName,
		TestnetConfig+"/"+GethTomlPath,
	),
}
