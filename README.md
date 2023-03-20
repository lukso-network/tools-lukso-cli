# LUKSO CLI

> ⚠️ DO NOT USE YET, this is WIP!

The `lukso` CLI serves the following purposes:

- easy installation of all node types (full installs into `/bin/` , not docker containers)
- easy starts and stops local nodes (as it runs as a daemon)
- easy access to nodes logs
- running a node as a validator
- making validator deposits

## Repository struct

- [`./cmd/lukso`](./cmd/lukso): code of LUKSO CLI
- [`./abis`](./abis) - collection of ABIs from smart contracts that are being interacted with
- [`./contracts`](./contracts) - collection of said smart contracts
- [`./contracts/bindings`](./contracts/bindings) - bindings generated from ABIs - to generate new bindings see [Generate bindings](#generate-bindings) section.

## Installation ( Linux/MacOS )

> 🛠️ Work In Progress, available soon.

## Installation ( Windows )

> 🛠️ Work In Progress, available soon.

## Installing LUKSO

```bash
# 1. The command below installs the LUKSO client and prompts the user with default settings to get started as a validator
$ curl https://install.lukso.network | sh

# 2. You will need to agree to the lukso-cli terms before continuing. Simply type "Y" to agree
https://github.com/lukso-network/tools-lukso-cli/blob/main/TERMS.md

# 3. Create a working folder where you want your clients to store their data:
$ mkdir myLUKSOFolder && cd ./myLUKSOFolder

# 4. This command downloads the network configs from
# https://github.com/lukso-network/network-configs to your recently created "myLUKSOFolder/configs"
# It will not overwrite any existing config folders, data or keystore the user might have
$ lukso init

# 5. Install your desired clients using
$ lukso install
# Or simply type to accept all Terms & Conditions of LUKSO's clients.
# You don't need to run this if you already installed LUKSO.
$ lukso install --agree-terms

```

## Available parameters

`lukso <command> [geth, prysm, validator, *all*] [--flags]`

> _all_ means that you can skip an argument for all possible options to run (default, only option for download)

| Command        | Description                                  |
| -------------- | -------------------------------------------- |
| init           | Initializes configuration files              |
| install        | Downloads all default client(s)              |
| log            | Show logs                                    |
| reset          | Resets data directories                      |
| start          | Starts up all or specific client(s)          |
| status         | Shows status of all or specified client      |
| stop           | Stops all or specific client(s)              |
| update         | sets client(s) to desired version            |
| validator      | Manages validator-related commands           |
| validator init | Initializes your validator with deposit keys |
| version        | Display version of LUKSO CLI                 |

## Starting your node as a Genesis Validator

```bash
# The LUKSO start command for Genesis validators should be run as the following:
$ lukso start --genesis-ssz "./config/mainnet/shared/genesis.ssz" --genesis-json "./config/mainnet/geth/genesis.json"
```

## Starting your node

```bash

# starts your currently installed default clients and connects to LUKSO mainnet
$ lukso start

# starts your nodes connecting to the testnet
$ lukso start --testnet

# starts your nodes connecting to the mainet as a validator
# use default keystore folder (/mainnet-keystore)
$ lukso start --validator

```

The following flags are available to the `$ lukso start` command:

| Flag                                 | Description                                                                                               |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| --mainnet                            | Run for LUKSO mainnet (default network)                                                                   |
| --testnet                            | Run for testnet                                                                                           |
| --validator                          | Starts your node as validator node                                                                        |
| --validator -datadir                 | A path of validator's data directory (./validator_data)                                                   |
| --validator -verbosity               | Verbosity for validator logs                                                                              |
| --validator -wallet-dir              | Location of a generated wallet (./mainet_keystore)                                                        |
| --validator -wallet-password-file    | Location for password used for wallet generation (./config/mainnet/shared/secrets/validator-password.txt) |
| --validator -chain-config-file       | Path to config.yaml file (./config/mainnet/shared/config.yaml)                                            |
| --validator -monitoring-host         | Host used for interacting with Prometheus metrics (IP Address)                                            |
| --validator -grpc-gateway-host       | Host for gRPC gateway (IP address)                                                                        |
| --validator -rpc-host                | RPC server host (IP address)                                                                              |
| --validator -suggested-fee-recipient | Address that receives block fees (0x12345..abf)                                                           |
| --validator -output-dir              | Directory where the logs are created (./logs/consensus/validator)                                         |
| --validator -std-output              | Set output to console                                                                                     |

## LUKSO CLI Clients Flags

All the flags included in each of the clients (geth, erigon, prysm and lighthouse) are also available and can be passed as valid commands.
Documentation for other client's CLIs can be found:

- Geth: https://geth.ethereum.org/docs/fundamentals/command-line-options
- Erigon: https://github.com/ledgerwatch/erigon
- Prysm: https://docs.prylabs.network/docs/prysm-usage/parameters
- Lighthouse: https://lighthouse-book.sigmaprime.io/advanced-datadir.html

## Update

```bash

# updates lukso client to the latest available version
$ lukso update

# updates to the specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso update --geth-tag

#  updates a specific version of the validator
$ lukso update --validator-tag
```

## Log

```bash
# displays the logs of LUKSO's Execution client
$ lukso log execution

#displays the LUKSO's consensus client's logs
$ lukso log consensus

#displays the LUKSO's validator client's logs
$ lukso log validator

#displays the testnet client's logs
$ lukso log --testnet

#displays the devnet client's logs
$ lukso log --devnet

#Log files are by default un-capped, be aware that these files can grow very large.
#You can use:
$ lukso start --log-size n
#to cap the file size in n MB. Example: lukso start --log-size 3     to cap the file to 3MB of data

```

| Name                    | Description                                     | Argument | Default          |
| ----------------------- | ----------------------------------------------- | -------- | ---------------- |
| --geth-output-file      | Path to geth log file that you want to log      | Path     | "./mainnet-logs" |
| --prysm-output-file     | Path to prysm log file that you want to log     | Path     | "./mainnet-logs" |
| --validator-output-file | Path to validator log file that you want to log | Path     | "./mainnet-logs" |
| --mainnet               | Run for mainnet (default network)               | Bool     | false            |
| --testnet               | Run for testnet                                 | Bool     | false            |
| --devnet                | Run for devnet                                  | Bool     | false            |

## reset

| Name                | Description                       | Argument | Default                    |
| ------------------- | --------------------------------- | -------- | -------------------------- |
| --geth-datadir      | geth datadir                      | Path     | "./mainnet-data/execution" |
| --prysm-datadir     | prysm datadir                     | Path     | "./mainnet-data/consensus" |
| --validator-datadir | validator datadir                 | Path     | "./mainnet-data/validator" |
| --mainnet           | Run for mainnet (default network) | Bool     | false                      |
| --testnet           | Run for testnet                   | Bool     | false                      |
| --devnet            | Run for devnet                    | Bool     | false                      |

## validator

| Name                | Description                                                                       | Argument | Default                              |
| ------------------- | --------------------------------------------------------------------------------- | -------- | ------------------------------------ |
| --deposit           | Path to your deposit file - makes a deposit to a deposit contract                 | Path     | ""                                   |
| --genesis-deposit   | Path to your genesis deposit file - makes a deposit to genesis validator contract | Path     | ""                                   |
| --rpc               | Your RPC provider                                                                 | URL      | "https://rpc.2022.l16.lukso.network" |
| --gas-price         | Gas price provided by user                                                        | Int      | 1000000000                           |
| --max-txs-per-block | Maximum amount of txs sent per single block                                       | Int      | 10                                   |

## validator init

| Name                                   | Description                           | Argument | Default              |
| -------------------------------------- | ------------------------------------- | -------- | -------------------- |
| --validator-wallet-dir value           | location of generated wallet          | Path     | "./mainnet-keystore" |
| --validator-keys-dir value             | Path to your validator keys           | Path     |                      |
| --validator-wallet-password-file value | Path to your password file            | Path     |                      |
| --mainnet                              | Run for mainnetFlag (default network) | Bool     | false                |
| --testnet                              | Run for testnetFlag                   | Bool     | false                |
| --devnet                               | Run for devnet                        | Bool     | false                |

## Extras - Other Installation Options & Custom Commands

```bash

# downloads lukso client accepting terms provided by clients you want to download
$ lukso install --agree-terms
# downloads a specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso install --geth-tag
# downloads a specific tagged commit of (geth/prysm/erigon/lighthouse) client- Example Geth v1.11.4
$ lukso install --geth-commit-hash
#  downloads a specific version of the validator
$ lukso install --validator-tag
```

```bash

#User can set up a custom log directory when starting lukso client
$ lukso start --log-folder "./myCustomLogFolder"
#in this case, to access their logs user needs to indicate the folder
$ lukso log --log-folder "./myCustomLogFolder"
```

Note difference in tags between geth and prysm/validator (`v` at the beginning)

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
