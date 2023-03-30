# LUKSO CLI

> ⚠️ DO NOT USE IN PRODUCTION, SCRIPTS ARE NOT DEPLOYED YET.

The LUKSO CLI is a command line tool to install, manage and set up validators of different types of nodes for the LUKSO network.

## Features

- Installation of Execution, Beacon, and Validator Node Clients
- starting and stoping local nodes running as a daemon
- accessing various node logs
- running a node as a validator
- making validator deposits

## Repository Structure

```
tools-lukso-cli
│
└───abis                  // Binary Interfaces of LUKSO Smart Contracts
│
└───cmd                   // Commands
│   └───lukso             // LUKSO CLI
│
└───contracts             // Solidity Contracts
│   └───bindings          // Bindings generated from ABIs
|
└───docs                  // Updates for Installation Progress
|
└───install               // Mandatory Installation Tools
│   └───cf-wrangler       // Manager for Cloudflare Workers
│   └───docs-processor    // Markdown to Page Converter
│   └───macos_packages    // MacOS Codesigning Scripts
|
└───pid                   // Process ID Management
```

## CLI Installation Script

```sh
https://install.lukso.network
```

**Running this script will install the full LUKSO CLI Tool on Mac and Linux.**
**Installation directory: `usr/local/bin/lukso`**

## Node Folder Structure

> Initializing a LUKSO node will not overwrite existing config, data or keystore folders

```
lukso-node
│
└───configs               // Blockchain Configuration
|   | config.yaml         // Network Configuration
|   | genesis.json        // Genesis JSON Data
|   | genesis.ssz         // Genesis Validator File
│
└───data                  // Blockchain Data Storage
│   └───consensus_data    // Storage of Consensus Client
│   └───execution_data    // Storage of Execution Client
│   └───validator_data    // Storage of Validator Client
│
└───keystore              // Validator Secrets and Keys
│   └───keys              // Encrypted Private Keys
│   └───...               // Folders & Files for Signature Creation
|   | pubkeys.json        // Validator Public Keys
|
└───transaction_wallet    // Validators Transaction Wallet
|
| deposit_data.json       // Deposit JSON for Validators
| node_config.yaml        // Node Configuration File
```

## External Sources

- The network configuration is fetched from [lukso-network/network-configs](https://github.com/lukso-network/network-configs)

- Deposit key can be generated using [lukso-network/tools-key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli)

## Client Clarification

> WIP: More client setups will be added

The LUKSO CLI is able to install multiple node clients.
They cover the full node functionality of an EVM PoS Blockchain.

- Supported Execution Node Clients: Geth
- Supported Beacon Node Clients: Prysm
- Validator Client for Staking

## Setting up the Node

Process of setting up the node using the LUSKO CLI Tool

### Installing cURL

Installing a tool to fetch the LUKSO CLI Installation Script from the server.

#### MacOS

```sh
# Install the Homebrew package manager
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install cURL through Homebrew
brew install curl

# Check the cURL version
curl --version
```

#### Linux

```sh
# Reload the debian package list
sudo apt-get update -y

# Install cURL through debian package manager
sudo apt-get install curl -y

# Check the cURL version
curl --version
```

### Installing the LUKSO CLI Tool

Download and execute the LUKSO CLI Installation Script

```sh
# Might need admin access by typing `sudo` in front of the command
$ curl https://install.lukso.network | sh
```

### Setting up the Clients

1. Create and move into a working directory for your node client data

```sh
# Exchange [folder_name] with the folder name you want
$ mkdir [folder_name] && cd ./[folder_name]
```

2. Initialize the working directory

```sh
# Downloads all network configs from
$ lukso init
```

3. Install choosen LUKSO node clients into the working directory

```sh
# Installing Execution Chain, Beacon Chain, and Validator Client
# Might need admin access by typing `sudo` in front of the command
$ lukso install
```

## LUKSO CLI Usage

```sh
lukso [global options] [command] [command options] [arguments...] [global options]
```

> Global Options can be placed at the beginning or end of the terminal entry

## Available Commands

| Command   | Description                                                                                      |
| --------- | ------------------------------------------------------------------------------------------------ |
| install   | Installs choosen LUKSO node clients (Execution, Beacon, Validator) and their binary dependencies |
| init      | Initializes the node working directory, it's structure, and network configuration                |
| update    | Updates all or specific LUKSO node clients in the working directory to the newest version        |
| start     | Starts all or specific LUKSO node clients and connects to the specified network                  |
| stop      | Stops all or specific LUKSO node clients that are currently running                              |
| log       | Listens to all log events from a specific client in the current terminal window                  |
| status    | Shows the LUKSO node client processes that are currently running                                 |
| reset     | Resets all or specific client data directories and logs excluding the validator keys             |
| validator | Manages the LUKSO validator keys including their initialization and deposits                     |
| version   | Display the version of the LUKSO CLI Tool that is currently installed                            |
| help, h   | Shows the full list of commands, global options, and their usage                                 |

## Available Global Options

Global options can be added behind a command to allow different modifications to its execution.

| Global Option                    | Description                                        |
| -------------------------------- | -------------------------------------------------- |
| --accept-terms-of-use            | Automatically accept upcoming terms of use         |
| --help, --h, -h, -help, help, h, | Show help page of the command provided in the call |

## Available Flags

Flags can be added behind a command to allow further command specifications.

> Global Flags containting [string] are awaiting a string input in quotes.
> Global Flags containting [int] are awaiting a string input without quotes.

| Flag                                 | Used in Commands   | Description                                         |
| ------------------------------------ | ------------------ | --------------------------------------------------- |
| --validator                          | start, stop        | Starts or stops the validator client                |
| --execution                          | start, stop        | Starts or stops the execution client                |
| --consensus                          | start, stop        | Starts or stops the consensus client                |
| --geth                               | start              | Starts Geth client                                  |
| --prysm                              | start              | Starts Prysm client                                 |
| --lighthouse                         | start              | Starts Lighthouse client                            |
| --erigon                             | start              | Starts Erigon client                                |
| --mainnet                            | init, start, reset | Initializes, starts, or resets LUKSO's mainnet data |
| --testnet                            | init, start, reset | Initializes, starts, or resets LUKSO's testnet data |
| --devnet                             | init, start, reset | Initializes, starts, or resets LUKSO's devnet data  |
| --geth-config                        | start              | Defines the path to TOML config file                |
| --prysm-config [string]              | start              | Defines the path to the YAML config file            |
| --geth-bootnodes [string]            | start              | Sets a custom Geth bootnode name                    |
| --transaction-fee-recipient [string] | start              | Sets the address that receives block fees           |
| --validator-keys [string]            | start              | Passes the validator keys from a custom directory   |
| --validator-password [string]        | start              | Passes the assword from a custom directory          |
| --log-folder [string]                | start              | Sets up a custom log directory                      |
| --no-slasher                         | start              | Disables slasher                                    |
| --genesis-json [string]              | start              | Defines the path to genesis JSON file               |
| --genesis-ssz [string]               | start              | Defines the path to genesis SSZ file                |
| --deposit-data-json [string]         | validator deposit  | Defines the path to the deposit JSON file           |
| --gas-price [string]                 | validator deposit  | Defines the gas price in integers as string         |
| --rpc [string]                       | validator deposit  | Defines the RPC URL on deposit                      |
| --genesis                            | validator deposit  | Executes deposit to genesis validator contract      |
| --start-from-index [int]             | validator deposit  | Start deposit from specific block index             |
| --max-txs-per-block                  | validator          | Maximum amount of txs sent per single block         |
| --validator-wallet-dir               | validator          | Location of a generated wallet                      |
| --validator-keys-dir                 | validator          | Path to your validator keys                         |
| --validator-wallet-password-file     | validator          | Path to your password file                          |

## Available Subcommands

Subcommands and tags can be added behind commands to specify a certain function or dictate specific versioning.

| Subcommands | Superordinate Commands | Description                                     |
| ----------- | ---------------------- | ----------------------------------------------- |
| geth        | update, log, status    | Updates your Geth client to newest version      |
| prysm       | update, log, status    | Updates your Prysm client to newest version     |
| validator   | update, log, status    | Updates your Validator client to newest version |
| deposit     | validator              | Makes a deposit to the deposit bridge contract  |
| init        | validator              | Initializes the validator with keys             |

> When initializing the validator keys must have been generated using LUKSO's [key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli) tool

## How to install LUKSO CLI

```bash
# Installs the LUKSO CLI and prompts user to select its Consensus and Execution clients.
# Install also detects if you have any pre-installed client and confirms an override to a newer version in case needed.
$ lukso install

# Installs clients and agrees with Terms & Conditions automatically
$ lukso install --agree-terms
```

## How to view logs

```bash
# Displays the logs of execution client
$ lukso log execution

# Displays the consensus client's logs
$ lukso log consensus

# Displays the validator client's logs
$ lukso log validator
```

## How to reset your data directory

```bash
# Resets LUKSO mainnet data directory
$ lukso reset

# Resets LUKSO's testnet data
$ lukso reset --testnet
```

## How to start a node

```bash
# Starts your currently installed default clients and connects to LUKSO mainnet.
# Takes the default config files from the path "./config/mainnet/geth/config.toml"
$ lukso start

# Starts your nodes connecting to the testnet
$ lukso start --testnet

# Starts your nodes connecting to the devnet
$ lukso start --devnet

# Starts your nodes connecting to mainnet as a validator, using the default keystore folder (/mainnet-keystore)
$ lukso start --validator


# How to start a Genesis Validator node


# Start command for Genesis Validators should be run as the following:
$ lukso start --genesis-ssz "./config/mainnet/shared/genesis.ssz" --genesis-json "./config/mainnet/geth/genesis.json"


# How to start your validator (keys & tx fee recipient)


# Starts your node as a validator node
$ lukso start --validator

# The transaction fee recipient; aka coinbase, is the address where the transactions fees are sent to.
$ lukso start --validator --transaction-fee-recipient  "0x12345678..."

# Validator keys and password
$ lukso start --validator --validator-keys "./mainnet-keystore" --validator-password "./myfile.txt"


# How to start a node using config files


# Geth Configs
$ lukso start --geth-config "./myconfig.toml"

# Prysm Configs
$ lukso start --prysm-config "./myconfig.yaml" --geth-bootnodes "mycustombootnode0000"

# An experienced user might also want to start custom clients
$ lukso start --lighthouse --erigon


# How to set & customize a log folder


# Setting up a custom log directory
$ lukso start --log-folder "./myCustomLogFolder"
```

## How to check the status of LUKSO node

```bash
# Shows you which processes are running
$ lukso status
```

## How to stop LUKSO node

```bash
# Stops currently running clients
$ lukso stop

# Only stops the validator client
$ lukso stop --validator

# Only stops the execution client
$ lukso stop --execution

# Only stops the consensus client
$ lukso stop --consensus
```

## How to update lukso-cli

```bash
# Updates installed clients
$ lukso update

# Updates to the specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso update geth-tag
```

| Flag          | Description                                                                                                                 |
| ------------- | --------------------------------------------------------------------------------------------------------------------------- |
| geth-\*       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| erigon-\*     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| prysm-\*      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| lighthouse-\* | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |

## Running your validator

The main activity you can perform as a validator is depositing your keys.

#### How to deposit as a Validator and as a Genesis Validator

```bash
# Validator's deposits setting gas price and an RPC connection
$ lukso validator deposit --deposit-data-json "./validator-deposit-data.json" [--gasPrice "1000000000" --rpc "https://infura.io./apiv1"]

# Genesis validator's deposits setting gas price and an RPC connection
$ lukso validator deposit --genesis --deposit-data-json "./validator-deposit-data.json" --rpc "https://infura.io./apiv1" [--gas-price "1000000000" --start-from-index N]
```

All Genesis Validators will be prompted to vote for the initial token supply of LYX; determining how much the Foundation will receive. More details at: https://deposit.mainnet.lukso.network

Genesis Validators need to have at least 32 LYXe per validator and some ETH to pay for gas expenses.

## Checking the version

```bash
# Displays the current version of your lukso-cli
$ lukso version
```

## Development

## Generate bindings

### Prerequisites:

- solc (https://github.com/ethereum/solidity)
- abigen (https://geth.ethereum.org/docs/tools/abigen)

#### Steps

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
