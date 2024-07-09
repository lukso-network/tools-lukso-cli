package configs

var BesuConfigDependencies = map[string]ClientConfigDependency{
	besuMainnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/besu/besu.yaml",
		name:     besuMainnetConfigName,
		filePath: MainnetConfig + BesuYamlPath,
	},
	besuTestnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/besu/besu.yaml",
		name:     besuTestnetConfigName,
		filePath: TestnetConfig + BesuYamlPath,
	},
}
