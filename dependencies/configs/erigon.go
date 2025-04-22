package configs

var ErigonConfigDependencies = map[string]ClientConfigDependency{
	erigonMainnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/erigon/erigon.toml",
		erigonMainnetConfigName,
		MainnetConfig+"/"+ErigonTomlPath,
	),
	erigonTestnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/erigon/erigon.toml",
		erigonTestnetConfigName,
		TestnetConfig+"/"+ErigonTomlPath,
	),
}
