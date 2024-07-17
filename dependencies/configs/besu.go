package configs

var BesuConfigDependencies = map[string]ClientConfigDependency{
	besuMainnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/besu/besu.toml",
		name:     besuMainnetConfigName,
		filePath: MainnetConfig + "/" + BesuTomlPath,
	},
	besuTestnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/besu/besu.toml",
		name:     besuTestnetConfigName,
		filePath: TestnetConfig + "/" + BesuTomlPath,
	},
}
