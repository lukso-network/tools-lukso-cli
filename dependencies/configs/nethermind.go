package configs

var NethermindConfigDependencies = map[string]ClientConfigDependency{
	nethermindMainnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nethermind/nethermind.cfg",
		nethermindMainnetConfigName,
		MainnetConfig+"/"+NethermindJsonPath,
	),
	nethermindTestnetConfigName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nethermind/nethermind.cfg",
		nethermindTestnetConfigName,
		TestnetConfig+"/"+NethermindJsonPath,
	),
}
