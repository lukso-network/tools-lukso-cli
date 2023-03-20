# LUKSO CLI

> ⚠️ DO NOT USE YET, this is WIP!

The `lukso` CLI serves the following purposes:

- easy installation of all node types (full installs into `/bin/` , not docker containers)
- easy starts and stops local nodes (as it runs as a daemon)
- easy access to nodes logs

## Repository struct

- [`./cmd/lukso`](./cmd/lukso): code of LUKSO CLI
- [`./abis`](./abis) - collection of ABIs from smart contracts that are being interacted with
- [`./contracts`](./contracts) - collection of said smart contracts
- [`./contracts/bindings`](./contracts/bindings) - bindings generated from ABIs - to generate new bindings see [Generate bindings](#generate-bindings) section.
- [`./cf-wrangler`](./cf-wrangler/) A small cloudflare proxy to handle redirect to install.sh from install.lukso.network
- [`./docs`](./docs) Some small content to because part of the mac pgk file

## Installation ( Linux/MacOS )

```sh
curl https://install.lukso.network | sh
```

## Running

Enter `lukso start` to start a node.

## Available parameters

`lukso <command> [geth, prysm, validator, *all*] [--flags]`

> _all_ means that you can skip an argument for all possible options to run (default, only option for download)

| Command          | Description                                  |
| ---------------- | -------------------------------------------- |
| `install`        | Downloads all client(s)                      |
| `init`           | Initializes configuration files              |
| `update`         | sets client(s) to desired version            |
| `start`          | Starts up all or specific client(s)          |
| `stop`           | Stops all or specific client(s)              |
| `log`            | Show logs                                    |
| `status`         | Shows status of all or specified client      |
| `reset`          | Resets data directories                      |
| `validator`      | Manages validator-related commands           |
| `validator init` | Initializes your validator with deposit keys |
| `version`        | Display version of LUKSO CLI                 |

##### geth | erigon | prysm | lighthouse - command lines

These flags are documented within the projects they are for; the cli passed the flags along to the underlying client software.

Geth: https://geth.ethereum.org/docs/fundamentals/command-line-options
Erigon: https://github.com/ledgerwatch/erigon
Prysm: https://docs.prylabs.network/docs/prysm-usage/parameters
Lighthouse: https://lighthouse-book.sigmaprime.io/advanced-datadir.html

#### install

Installs LUKSO client: `lukso install --datadir /data/network-node`

| Name                 | Description                                                                                                                                            |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `install`            | Installs LUKSO Mainnet and allows user to select which consensus and execution client they wish to install and prompts agreement to Terms & Conditions |
| `install prysm`      | Installs only the Prysm client                                                                                                                         |
| `install lighthouse` | Installs only the Lighthouse client                                                                                                                    |
| `install geth`       | Installs only the Geth client                                                                                                                          |
| `install erigon`     | Installs only the Erigon client                                                                                                                        |
| `install testnet`    | Installs LUKSO Testnet                                                                                                                                 |

#### init

`lukso init` - downloads the network configs from [https://github.com/lukso-newtork/network-configs]() to `./myLUKSOMainnetFolder/configs` . It won't overwrite any existing config folders, data or keystore the user might have

#### update

`lukso update` - updates LUKSO client to the latest version available
`lukso update prysm` - updates only the Prysm client to the latest version available
`lukso update prysm@3.2.2` - updates only the Prysm client to a specific release version
`lukso update prysm --commit "123456"` - updates only the Prysm client to a specific version
`lukso update prysm --tag "v1.0.0"` - updates only the Prysm client to a specific version

#### start

How to use flags with values? Provide a flag and value like: `lukso start --datadir /data/network-node`

`lukso start --genesis-ssz "./config/mainnet/shared/genesis.szz" --genesis-json "./config/mainnet/geth/genesis.json`

| Name      | Description                       |
| --------- | --------------------------------- |
| --mainnet | Run for mainnet (default network) |
| --testnet | Run for testnet                   |
| --devnet  | Run for devnet                    |

##### start + custom client config files

`lukso start` - Starts LUKSO client and takes the default config files from the default path located at `./config/mainnet/geth/config.toml`

Lukso can also start custom clients; with the user's custom config files for example:

`lukso start --geth-config "./myconfig.toml"`
`lukso start --geth-config "./pathtofile.toml`

`lukso start --prysm-config "./myconfig.yaml" --geth-bootnodes "mycusomtbootnode00000000"`

`lukso start --lighthouse --erigon`

##### start + validator functions

`lukso start --validator --transaction-fee-recipient "0x12345678.."` - uses the default keystore folder and asks for the validator key password; also requires the user to provide a coinbase address to receive the transaction fees.

it also checks if the

- deposit `data-xxxx.json` and key files `(keystore-m*.....json)` are present and tries to import them

- if there is NO `deposit_data-xxxx.json`, AND no accounts in the validator it will throw an error:

"No available validator keys! Please place a `deposit_data-xxxxx.json` inside the `"./mainnet-keystore"` folder

- if there is NO deposit_data-xxxx.json, AND no accounts, just start the validator

`lukso start --validator --validator-keys "./mainnet-keystore/" --validator-password "./myfile.txt"` - uses the mainnet keystore folder and asks for the user to point to the validator key password

`lukso start --validator --validator-keys "./mainnet-keystore/"` - users the mainnet keystore folder and prompts the user to type in the password.

| Name                                  | Description                                       | Argument       | Default value                                          |
| ------------------------------------- | ------------------------------------------------- | -------------- | ------------------------------------------------------ |
| `--validator-datadir`                 | A path of validator's data directory              | Path           | ./validator_data                                       |
| `--validator-verbosity`               | Verbosity for validator logs                      | Log level      | info                                                   |
| `--validator-wallet-dir`              | Location of generated wallet                      | Path           | ./mainnet_keystore                                     |
| `--validator-wallet-password-file`    | Location of password used for wallet generation   | Path           | ./config/mainnet/shared/secrets/validator-password.txt |
| `--validator-chain-config-file`       | Path to config.yaml file                          | Path           | ./config/mainnet/shared/config.yaml                    |
| `--validator-monitoring-host`         | Host used for interacting with Prometheus metrics | IP address     | 0.0.0.0                                                |
| `--validator-grpc-gateway-host`       | Host for gRPC gateway                             | IP address     | 0.0.0.0                                                |
| `--validator-rpc-host`                | RPC server host                                   | IP address     | 0.0.0.0                                                |
| `--validator-suggested-fee-recipient` | Address that receives block fees                  | Public address | 0x12345678xxxxabcde                                    |
| `--validator-output-dir `             | Directory where logs are created                  | Path           | ./logs/consensus/validator                             |
| `--validator-std-output`              | Set output to console                             | None           | False                                                  |

`lukso start --log-folder "./myCustomLogFolder"` - starts LUKSO client from a log folder and displays a log file: `geth_2023-02-10_12-21-52.log`

#### stop

`lukso stop` - stops LUKSO client
`lukso stop --validator` - stops only the validator
`lukso stop --execution` - stops only the execution client
`lukso stop --consensus` - stops only the consensus client
|

#### log

| Name                      | Description                                     | Argument | Default          |
| ------------------------- | ----------------------------------------------- | -------- | ---------------- |
| `--geth-output-file`      | Path to geth log file that you want to log      | Path     | "./mainnet-logs" |
| `--prysm-output-file`     | Path to prysm log file that you want to log     | Path     | "./mainnet-logs" |
| `--validator-output-file` | Path to validator log file that you want to log | Path     | "./mainnet-logs" |
| `--mainnet`               | Run for mainnet (default network)               | Bool     | false            |
| `--testnet`               | Run for testnet                                 | Bool     | false            |
| `--devnet`                | Run for devnet                                  | Bool     | false            |

#### status

`lukso status` - displays the status of each execution, consensus and validator clients

#### reset

| Name                | Description                       | Argument | Default                    |
| ------------------- | --------------------------------- | -------- | -------------------------- |
| --geth-datadir      | geth datadir                      | Path     | "./mainnet-data/execution" |
| --prysm-datadir     | prysm datadir                     | Path     | "./mainnet-data/consensus" |
| --validator-datadir | validator datadir                 | Path     | "./mainnet-data/validator" |
| --mainnet           | Run for mainnet (default network) | Bool     | false                      |
| --testnet           | Run for testnet                   | Bool     | false                      |
| --devnet            | Run for devnet                    | Bool     | false                      |

#### validator

<!--- TODO: please decide with format you think it works best. The info in this field is duplicated -->

`lukso validator --deposit` - Path to your deposit file - makes a deposit to a deposit contract
`lukso validator --genesis-deposit` - Path to your genesis deposit file - makes a deposit to genesis validator contract
`lukso validator --rpc` - Your RPC provider example: "https://rpc.2022.l16.lukso.network"  
`lukso validator --gas-price` -
`lukso validator --max-txs-per-block` - Maximum amount of txs sent per single block

| Name                | Description                                                                       | Argument | Default                              |
| ------------------- | --------------------------------------------------------------------------------- | -------- | ------------------------------------ |
| --deposit           | Path to your deposit file - makes a deposit to a deposit contract                 | Path     | ""                                   |
| --genesis-deposit   | Path to your genesis deposit file - makes a deposit to genesis validator contract | Path     | ""                                   |
| --rpc               | Your RPC provider                                                                 | URL      | "https://rpc.2022.l16.lukso.network" |
| --gas-price         | Gas price provided by user                                                        | Int      | 1000000000                           |
| --max-txs-per-block | Maximum amount of txs sent per single block                                       | Int      | 10                                   |

#### validator init

| Name                                   | Description                           | Argument | Default              |
| -------------------------------------- | ------------------------------------- | -------- | -------------------- |
| --validator-wallet-dir value           | location of generated wallet          | Path     | "./mainnet-keystore" |
| --validator-keys-dir value             | Path to your validator keys           | Path     |                      |
| --validator-wallet-password-file value | Path to your password file            | Path     |                      |
| --mainnet                              | Run for mainnetFlag (default network) | Bool     | false                |
| --testnet                              | Run for testnetFlag                   | Bool     | false                |
| --devnet                               | Run for devnet                        | Bool     | false                |

#### version

=============

## PR Testing ( MacOS/Linux )

For PR builds please use the PR gh deployment for example

```sh
curl https://install.lukso.network/25 | sh
```

where 25 is the sample pull request ID. You can also directly use the PR preview URL mentioned inside
of the PR status similar to:

```sh
curl https://lukso-network.github.io/tools-lukso-cli/pr-preview/pr-25 | sh
```

### start

geth commands & prysm commands

                                                 |

<!---
| --geth-datadir | A path of geth's data directory | Path | ./execution_data |
| --geth-ws | Enable WS server | None | true |
| --geth-ws-addr | Address of WS server | IP Address | 0.0.0.0 |
| --geth-ws-origins | Origins to accept requests from | WS Origins OR wildcard | \* |
| --geth-ws-apis | Comma separated apis | String of apis separated by comma | "net,eth,debug,engine" |
| --geth-bootnodes | Bootnode addresses | Bootnode addresses | See [Bootnodes](#bootnodes) |
| --geth-networkid | Network ID | Integer | 2022 |
| --geth-nat | Sets HTTP nat to assign static IP for geth | Example: "extip:0.0.0.0" | extip:83.144.95.19 |
| --geth-http | Enable HTTP server | None | true |
| --geth-http-apis | Comma separated apis | String of apis separated by comma | "net,eth,debug,engine,txlookup" |
| --geth-http-addr | Address used in HTTP comunication | IP address | 0.0.0.0 |
| --geth-http-corsdomain | Origins to accept requests from | HTTP Origins OR wildcard | \* |
| --geth-http-vhosts | Geth's virtual hosts | Virtual hostnames OR wildcard | \* |
| --geth-ipcdisable | Disable IPC communication | None | True |
| --geth-ethstats | URL of ethstats service | URL | "" |
| --geth-metrics | Enable metrics system | None | True |
| --geth-metrics-addr | Address of service managing collected metrics | IP Address | 0.0.0.0 |
| --geth-syncmode | Sync mode | Sync mode level | full |
| --geth-gcmode | Garbage colelction mode | Garbage collection level | archive |
| --geth-tx-look-up-limit | Number of blocks to maintain tx indexes from | Integer | 1 |
| --geth-cache-preimages | Enable cache preimaging | None | True |
| --geth-verbosity | Verbosity for geth logging | Verbosity level | 3 |
| --geth-port | Geth's port | Port | 30405 |
| --geth-http-port | Geth's HTTP port | Port | 8565 |
| --geth-mine | Enable mining | None | True |
| --geth-miner-threads | Number of CPU threads used for mining | Integer | 1 |
| --geth-miner-gaslimit | Gas ceiling | Integer | 60000000 |
| --geth-miner-etherbase | Your ECDSA public key used to get rewards on geth chain | Public address | 0x0000000000000000000000000000000000000000 |
| --geth-auth-jwt-secret | Path to JWT 32-byte secret | Path | ./config/mainnet/shared/secrets/jwt.hex |
| --geth-std-output | Set output to console | None | False |
| --geth-output-dir | Directory where logs are created | Path | ./logs/execution/geth _/

<!---

| --prysm-genesis-state | Genesis state file path | Path | ./config/mainnet/shared/genesis.ssz |
| --prysm-datadir | A path of prysm's beacon chain data directory | Path | ./consensus_data |
| --prysm-execution-endpoint | Execution endpoint | URL | http://localhost:8551 |
| --prysm-bootstrap-nodes | Bootnode addresses | Bootnode addresses | See [Bootnodes](#bootnodes) |
| --prysm-jwt-secret | Path to JWT 32-byte secret | Path | ./config/mainnet/shared/secrets/jwt.hex |
| --prysm-suggested-fee-recipient | Address that receives block fees | Public address | 0x0000000000000000000000000000000000000000 |
| --prysm-min-sync-peers | Minimum sync peers number for prysm | Integer | 0 |
| --prysm-p2p-host | P2P host IP | IP address | Empty |
| --prysm-deposit-deployment | Deployemnt height of deposit contract | Integer | 0 |
| --prysm-chain-config-file | Path to config.yaml file | Path | ./config/mainnet/shared/config.yaml |
| --prysm-monitoring-host | Host used for interacting with Prometheus metrics | IP address | 0.0.0.0 |
| --prysm-grpc-gateway-host | Host for gRPC gateway | IP address | 0.0.0.0 |
| --prysm-rpc-host | RPC server host | IP address | 0.0.0.0 |
| --prysm-verbosity | Verbosity for Prysm logs | Log level | info |
| --prysm-p2p-max-peers | Max peers for prysm | Integer | 250 |
| --prysm-subscribe-all-subnets | Subscribe to all possible subnets | None | True |
| --prysm-minimum-peers-per-subnet | Minimum peers per subnet | Integer | 0 |
| --prysm-enable-rpc-debug-endpoints | Enable debugging RPC endpoints | None | True |
| --prysm-output-dir | Directory where logs are created | Path | ./logs/consensus/beacon_chain |
| --prysm-std-output | Set output to console | None | False |
-->

#### Bootnodes

### download

| Name               | Description                                           | Argument                    |
| ------------------ | ----------------------------------------------------- | --------------------------- |
| --accept-terms     | Accept Terms provided by clients you want to download | None                        |
| --geth-tag         | Tag of geth's version that you want to download       | Tag, ex. `1.0.0`            |
| --geth-commit-hash | Commit hash that matches provided tag commit          | Commit Hash, ex. `12345678` |
| --validator-tag    | Tag of validator's version that you want to download  | Tag, ex. `v1.0.0`           |
| --prysm-tag        | Tag of prysm's version that you want to download      | Tag, ex. `v1.0.0`           |

Note difference in tags between geth and prysm/validator (`v` at the beginning)

### update

| Name            | Description                                          | Argument          |
| --------------- | ---------------------------------------------------- | ----------------- |
| --geth-tag      | Tag of geth's version that you want to download      | Tag, ex. `1.0.0`  |
| --validator-tag | Tag of validator's version that you want to download | Tag, ex. `v1.0.0` |
| --prysm-tag     | Tag of prysm's version that you want to download     | Tag, ex. `v1.0.0` |

### log

### reset

### validator

## Generate bindings

### Prerequisites:

- solc (https://github.com/ethereum/solidity)
- abigen (https://geth.ethereum.org/docs/tools/abigen)

### Steps

1. Paste your smart contract that you want to interact with into [`./contracts`](./contracts) directory
2. Generate ABI from your smart contract:

```bash
$ solcjs --output-dir abis --abi contracts/depositContract.sol
```

3. Generate bindings using newly generated ABI

```bash
abigen --abi abis/your-abi-file --pkg bindings --out contracts/bindings/yourBindingFile.go --type TypeName
```

4. To use binding in code type in:

```go
bind, err := bindings.NewTypeName(common.HexToAddress(contractAddress), ethClient)
if err != nil {
	return
}

tx, err := bind.DoSomething(...)
```
