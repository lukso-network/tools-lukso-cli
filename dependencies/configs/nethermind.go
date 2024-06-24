package configs

var NethermindConfigDependencies = map[string]ClientConfigDependency{
	nethermindMainnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nethermind/nethermind.cfg",
		name:     nethermindMainnetConfigName,
		filePath: MainnetConfig + "/" + NethermindCfgPath,
	},
	nethermindTestnetConfigName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nethermind/nethermind.cfg",
		name:     nethermindTestnetConfigName,
		filePath: TestnetConfig + "/" + NethermindCfgPath,
	},
}
