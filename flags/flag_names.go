package flags

const (
	GethTagFlag        = "geth-tag"
	GethCommitHashFlag = "geth-commit-hash"
	GethDatadirFlag    = "geth-datadir"
	GethConfigFileFlag = "geth-config"
	GenesisJsonFlag    = "genesis-json"

	ErigonTagFlag        = "erigon-tag"
	ErigonConfigFileFlag = "erigon-config"
	ErigonDatadirFlag    = "erigon-datadir"

	NethermindTagFlag        = "nethermind-tag"
	NethermindCommitHashFlag = "nethermind-commit-hash"
	NethermindConfigFileFlag = "nethermind-config"
	NethermindDatadirFlag    = "nethermind-datadir"

	BesuTagFlag        = "besu-tag"
	BesuConfigFileFlag = "besu-config"
	BesuDatadirFlag    = "besu-datadir"

	PrysmTagFlag             = "prysm-tag"
	GenesisStateFlag         = "genesis-ssz"
	PrysmChainConfigFileFlag = "prysm-chain-config"
	PrysmConfigFileFlag      = "prysm-config"
	PrysmDatadirFlag         = "prysm-datadir"
	NoSlasherFlag            = "no-slasher"

	LighthouseTagFlag                 = "lighthouse-tag"
	LighthouseConfigFileFlag          = "lighthouse-config"
	LighthouseValidatorConfigFileFlag = "lighthouse-validator-config"
	LighthouseDatadirFlag             = "lighthouse-datadir"
	TestnetDirFlag                    = "testnet-dir"

	ValidatorTagFlag                = "validator-tag"
	ValidatorDatadirFlag            = "validator-datadir"
	ValidatorWalletPasswordFileFlag = "validator-wallet-password"
	ValidatorWalletDirFlag          = "validator-wallet-dir"
	ValidatorConfigFileFlag         = "validator-config"
	ValidatorChainConfigFileFlag    = "validator-chain-config"

	TekuTagFlag                 = "teku-tag"
	TekuConfigFileFlag          = "teku-config"
	TekuValidatorConfigFileFlag = "teku-validator-config"
	MainnetFlag                 = "mainnet"
	TestnetFlag                 = "testnet"
	DevnetFlag                  = "devnet"
	ValidatorFlag               = "validator"
	ConsensusFlag               = "consensus"
	ExecutionFlag               = "execution"

	LogFolderFlag = "logs-folder"

	ValidatorKeysFlag           = "validator-keys"
	ValidatorPasswordFlag       = "validator-password"
	KeystoreFlag                = "keystore"
	RpcAddressFlag              = "rpc-address"
	CheckpointSyncFlag          = "checkpoint-sync"
	TransactionFeeRecipientFlag = "transaction-fee-recipient"
	AgreeTermsFlag              = "agree-terms"
	AllFlag                     = "all"

	ExecutionClientHost = "execution-client-host"
	ConsensusClientHost = "consensus-client-host"
	ValidatorClientHost = "validator-client-host"
	ExecutionClientPort = "execution-client-port"
	ConsensusClientPort = "consensus-client-port"
	ValidatorClientPort = "validator-client-port"
)
