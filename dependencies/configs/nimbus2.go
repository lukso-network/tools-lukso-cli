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
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/nimbus2.yaml",
		name:     nimbus2MainnetChainConfigDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2TomlPath,
	},
	nimbus2TestnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/nimbus2.yaml",
		name:     nimbus2TestnetChainConfigDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2TomlPath,
	},
	nimbus2MainnetDepositBlockHashDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/nimbus2/nimbus2.yaml",
		name:     nimbus2MainnetDepositBlockHashDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2DepositContractHashPath,
	},
	nimbus2TestnetDepositBlockHashDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/nimbus2/nimbus2.yaml",
		name:     nimbus2TestnetDepositBlockHashDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2DepositContractHashPath,
	},
	nimbus2ValidatorMainnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/mainnet/prysm/validator.yaml",
		name:     nimbus2ValidatorMainnetConfigDependencyName,
		filePath: MainnetConfig + "/" + Nimbus2ValidatorTomlPath,
	},
	nimbus2ValidatorTestnetConfigDependencyName: &clientConfig{
		url:      "https://raw.githubusercontent.com/lukso-network/network-configs/main/testnet/prysm/validator.yaml",
		name:     nimbus2ValidatorTestnetConfigDependencyName,
		filePath: TestnetConfig + "/" + Nimbus2ValidatorTomlPath,
	},
}
