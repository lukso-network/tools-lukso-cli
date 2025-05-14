package configs

var Nimbus2ConfigDependencies = map[string]ClientConfigDependency{
	nimbus2MainnetChainConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/config.yaml",
		nimbus2MainnetChainConfigDependencyName,
		MainnetConfig+"/"+Nimbus2ChainConfigYamlPath,
	),
	nimbus2TestnetChainConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/config.yaml",
		nimbus2TestnetChainConfigDependencyName,
		TestnetConfig+"/"+Nimbus2ChainConfigYamlPath,
	),
	nimbus2MainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/nimbus.toml",
		nimbus2MainnetConfigDependencyName,
		MainnetConfig+"/"+Nimbus2TomlPath,
	),
	nimbus2TestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/nimbus.toml",
		nimbus2TestnetConfigDependencyName,
		TestnetConfig+"/"+Nimbus2TomlPath,
	),
	nimbus2MainnetDepositContractBlockHashDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/deposit_contract_block_hash.txt",
		nimbus2MainnetDepositContractBlockHashDependencyName,
		MainnetConfig+"/"+Nimbus2DepositContractBlockHashPath,
	),
	nimbus2MainnetBootnodesDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/bootnodes.txt",
		nimbus2MainnetBootnodesDependencyName,
		MainnetConfig+"/"+Nimbus2Bootnodes,
	),
	nimbus2TestnetDepositContractBlockHashDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/deposit_contract_block_hash.txt",
		nimbus2TestnetDepositContractBlockHashDependencyName,
		TestnetConfig+"/"+Nimbus2DepositContractBlockHashPath,
	),
	nimbus2MainnetDeployBlockDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/deploy_block.txt",
		nimbus2MainnetDeployBlockDependencyName,
		MainnetConfig+"/"+Nimbus2DeployBlockPath,
	),
	nimbus2TestnetDeployBlockDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/deploy_block.txt",
		nimbus2TestnetDeployBlockDependencyName,
		TestnetConfig+"/"+Nimbus2DeployBlockPath,
	),
	nimbus2ValidatorMainnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/validator.toml",
		nimbus2ValidatorMainnetConfigDependencyName,
		MainnetConfig+"/"+Nimbus2ValidatorTomlPath,
	),
	nimbus2ValidatorTestnetConfigDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/validator.toml",
		nimbus2ValidatorTestnetConfigDependencyName,
		TestnetConfig+"/"+Nimbus2ValidatorTomlPath,
	),
	nimbus2TestnetBootnodesDependencyName: newClientConfig(
		"https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/bootnodes.txt",
		nimbus2TestnetBootnodesDependencyName,
		TestnetConfig+"/"+Nimbus2Bootnodes,
	),
}
