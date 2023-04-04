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
└───cmd                   // Commands
│   └───lukso             // LUKSO CLI
│
└───contracts             // Solidity Contracts
│   └───bindings          // Bindings generated from ABIs
│   └───abis              // Binary Interfaces of LUKSO Smart Contracts
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
**Installation directory: `/usr/local/bin/lukso`**

## Node Folder Structure

> Initializing a LUKSO node will not overwrite existing config, data or keystore folders

```
lukso-node
│
└───configs                 		// Blockchain Configuration
│   └───mainnet 	        		// Mainnet Config
│   |   └───shared
|   |   |   | genesis.json  		// Genesis JSON Data
|   |   |   | genesis.ssz   		// Genesis Validator File
|   |	|	| config.yaml   		// Network Configuration
│   |   └───geth	        		// Storage of Geth Client
│   |   └───prysm	      		 	// Storage of Prysm Client
│   |   └───erigon	        		// Storage of Erigon Client
│   |   └───lighthouse	    		// Storage of Lighthouse Client
|   |
│   └───testnet			    		// Testnet Config
│   |   └───...		    			// Similar to Mainnet Config
|   |
│   └───devnet    		    		// Devnets Config
│       └───...		    			// Similar to Mainnet Config
│
└───mainnet-validator        		// Mainnet Validator Secrets and Keys
|   └───mainnet-transaction-wallet  // Validators Transaction Wallet
|   └───mainnet-keystore        	// Mainnet Validator Secrets and Keys
│   |   └───keys               		// Encrypted Private Keys
│   |   └───...                 	// Folders & Files for Signature Creation
|   |   | pubkeys.json          	// Validator Public Keys
|	| deposit_data.json         	// Deposit JSON for Validators
| 	| node_config.yaml          	// Node Configuration File
|
└───mainnet-data            		// Mainnet Data Storage
│   └───consensus_data      		// Storage of used Consensus Client
│   └───execution_data      		// Storage of used Execution Client
│   └───validator_data      		// Storage of Validator Client
│
└───mainnet-logs            		// Mainnet Log Data
|
└───testnet-validator        		// Testnet Validator Secrets and Keys
|   └───...  						// Similar to Mainnet Validator
|
└───testnet-data            		// Testnet Data Storage
|   └───...  						// Similar to Mainnet Data
│
└───testnet-logs            		// Testnet Log Data
|
└───devnet-validator        		// Devnet Validator Secrets and Keys
|   └───...  						// Similar to Mainnet Validator
|
└───devnet-data            			// Devnet Data Storage
|   └───...  						// Similar to Mainnet Data
│
└───devnet-logs            			// Devnet Log Data
|
| cli-config.yaml           		// LUKSO CLI Configuration
```

## External Sources

- The network configuration is fetched from [lukso-network/network-configs](https://github.com/lukso-network/network-configs)

- Deposit keys can be generated using [lukso-network/tools-key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli)

- Tool for generating validator keys [wagyu-key-gen](https://github.com/lukso-network/tools-wagyu-key-gen)

## Client Clarification

> WIP: More client setups will be added

The LUKSO CLI is able to install multiple node clients.
They cover the full node functionality of an EVM PoS Blockchain.

- Supported Execution Node Clients: Geth
- Supported Beacon Node Clients: Prysm
- Validator Client for Staking

## Setting up the Node

Process of setting up the node using the LUSKO CLI Tool

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

### `install`

| Tag                | Used in Commands | Description                              |
| ------------------ | ---------------- | ---------------------------------------- |
| geth-tag           | install          | Installs Geth client to latest version   |
| prysm-tag          | install          | Installs Prysm client to latest version  |
| geth-tag [string]  | install          | Installs Geth client to certain version  |
| prysm-tag [string] | install          | Installs Prysm client to certain version |

### `init`

### `update`

| Tag                     | Description                                  |
| ----------------------- | -------------------------------------------- |
| geth-tag                | Updates Geth client to latest version        |
| prysm-tag               | Updates Prysm client to latest version       |
| geth-tag [string]       | Updates Geth client to certain version       |
| prysm-tag [string]      | Updates Prysm client to certain version      |
| erigon-tag              | Updates Erigon client to latest version      |
| lighthouse-tag          | Updates Lighthouse client to latest version  |
| erigon-tag [string]     | Updates Erigon client to certain version     |
| lighthouse-tag [string] | Updates Lighthouse client to certain version |

### `start`

| Flag                                 | Description                                       |
| ------------------------------------ | ------------------------------------------------- |
| --validator                          | Starts the validator client                       |
| --execution                          | Starts the execution client                       |
| --consensus                          | Starts the consensus client                       |
| --geth                               | Starts Geth client                                |
| --prysm                              | Starts Prysm client                               |
| --lighthouse                         | Starts Lighthouse client                          |
| --erigon                             | Starts Erigon client                              |
| --mainnet                            | Starts LUKSO's mainnet data                       |
| --testnet                            | Starts LUKSO's testnet data                       |
| --devnet                             | Starts LUKSO's devnet data                        |
| --geth-config                        | Defines the path to TOML config file              |
| --prysm-config [string]              | Defines the path to the YAML config file          |
| --geth-bootnodes [string]            | Sets a custom Geth bootnode name                  |
| --transaction-fee-recipient [string] | Sets the address that receives block fees         |
| --validator-keys [string]            | Passes the validator keys from a custom directory |
| --validator-password [string]        | Passes the assword from a custom directory        |
| --log-folder [string]                | Sets up a custom log directory                    |
| --no-slasher                         | Disables slasher                                  |
| --genesis-json [string]              | Defines the path to genesis JSON file             |
| --genesis-ssz [string]               | Defines the path to genesis SSZ file              |

### `stop`

| Flag        | Description                |
| ----------- | -------------------------- |
| --validator | Stops the validator client |
| --execution | Stops the execution client |
| --consensus | Stops the consensus client |

### `log`

### `status`

### `reset`

| Flag      | Description                 |
| --------- | --------------------------- |
| --mainnet | Resets LUKSO's mainnet data |
| --testnet | Resets LUKSO's testnet data |
| --devnet  | Resets LUKSO's devnet data  |

### `validator init`

Initializes the validator with keys.

| Flag                             | Description                                 |
| -------------------------------- | ------------------------------------------- |
| --max-txs-per-block              | Maximum amount of txs sent per single block |
| --validator-wallet-dir           | Location of a generated wallet              |
| --validator-keys-dir             | Path to your validator keys                 |
| --validator-wallet-password-file | Path to your password file                  |

### `validator deposit`

Makes a deposit to the deposit bridge contract.

| Flag                             | Description                                    |
| -------------------------------- | ---------------------------------------------- |
| --deposit-data-json [string]     | Defines the path to the deposit JSON file      |
| --gas-price [string]             | Defines the gas price in integers as string    |
| --rpc [string]                   | Defines the RPC URL on deposit                 |
| --genesis                        | Executes deposit to genesis validator contract |
| --start-from-index [int]         | Start deposit from specific block index        |
| --max-txs-per-block              | Maximum amount of txs sent per single block    |
| --validator-wallet-dir           | Location of a generated wallet                 |
| --validator-keys-dir             | Path to your validator keys                    |
| --validator-wallet-password-file | Path to your password file                     |

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

For Client Tags, please visit their official documentations:

- [Geth Client Specification](https://geth.ethereum.org/docs/fundamentals/command-line-options)
- [Prysm Client Specification](https://docs.prylabs.network/docs/prysm-usage/parameters)
- [Erigon Client Specification](https://github.com/ledgerwatch/erigon)
- [Lighthouse Client Specification](https://lighthouse-book.sigmaprime.io/advanced-datadir.html)

## General Examples and Explanations

#### How to view logs of the node clients

```sh
# Displays the logs of execution client
$ lukso log execution

# Displays the consensus client's logs
$ lukso log consensus

# Displays the validator client's logs
$ lukso log validator
```

#### How to install the node clients

```sh
# User is able to select its Consensus and Execution clients.
# Detects pre-installed clients and will ask for overrides
# in case install is called multiple times.
$ lukso install

# Installs clients and agrees with Terms & Conditions automatically
$ lukso install --accept-terms-of-use
```

#### How to call help interface of an command

```sh
# Call help interface for installation
$ lukso install --help

# Call help interface for starting the node
$ lukso start --help
```

#### How to reset the node's data directory

```sh
# Resets LUKSO's mainnet data directory
$ lukso reset

# Resets LUKSO's testnet data
$ lukso reset --testnet

# Resets LUKSO's devnet data
$ lukso reset --devnet

```

#### How to start the node clients

```sh
# Starts your node clients and connects to LUKSO mainnet
# Uses the default config files from configs/mainnet folder
$ lukso start

# Starts your node clients and connects to the LUKSO testnet
$ lukso start --testnet

# Starts your node clients and connects to the LUKSO devnet
$ lukso start --devnet

# Starts your node clients and connects to mainnet as a validator
$ lukso start --validator
```

#### How to start a genesis validator node

```sh
# Example using Geth client and its folder
# Command split across multipl lines for readability
# Make sure that both SSZ and JSON files are placed correctly
$ lukso start \
--genesis-ssz "./config/mainnet/shared/genesis.ssz" \
--genesis-json "./config/mainnet/geth/genesis.json"
```

#### How to start your validator (keys & tx fee recipient)

```sh
# Start your node as a validator node
$ lukso start --validator

# Specify the transaction fee recipient, also knows as coinbase
# Its the address where the transactions fees are sent to
$ lukso start --validator --transaction-fee-recipient  "0x12345678..."

# Validator keys and password
# Command split across multipl lines for readability
# Change [file_name] with the your password text file's name
$ lukso start --validator \
--validator-keys "./mainnet-keystore" \
--validator-password "./[file_name].txt"
```

#### How to start a node using config files

```sh
# Geth Configutation
# Change [config] to the name of your configuration file
$ lukso start --geth-config "./[config].toml"

# Prysm Configutation
# Change [config] to the name of your configuration file
# Change [custom_bootnode] to the bootnode's name
$ lukso start --prysm-config "./[config].yaml" \
--geth-bootnodes "[custom_bootnode]"

# An experienced user can also start custom clients
# Example with Lighthouse and Erigon clients
$ lukso start --lighthouse --erigon
```

#### How to set up and customize a log folder

```sh
# Setting up a custom log directory
# Change [folder path] to a static or dynamic directory path
$ lukso start --log-folder "[folder_path]"
```

#### How to check the status of the node

```sh
# Shows which client processes are currently running
$ lukso status
```

#### How to stop all or specific node clients

```sh
# Stops all running node clients
$ lukso stop

# Only stops the validator client
$ lukso stop --validator

# Only stops the execution client
$ lukso stop --execution

# Only stops the consensus client
$ lukso stop --consensus
```

#### How to update lukso-cli

```sh
# Updates installed clients
$ lukso update

# Update Geth client to a latest version
$ lukso update geth-tag

# Update Geth client to version 1.0.0
$ lukso update geth --tag "v1.0.0"
```

#### Checking the version of the LUKSO CLI

```sh
# Displays the currently installed version of the LUKSO CLI
$ lukso version
```

## Running a Validator

#### Genesis Amounts

- All Genesis Validators will be prompted to vote for the initial token supply of LYX
- The initial token supply will determie how much LYX the Foundation will receive
- More details at: https://deposit.mainnet.lukso.network

#### Validator Stake

- Genesis Validators need to have at least 32 LYXe per validator
- Validators in general also need some ETH to pay for gas expenses

### Performing Deposits

The following example snippets show how to perform deposits with your keys to be able to stake.

#### How to deposit as a validator

> WIP: Inconsistant Object Usage

```sh
# Executing deposit by setting gas price and RPC connection
# Change [rpc_address] to the actual RPC URL
$ lukso validator deposit \
--deposit-data-json "./deposit_data.json" \
[--gasPrice "1000000000" --rpc "[rpc_address]"]
```

#### How to deposit as a genesis validator

> WIP: Inconsistant Object Usage

```sh
# Genesis validator's deposits setting gas price and an RPC connection
# Change [rpc_address] to the actual RPC URL
# Change N to specific block index
$ lukso validator deposit \
--genesis \
--deposit-data-json "./deposit-data.json" \
--rpc "[rpc_address]" [--gas-price "1000000000" --start-from-index N]
```

## Development & Testing

### Install CLI from Codebase

#### Download the LUKSO CLI Repository

```sh
# Clone tools-lukso-cli repository
git clone git@github.com:lukso-network/tools-lukso-cli.git
```

#### Install Go Language

Head over to the [Official Go Installation Page](https://go.dev/doc/install) and follow along.

After Installation, check if everything is set up correctly by querying the version.

```sh
# Check the installed Go version
go version
```

#### Build the Executable

```sh
# Build the local project within the tools-lukso-cli repository
cd cmd/lukso/ && go build -o lukso
```

#### Run the generated LUKSO CLI

```sh
# Within the tools-lukso-cli repository
# Alternatively a static path can be used
cd cmd/lukso/
./lukso <command>
```

### Bindings

#### Prerequisites:

- Solidity Compiler [solc.js](https://github.com/ethereum/solidity)
- ABI Generator [abigen](https://geth.ethereum.org/docs/tools/abigen)

#### Generate Bindings

1. Paste the interacting smart contract into the [`contracts`](./contracts) directory
2. Generate the ABIs from your smart contracts

```sh
# Generating ABI into abis directory
# Change [contract_name] to the contract file name
$ solcjs \
--output-dir contracts/abis \
--abi contracts/[contract_name].sol
```

3. Generate bindings using the newly generated ABIs

```sh
# Generate Go Bindings for CLI
# Change [abi_name] to the contract file name
# Change [binding_name] to the wanted output file name
abigen \
--abi contracts/abis/[abi_name] \
--pkg bindings \
--out contracts/bindings/[binding_name].go \
--type TypeName
```

4. Use generated bindings in code

```go
# Sample Implementation
bind, err := bindings.NewTypeName(common.HexToAddress(contractAddress), ethClient)
if err != nil {
	return
}

tx, err := bind.DoSomething(...)
```

### PR Testing

This repository has a CI/CD pipeline set up that will automatically build a script to install a new version of the LUKSO CLI globally.

- Should only be used for testing
- Will overwrite the LUKSO CLI that is currently installed

> For PR builds please use the separate PR's GH deployment

#### Using the LUKSO CLI URL

```sh
# Might need admin access by typing `sudo` in front of the command
curl https://install.lukso.network/36 | sh
```

> 36 is the sample pull request ID that can be changed

#### Using PR Preview URL

```sh
## Might need admin access by typing `sudo` in front of the command
# Using the live environment
curl https://lukso-network.github.io/tools-lukso-cli/pr-preview/pr-36 | sh
```

> 36 is the sample pull request ID that can be changed
