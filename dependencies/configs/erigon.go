package configs

var (
	ErigonConfigDependencies = map[string]ClientConfigDependency{
		erigonMainnetConfigName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/erigon/erigon.toml",
			name:     erigonMainnetConfigName,
			filePath: MainnetConfig + "/" + ErigonTomlPath,
		},
		erigonTestnetConfigName: &clientConfig{
			url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/erigon/erigon.toml",
			name:     erigonTestnetConfigName,
			filePath: TestnetConfig + "/" + ErigonTomlPath,
		},
	}
)
