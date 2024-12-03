package flags

const (
	// Network
	GenesisStateFlag = "genesis-ssz"
	GenesisJsonFlag  = "genesis-json"

	// Geth
	GethTagFlag        = "geth-tag"
	GethCommitHashFlag = "geth-commit-hash"
	GethConfigFileFlag = "geth-config"

	// Erigon
	ErigonTagFlag        = "erigon-tag"
	ErigonConfigFileFlag = "erigon-config"

	// Besu
	BesuTagFlag        = "besu-tag"
	BesuConfigFileFlag = "besu-config"

	// Nethermind
	NethermindTagFlag        = "nethermind-tag"
	NethermindCommitHashFlag = "nethermind-commit-hash"
	NethermindConfigFileFlag = "nethermind-config"

	// Prysm
	PrysmTagFlag             = "prysm-tag"
	PrysmChainConfigFileFlag = "prysm-chain-config"
	PrysmConfigFileFlag      = "prysm-config"
	NoSlasherFlag            = "no-slasher"

	// Lighthouse
	LighthouseTagFlag                 = "lighthouse-tag"
	LighthouseConfigFileFlag          = "lighthouse-config"
	LighthouseValidatorConfigFileFlag = "lighthouse-validator-config"
	TestnetDirFlag                    = "testnet-dir"

	// Teku
	TekuTagFlag                 = "teku-tag"
	TekuConfigFileFlag          = "teku-config"
	TekuValidatorConfigFileFlag = "teku-validator-config"

	// Nimbus-eth2
	Nimbus2TagFlag                 = "nimbus2-tag"
	Nimbus2CommitHashFlag          = "nimbus2-commit-hash"
	Nimbus2NetworkFlag             = "nimbus2-network"
	Nimbus2ConfigFileFlag          = "nimbus2-config"
	Nimbus2ValidatorConfigFileFlag = "nimbus2-validator-config"

	// Validator
	ValidatorTagFlag                = "validator-tag"
	ValidatorWalletPasswordFileFlag = "validator-wallet-password"
	ValidatorWalletDirFlag          = "validator-wallet-dir"
	ValidatorConfigFileFlag         = "validator-config"
	ValidatorChainConfigFileFlag    = "validator-chain-config"

	// Node folder
	MainnetFlag   = "mainnet"
	TestnetFlag   = "testnet"
	ValidatorFlag = "validator"
	ConsensusFlag = "consensus"
	ExecutionFlag = "execution"

	LogFolderFlag        = "logs-folder"
	DatadirFlag          = "datadir"
	ExecutionDatadirFlag = "execution-datadir"
	ConsensusDatadirFlag = "consensus-datadir"
	ValidatorDatadirFlag = "validator-datadir"

	ValidatorKeysFlag           = "validator-keys"
	ValidatorPasswordFlag       = "validator-password"
	KeystoreFlag                = "keystore"
	RpcAddressFlag              = "rpc-address"
	CheckpointSyncFlag          = "checkpoint-sync"
	TransactionFeeRecipientFlag = "transaction-fee-recipient"
	AgreeTermsFlag              = "agree-terms"
	AllFlag                     = "all"

	// Misc.
	// Status host:port
	ExecutionClientHost = "execution-client-host"
	ConsensusClientHost = "consensus-client-host"
	ValidatorClientHost = "validator-client-host"
	ExecutionClientPort = "execution-client-port"
	ConsensusClientPort = "consensus-client-port"
	ValidatorClientPort = "validator-client-port"
)
