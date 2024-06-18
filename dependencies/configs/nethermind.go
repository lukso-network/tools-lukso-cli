package configs

var NethermindConfigDependencies = map[string]ClientConfigDependency{
	nethermindMainnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/geth/geth.toml",
		name:     nethermindMainnetConfigName,
		filePath: MainnetConfig + "/" + NethermindCfgPath,
	},
	nethermindTestnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/geth/geth.toml",
		name:     nethermindTestnetConfigName,
		filePath: TestnetConfig + "/" + NethermindCfgPath,
	},
}
