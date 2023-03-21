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

## init

This command downloads the network configs from https://github.com/lukso-network/network-configs to your recently created "myLUKSOFolder/configs"
It will not overwrite any existing config folders, data or keystore the user might have.

The init command should be run from the user's "myLUKSOFolder" directory.

```bash

$ cd myLUKSOFolder/

myLUKSOFolder $ lukso init
```

## install

Installs will install LUKSO's default clients. This command can also be used to install any other desired client of your choice (Geth/Erigon - Prysm/Lighthouse)

```bash
$ lukso install
```

```bash
# Or simply type to accept all Terms & Conditions of LUKSO's clients.
# You don't need to run this if you already installed LUKSO.
$ lukso install --agree-terms

```

| Flag             | Description                                                                                                                 |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --agree-terms    | installs LUKSO client and agrees with Terms & Conditions                                                                    |
| --geth -\*       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon -\*     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse -\* | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |

#### Examples

```bash
# downloads a specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso install --geth-tag
# downloads a specific tagged commit of (geth/prysm/erigon/lighthouse) client- Example Geth v1.11.4
$ lukso install --geth-commit-hash
#  downloads a specific version of the validator
$ lukso install --validator-tag
```

## log

Log displays the logs of LUKSO's execution/consensus/validator clients. Here are the common flags:

```bash
# displays the logs of LUKSO's Execution client
$ lukso log execution
# displays the LUKSO's consensus client's logs
$ lukso log consensus
# displays the LUKSO's validator client's logs
$ lukso log validator
```

| Flag             | Description                                                                                                                 |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --testnet        | displays client's testnet logs                                                                                              |
| --devnet         | displays client's devnet logs                                                                                               |
| --log-folder     | user can access their custom log folder "./myCustomLogFolder" [learn how to setup log folder]()                             |
| --geth -\*       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon -\*     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse -\* | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |

#### Examples

```bash
$ lukso log --log-folder "./myCustomLogFolder"
# Path to geth log file that you want to log
$ lukso log --geth-output-file "./mainnet-logs"
# Path to prysm log file that you want to log
$ lukso log --prysm-output-file "./mainnet-logs"
# Path to validator log file that you want to log
$ lukso log --validator-output-file "./mainnet-logs"
```

## reset

LUKSO reset will reset the mainnet data directory, not the keys

```bash
# resets LUKSO data
$ lukso reset
# resets LUKSO's testnet data
$ lukso reset --testnet
```

| Flag                  | Description                                                                                                                 |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --testnet             | resets the client's testnet                                                                                                 |
| --devnet              | resets the client's devnet                                                                                                  |
| --geth -\*            | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --geth -data-dir      | resets the "./mainnet/data/execution" directory                                                                             |
| --erigon -\*          | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*           | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --prysm -data-dir     | resets the "./mainnet/data/consensus" directory                                                                             |
| --lighthouse -\*      | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --validator -data-dir | resets the "./mainnet/data/validator" directory                                                                             |

## start

Starts your currently installed default clients and connects, by default, to LUKSO mainnet:

```bash

$ lukso start
```

```bash
# starts your nodes connecting to the testnet
$ lukso start --testnet

# starts your nodes connecting to the mainet as a validator
# use default keystore folder (/mainnet-keystore)
$ lukso start --validator
```

The LUKSO start command for Genesis validators should be run as the following:

```bash

$ lukso start --genesis-ssz "./config/mainnet/shared/genesis.ssz" --genesis-json "./config/mainnet/geth/genesis.json"
```

| Flag                                 | Description                                                                                                                 |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet                            | starts LUKSO's mainnet                                                                                                      |
| --testnet                            | starts LUKSO's testnet                                                                                                      |
| --devnet                             | starts LUKSO's devnet                                                                                                       |
| --geth -\*                           | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon -\*                         | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*                          | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse -\*                     | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --validator -data-dir                | A path of validator's data directory (./validator_data)                                                                     |
| --validator -verbosity               | verbosity for validator logs                                                                                                |
| --validator -wallet-dir              | location of a generated wallet (./mainnet_keystore)                                                                         |
| --validator -wallet-password-file    | location of password used for wallet generation (./config/mainnet/shared/secrets/validator-password.txt)                    |
| --validator -chain-config-file       | path to config.yaml file (./config/mainnet/shared/config.yaml)                                                              |
| --validator -monitor-host            | host used for interacting with Prometheus metrics (IP address)                                                              |
| --validator -grpc-gateway-host       | host for gRPC gateway (IP address)                                                                                          |
| --validator -rpc-host                | RPC server host (IP address)                                                                                                |
| --validator -suggested-fee-recipient | address that receives block fees (0x12345..abcd)                                                                            |
| --validator -output-dir              | directory where logs are created                                                                                            |
| --validator -std-output              | set output to console                                                                                                       |
| --log -folder "./myCustomLogFolder"  | user can setup a custom log directory when starting LUKSO client                                                            |
| --log -size n                        | Log files capped to the size in MB                                                                                          |

#### Examples:

```bash
# log files are by default un-capped, be aware that # these files can grow very large.
# You can use to cap the file size in n MB. Example: # lukso start --log-size 3     to cap the file to 3MB of data
$ lukso start --log-size n
# user can set up a custom log directory when starting lukso client
$ lukso start --log-folder "./myCustomLogFolder"
# in this case, to access their logs user needs to indicate the folder
```

## status

Displays the most recent status of LUKSO's node

```bash
$ lukso status
```

## stop

Stops all client's activities. usually used when upgrading the client or running maintenance tasks.

```bash
$ lukso stop
```

## update

Updates LUKSO client to the latest available version

```bash
$ lukso update
```

| Flag             | Description                                                                                                                 |
| ---------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet        | updates LUKSO's mainnet                                                                                                     |
| --testnet        | updates LUKSO's testnet                                                                                                     |
| --devnet         | updates LUKSO's devnet                                                                                                      |
| --geth -\*       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon -\*     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse -\* | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --validator      | updates a specific version of the validator                                                                                 |

#### Examples:

```bash
# updates to the specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso update --geth-tag

#  updates a specific version of the validator
$ lukso update --validator-tag
```

## validator

Starts your node as a validator node

```bash
$ lukso start --validator
```

| Flag                             | Description                                                                                                                 |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet                        | runs a validator on LUKSO's mainnet (default)                                                                               |
| --testnet                        | runs a validator on LUKSO's testnet                                                                                         |
| --devnet                         | runs a validator on LUKSO's devnet                                                                                          |
| --geth -\*                       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon -\*                     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm -\*                      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse -\*                 | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --deposit                        | path to your deposit file. Makes a deposit to a deposit contract                                                            |
| --genesis-deposit                | path to your genesis deposit file; makes a deposit to genesis validator contract                                            |
| --rpc                            | your RPC provider (URL) - "https//rpc.2022.l16.lukso.network"                                                               |
| --gas-price                      | Gas price provided by user (int) 1000000000                                                                                 |
| --max-txs-per-block              | Maximum amount of txs sent per single block (int) 10                                                                        |
| --validator-wallet-dir           | location of a generated wallet "./mainnet/keystore"                                                                         |
| --validator-keys-dir             | path to your validator keys                                                                                                 |
| --validator-wallet-password-file | path to your password file                                                                                                  |

## version

Displays the current version of your LUKSO client

```bash
$ lukso version
```

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
