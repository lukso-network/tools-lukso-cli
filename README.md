# LUKSO CLI
>⚠️ This page may change. Not everything is ready yet.


## Repository struct
In `./cmd/lukso` lies code of Lukso CLI

## Installation ( Linux/MacOS )

>🛠️ Work In Progress, available soon.

## Installation ( Windows )
>🛠️ Work In Progress, available soon.

## Running
Enter `lukso start` to start a node.

## Available parameters
`lukso <command> [geth, prysm, validator, *all*] [--flags]`
> *all* means that you can skip an argument for all possible options to run (default, only option for download)

| Command  | Description                             |
|----------|-----------------------------------------|
| download | Downloads all client(s)                 |
| init     | Initializes configuration files         |
| update   | sets client(s) to desired version       |
| start    | Starts up all or specific client(s)     |
| stop     | Stops all or specific client(s)         |
| logs     | Show logs                               |
| status   | Shows status of all or specified client |


### start

How to use flags with values? Provide a flag and value like: `lukso start --datadir /data/network-node`

| Name                                | Description                                             | Argument                          | Default value                                          |
|-------------------------------------|---------------------------------------------------------|-----------------------------------|--------------------------------------------------------|
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

### logs
| Name                    | Description                                     | Argument                                   |
|-------------------------|-------------------------------------------------|--------------------------------------------|
| --geth-output-file      | Path to geth log file that you want to log      | Path, ex. `./logs/log_folder/log_file.log` |
| --prysm-output-file     | Path to prysm log file that you want to log     | Path, ex. `./logs/log_folder/log_file.log` |
| --validator-output-file | Path to validator log file that you want to log | Path, ex. `./logs/log_folder/log_file.log` |

NOTE: `logs` command is broken after changing structure of log files (adding timestamp to filename).
