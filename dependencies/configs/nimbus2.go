package configs

var Nimbus2ConfigDependencies = map[string]ClientConfigDependency{
	nimbus2MainnetChainConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/config.yaml",
		name:     nimbus2MainnetChainConfigDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2ChainConfigYamlPath,
	},
	nimbus2TestnetChainConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/config.yaml",
		name:     nimbus2TestnetChainConfigDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2ChainConfigYamlPath,
	},
	nimbus2MainnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/nimbus.toml",
		name:     nimbus2MainnetConfigDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2TomlPath,
	},
	nimbus2TestnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/nimbus.toml",
		name:     nimbus2TestnetConfigDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2TomlPath,
	},
	nimbus2MainnetDepositContractBlockHashDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/deposit_contract_block_hash.txt",
		name:     nimbus2MainnetDepositContractBlockHashDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2DepositContractBlockHashPath,
	},
	nimbus2TestnetDepositContractBlockHashDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/deposit_contract_block_hash.txt",
		name:     nimbus2TestnetDepositContractBlockHashDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2DepositContractBlockHashPath,
	},
	nimbus2MainnetDeployBlockDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/deploy_block.txt",
		name:     nimbus2MainnetDeployBlockDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2DeployBlockPath,
	},
	nimbus2TestnetDeployBlockDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/deploy_block.txt",
		name:     nimbus2TestnetDeployBlockDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2DeployBlockPath,
	},
	nimbus2ValidatorMainnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/validator.toml",
		name:     nimbus2ValidatorMainnetConfigDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2ValidatorTomlPath,
	},
	nimbus2ValidatorTestnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/validator.toml",
		name:     nimbus2ValidatorTestnetConfigDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2ValidatorTomlPath,
	},
}
