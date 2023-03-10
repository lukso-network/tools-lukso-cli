# LUKSO CLI
>⚠️ This page may change. Not everything is ready yet.


## Repository struct
- [`./cmd/lukso`](./cmd/lukso): code of LUKSO CLI
- [`./abis`](./abis) - collection of ABIs from smart contracts that are being interacted with
- [`./contracts`](./contracts) - collection of said smart contracts
- [`./contracts/bindings`](./contracts/bindings) - bindings generated from ABIs - to generate new bindings see [Generate bindings](#generate-bindings) section.

## Installation ( Linux/MacOS )

>🛠️ Work In Progress, available soon.

## Installation ( Windows )
>🛠️ Work In Progress, available soon.

## Running
Enter `lukso start` to start a node.

## Available parameters
`lukso <command> [geth, prysm, validator, *all*] [--flags]`
> *all* means that you can skip an argument for all possible options to run (default, only option for download)

| Command        | Description                                  |
|----------------|----------------------------------------------|
| install        | Downloads all client(s)                      |
| init           | Initializes configuration files              |
| update         | sets client(s) to desired version            |
| start          | Starts up all or specific client(s)          |
| stop           | Stops all or specific client(s)              |
| log            | Show logs                                    |
| status         | Shows status of all or specified client      |
| reset          | Resets data directories                      |
| validator      | Manages validator-related commands           |
| validator init | Initializes your validator with deposit keys | 


### start

How to use flags with values? Provide a flag and value like: `lukso start --datadir /data/network-node`

| Name                                | Description                                             | Argument                          | Default value                                          |
|-------------------------------------|---------------------------------------------------------|-----------------------------------|--------------------------------------------------------|
| --mainnet                           | Run for mainnet (default network)                       | Bool                              | false                                                  |
| --testnet                           | Run for testnet                                         | Bool                              | false                                                  |
| --devnet                            | Run for devnet                                          | Bool                              | false                                                  |
| --geth-datadir                      | A path of geth's data directory                         | Path                              | ./execution_data                                       |
| --geth-ws                           | Enable WS server                                        | None                              | true                                                   |
| --geth-ws-addr                      | Address of WS server                                    | IP Address                        | 0.0.0.0                                                |
| --geth-ws-origins                   | Origins to accept requests from                         | WS Origins OR wildcard            | *                                                      |
| --geth-ws-apis                      | Comma separated apis                                    | String of apis separated by comma | "net,eth,debug,engine"                                 |
| --geth-bootnodes                    | Bootnode addresses                                      | Bootnode addresses                | See [Bootnodes](#bootnodes)                            |
| --geth-networkid                    | Network ID                                              | Integer                           | 2022                                                   |
| --geth-nat                          | Sets HTTP nat to assign static IP for geth              | Example: "extip:0.0.0.0"          | extip:83.144.95.19                                     |
| --geth-http                         | Enable HTTP server                                      | None                              | true                                                   |
| --geth-http-apis                    | Comma separated apis                                    | String of apis separated by comma | "net,eth,debug,engine,txlookup"                        |
| --geth-http-addr                    | Address used in HTTP comunication                       | IP address                        | 0.0.0.0                                                |
| --geth-http-corsdomain              | Origins to accept requests from                         | HTTP Origins OR wildcard          | *                                                      |
| --geth-http-vhosts                  | Geth's virtual hosts                                    | Virtual hostnames OR wildcard     | *                                                      |
| --geth-ipcdisable                   | Disable IPC communication                               | None                              | True                                                   |
| --geth-ethstats                     | URL of ethstats service                                 | URL                               | ""                                                     |
| --geth-metrics                      | Enable metrics system                                   | None                              | True                                                   |
| --geth-metrics-addr                 | Address of service managing collected metrics           | IP Address                        | 0.0.0.0                                                |
| --geth-syncmode                     | Sync mode                                               | Sync mode level                   | full                                                   |
| --geth-gcmode                       | Garbage colelction mode                                 | Garbage collection level          | archive                                                |
| --geth-tx-look-up-limit             | Number of blocks to maintain tx indexes from            | Integer                           | 1                                                      |
| --geth-cache-preimages              | Enable cache preimaging                                 | None                              | True                                                   |
| --geth-verbosity                    | Verbosity for geth logging                              | Verbosity level                   | 3                                                      |
| --geth-port                         | Geth's port                                             | Port                              | 30405                                                  |
| --geth-http-port                    | Geth's HTTP port                                        | Port                              | 8565                                                   |
| --geth-mine                         | Enable mining                                           | None                              | True                                                   |
| --geth-miner-threads                | Number of CPU threads used for mining                   | Integer                           | 1                                                      |
| --geth-miner-gaslimit               | Gas ceiling                                             | Integer                           | 60000000                                               |
| --geth-miner-etherbase              | Your ECDSA public key used to get rewards on geth chain | Public address                    | 0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b             |
| --geth-auth-jwt-secret              | Path to JWT 32-byte secret                              | Path                              | ./config/mainnet/shared/secrets/jwt.hex                |
| --geth-std-output                   | Set output to console                                   | None                              | False                                                  |
| --geth-output-dir                   | Directory where logs are created                        | Path                              | ./logs/execution/geth                                  |
| --prysm-genesis-state               | Genesis state file path                                 | Path                              | ./config/mainnet/shared/genesis.ssz                    |
| --prysm-datadir                     | A path of prysm's beacon chain data directory           | Path                              | ./consensus_data                                       |
| --prysm-execution-endpoint          | Execution endpoint                                      | URL                               | http://localhost:8551                                  |
| --prysm-bootstrap-nodes             | Bootnode addresses                                      | Bootnode addresses                | See [Bootnodes](#bootnodes)                            |
| --prysm-jwt-secret                  | Path to JWT 32-byte secret                              | Path                              | ./config/mainnet/shared/secrets/jwt.hex                |
| --prysm-suggested-fee-recipient     | Address that receives block fees                        | Public address                    | 0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b             |
| --prysm-min-sync-peers              | Minimum sync peers number for prysm                     | Integer                           | 0                                                      |
| --prysm-p2p-host                    | P2P host IP                                             | IP address                        | Empty                                                  |
| --prysm-deposit-deployment          | Deployemnt height of deposit contract                   | Integer                           | 0                                                      |
| --prysm-chain-config-file           | Path to config.yaml file                                | Path                              | ./config/mainnet/shared/config.yaml                    |
| --prysm-monitoring-host             | Host used for interacting with Prometheus metrics       | IP address                        | 0.0.0.0                                                |
| --prysm-grpc-gateway-host           | Host for gRPC gateway                                   | IP address                        | 0.0.0.0                                                |
| --prysm-rpc-host                    | RPC server host                                         | IP address                        | 0.0.0.0                                                |
| --prysm-verbosity                   | Verbosity for Prysm logs                                | Log level                         | info                                                   |
| --prysm-p2p-max-peers               | Max peers for prysm                                     | Integer                           | 250                                                    |
| --prysm-subscribe-all-subnets       | Subscribe to all possible subnets                       | None                              | True                                                   |
| --prysm-minimum-peers-per-subnet    | Minimum peers per subnet                                | Integer                           | 0                                                      |
| --prysm-enable-rpc-debug-endpoints  | Enable debugging RPC endpoints                          | None                              | True                                                   |
| --prysm-output-dir                  | Directory where logs are created                        | Path                              | ./logs/consensus/beacon_chain                          |
| --prysm-std-output                  | Set output to console                                   | None                              | False                                                  |
| --validator-datadir                 | A path of validator's data directory                    | Path                              | ./validator_data                                       |
| --validator-verbosity               | Verbosity for validator logs                            | Log level                         | info                                                   |
| --validator-wallet-dir              | Location of generated wallet                            | Path                              | ./mainnet_keystore                                     |
| --validator-wallet-password-file    | Location of password used for wallet generation         | Path                              | ./config/mainnet/shared/secrets/validator-password.txt |
| --validator-chain-config-file       | Path to config.yaml file                                | Path                              | ./config/mainnet/shared/config.yaml                    |
| --validator-monitoring-host         | Host used for interacting with Prometheus metrics       | IP address                        | 0.0.0.0                                                |
| --validator-grpc-gateway-host       | Host for gRPC gateway                                   | IP address                        | 0.0.0.0                                                |
| --validator-rpc-host                | RPC server host                                         | IP address                        | 0.0.0.0                                                |
| --validator-suggested-fee-recipient | Address that receives block fees                        | Public address                    | 0x8eFdC93aE5FEa9287e7a22B6c14670BfcCdA997b             |
| --validator-output-dir              | Directory where logs are created                        | Path                              | ./logs/consensus/validator                             |
| --validator-std-output              | Set output to console                                   | None                              | False                                                  |

#### Bootnodes

### download
| Name               | Description                                           | Argument                    |
|--------------------|-------------------------------------------------------|-----------------------------|
| --accept-terms     | Accept Terms provided by clients you want to download | None                        |
| --geth-tag         | Tag of geth's version that you want to download       | Tag, ex. `1.0.0`            |
| --geth-commit-hash | Commit hash that matches provided tag commit          | Commit Hash, ex. `12345678` |
| --validator-tag    | Tag of validator's version that you want to download  | Tag, ex. `v1.0.0`           |
| --prysm-tag        | Tag of prysm's version that you want to download      | Tag, ex. `v1.0.0`           |

Note difference in tags between geth and prysm/validator (`v` at the beginning)

### update
| Name            | Description                                          | Argument          |
|-----------------|------------------------------------------------------|-------------------|
| --geth-tag      | Tag of geth's version that you want to download      | Tag, ex. `1.0.0`  |
| --validator-tag | Tag of validator's version that you want to download | Tag, ex. `v1.0.0` |
| --prysm-tag     | Tag of prysm's version that you want to download     | Tag, ex. `v1.0.0` |

### log
| Name                    | Description                                     | Argument | Default          |
|-------------------------|-------------------------------------------------|----------|------------------|
| --geth-output-file      | Path to geth log file that you want to log      | Path     | "./mainnet-logs" |
| --prysm-output-file     | Path to prysm log file that you want to log     | Path     | "./mainnet-logs" |
| --validator-output-file | Path to validator log file that you want to log | Path     | "./mainnet-logs" |
| --mainnet               | Run for mainnet (default network)               | Bool     | false            |
| --testnet               | Run for testnet                                 | Bool     | false            |
| --devnet                | Run for devnet                                  | Bool     | false            |

### reset
| Name                | Description                       | Argument | Default                    |
|---------------------|-----------------------------------|----------|----------------------------|
| --geth-datadir      | geth datadir                      | Path     | "./mainnet-data/execution" |
| --prysm-datadir     | prysm datadir                     | Path     | "./mainnet-data/consensus" |
| --validator-datadir | validator datadir                 | Path     | "./mainnet-data/validator" |
| --mainnet           | Run for mainnet (default network) | Bool     | false                      |
| --testnet           | Run for testnet                   | Bool     | false                      |
| --devnet            | Run for devnet                    | Bool     | false                      |



### validator
| Name                | Description                                                                       | Argument | Default                              |
|---------------------|-----------------------------------------------------------------------------------|----------|--------------------------------------|
| --deposit           | Path to your deposit file - makes a deposit to a deposit contract                 | Path     | ""                                   |
| --genesis-deposit   | Path to your genesis deposit file - makes a deposit to genesis validator contract | Path     | ""                                   |
| --rpc               | Your RPC provider                                                                 | URL      | "https://rpc.2022.l16.lukso.network" |
| --gas-price         | Gas price provided by user                                                        | Int      | 1000000000                           |
| --max-txs-per-block | Maximum amount of txs sent per single block                                       | Int      | 10                                   |

### validator init
| Name                                   | Description                           | Argument | Default              |
|----------------------------------------|---------------------------------------|----------|----------------------|
| --validator-wallet-dir value           | location of generated wallet          | Path     | "./mainnet-keystore" |
| --validator-keys-dir value             | Path to your validator keys           | Path     |                      |
| --validator-wallet-password-file value | Path to your password file            | Path     |                      |
| --mainnet                              | Run for mainnetFlag (default network) | Bool     | false                |
| --testnet                              | Run for testnetFlag                   | Bool     | false                |
| --devnet                               | Run for devnet                        | Bool     | false                |

## Generate bindings
### Prerequisites:
- solc (https://github.com/ethereum/solidity)
- abigen (https://geth.ethereum.org/docs/tools/abigen)

### Steps

1) Paste your smart contract that you want to interact with into [`./contracts`](./contracts) directory
2) Generate ABI from your smart contract:
```bash
$ solcjs --output-dir abis --abi contracts/depositContract.sol
```
3) Generate bindings using newly generated ABI
```bash
abigen --abi abis/your-abi-file --pkg bindings --out contracts/bindings/yourBindingFile.go --type TypeName
```
4) To use binding in code type in:
```go
bind, err := bindings.NewTypeName(common.HexToAddress(contractAddress), ethClient)
if err != nil {
	return
}

tx, err := bind.DoSomething(...)
```
