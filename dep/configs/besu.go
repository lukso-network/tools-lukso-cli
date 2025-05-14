package configs

var BesuConfigDependencies = map[string]ClientConfigDependency{
	besuMainnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/besu/besu.toml",
		besuMainnetConfigName,
		MainnetConfig+"/"+BesuTomlPath,
	),
	besuTestnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/besu/besu.toml",
		besuTestnetConfigName,
		TestnetConfig+"/"+BesuTomlPath,
	),
}
