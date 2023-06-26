package configs

var (
	GethConfigDependencies = map[string]ClientConfigDependency{
		gethMainnetConfigName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
			name:     gethMainnetConfigName,
			filePath: MainnetConfig + "/" + GethTomlPath,
		},
		gethTestnetConfigName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
			name:     gethTestnetConfigName,
			filePath: TestnetConfig + "/" + GethTomlPath,
		},
	}
)
