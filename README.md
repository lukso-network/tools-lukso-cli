# LUKSO CLI

> ⚠️ DO NOT USE IN PRODUCTION, SCRIPTS ARE NOT DEPLOYED YET.

The LUKSO CLI is a command line tool to install, manage and set up validators of different types of nodes for the LUKSO network. For more information and tutorial, please check the [LUKSO Tech docs](https://docs.lukso.tech/).

## Features

- 🧰 Installation of Execution, Beacon and Validator Node Clients
- 📀 Running a node as a validator
- 📑 Accessing various node logs
- 💰 Making validator deposits

## Supported EVM Clients

> WIP: More clients will be added

The LUKSO CLI is able to install multiple clients for running the node.

- Execution Node Clients: [Geth](https://geth.ethereum.org/)
- Beacon Node Clients: [Prysm](https://github.com/prysmaticlabs/prysm)
- Validator Client for Staking: [Prysm](https://docs.prylabs.network/docs/how-prysm-works/prysm-validator-client)

## Setting up the Node

Process of setting up the node using the LUSKO CLI.

### Installing the LUKSO CLI

- Download and execute the LUKSO CLI installation script
- Running this script will install the full LUKSO CLI on Mac and Linux
- Installation directory: `/usr/local/bin/lukso`

```sh
# Might need admin access by typing `sudo` in front of the command
$ curl https://install.lukso.network | sh
```

### Initialise the node folder

1. Create and move into a working directory for your node client data

```sh
# Exchange [folder_name] with the folder name you want
$ mkdir [folder_name] && cd ./[folder_name]
```

2. Initialize the working directory

```sh
# Downloads all network configs from https://github.com/lukso-network/network-configs
$ lukso init
```

### Installing the Clients

3. Install choosen LUKSO node clients into the working directory

```sh
# Installing Execution Chain, Beacon Chain, and Validator Client
# Might need admin access by typing `sudo` in front of the command
$ lukso install
```

### Starting the Clients

Please refer to the `start` command below for more information.

## Node folder structure

As the LUKSO CLI is able to manage multiple clients for multiple blockchain networks in one folder, the structure of the node is set up in a generic way.

- When initializing the node (with `lukso init`), a global configuration folder is created, which holds shared and unique client information for each type of network.
- When executing commands, directories for the associated network type will be created accordingly.

Network Types: `mainnet`, `testnet`, `devnet`

> Even if multiple networks are set up, only one can be active at the time

```
lukso-node
│
├───configs                                 // Configuration
│   └───[network_type]                      // Network's Config Data
│       ├───shared
|       |   ├───genesis.json                // Genesis JSON Data
|       |   ├───genesis.ssz                 // Genesis Validator File
|       |   └───config.yaml                 // Global Client Config
│       ├───geth                            // Config for Geth Client
│       ├───prysm                           // Config for Prysm Client
│       ├───erigon                          // Config for Erigon Client
│       └───lighthouse                      // Config for Lighthouse Client
│
├───[network_type]-keystore                 // Network's Validator Data
│   ├───keys                                // Encrypted Private Keys
│   ├───...                                 // Files for Signature Creation
|   ├───pubkeys.json                        // Validator Public Keys
|   ├───deposit_data.json                   // Deposit JSON for Validators
|   └───node_config.yaml                    // Node Configuration File
|
├───[network_type]-wallet                   // Network's Transaction Data
|
├───[network_type]-data                     // Network's Blockchain Data
│   ├───consensus                           // Storage of used Consensus Client
│   ├───execution                           // Storage of used Execution Client
│   └───validator                           // Storage of Validator Client
│
├───[network_type]-logs                     // Network's Logged Data
|
└───cli-config.yaml                         // Global CLI Configuration
```

## Available Commands

| Command             | Description                                                                                      |
| ------------------- | ------------------------------------------------------------------------------------------------ |
| `install`           | Installs choosen LUKSO node clients (Execution, Beacon, Validator) and their binary dependencies |
| `init`              | Initializes the node working directory, it's structure, and network configuration                |
| `update`            | Updates all or specific LUKSO node clients in the working directory to the newest version        |
| `start`             | Starts all or specific LUKSO node clients and connects to the specified network                  |
| `stop`              | Stops all or specific LUKSO node clients that are currently running                              |
| `log`               | Listens to all log events from a specific client in the current terminal window                  |
| `status`            | Shows the LUKSO node client processes that are currently running                                 |
| `reset`             | Resets all or specific client data directories and logs excluding the validator keys             |
| `validator import`  | Import the validator keys in the wallet                                                          |
| `validator deposit` | Makes a deposit to the deposit bridge contract                                                   |
| `version`           | Display the version of the LUKSO CLI that is currently installed                                 |
| `help`, `h`         | Shows the full list of commands, global options, and their usage                                 |

## Global Help Flag

| Flag       | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| --help, -h | Can be added before or after a command to show it's help page |

## Examples and Explanations

For almost each command in the list, options can be added after it to modify or specify certain behavior.
Below, you can find examples and options tables for all available commands.

> Options containting [string] expects a string input in quotes.

> Options containting [int] expects an int input without quotes.

### `install`

| Option        | Description                               |
| ------------- | ----------------------------------------- |
| --agree-terms | Automatically accept Terms and Conditions |

#### How to install the node clients

```sh
# User is able to select its Consensus and Execution clients.
# Detects pre-installed clients and will ask for overrides
$ lukso install

# Installs clients and agrees with Terms & Conditions automatically
$ lukso install --agree-terms
```

### `update`

> Options WIP

#### How to update clients

```sh
# Updates installed clients
$ lukso update

# Update Geth client to a latest version
$ lukso update geth
```

### `start`

| Option                               | Description                                                                  |
| ------------------------------------ | ---------------------------------------------------------------------------- |
| --geth-config [string]               | Defines the path to geth TOML config file                                    |
| --prysm-config [string]              | Defines the path to pryms YAML config file                                   |
| --genesis-json [string]              | Defines the path to genesis JSON file                                        |
| --genesis-ssz [string]               | Defines the path to genesis SSZ file                                         |
| --log-folder [string]                | Sets up a custom log directory (default: "./\[network_type\]-logs")          |
| --no-slasher                         | Disables slasher                                                             |
| **VALIDATOR**                        |                                                                              |
| --validator                          | Starts the validator client                                                  |
| --transaction-fee-recipient [string] | Sets the address that receives block fees                                    |
| --validator-keys [string]            | Directory of the validator keys (default: "./\[network_type\]-keystore")     |
| --validator-password [string]        | Location of password file that you used for generation keys from deposit-cli |
| --validator-config [string]          | Path to prysm.yaml config file                                               |
| **NETWORK**                          |                                                                              |
| --mainnet                            | Starts the LUKSO node with mainnet data (default) (./configs/mainnet)        |
| --testnet                            | Starts the LUKSO node with testnet data (./configs/tesnet)                   |
| --devnet                             | Starts the LUKSO node with devnet data (./configs/devnet)                    |

For specific client options, please visit their official documentations:

- [Geth Client Specification](https://geth.ethereum.org/docs/fundamentals/command-line-options)
- [Prysm Client Specification](https://docs.prylabs.network/docs/prysm-usage/parameters)
- [Erigon Client Specification](https://github.com/ledgerwatch/erigon)
- [Lighthouse Client Specification](https://lighthouse-book.sigmaprime.io/advanced-datadir.html)

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

### `stop`

| Option      | Description                |
| ----------- | -------------------------- |
| --validator | Stops the validator client |
| --execution | Stops the execution client |
| --consensus | Stops the consensus client |

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

### `log`

> Options WIP

#### How to view logs of the node clients

```sh
# Displays the logs of the execution client
$ lukso log execution

# Displays the logs of the consensus client
$ lukso log consensus

# Displays the logs of the validator
$ lukso log validator
```

### `status`

> Options WIP

#### How to check the status of the node

```sh
# Shows which client processes are currently running
$ lukso status
```

### `reset`

| Option              | Description                 |
| ------------------- | --------------------------- |
| --mainnet [default] | Resets LUKSO's mainnet data |
| --testnet           | Resets LUKSO's testnet data |
| --devnet            | Resets LUKSO's devnet data  |

#### How to reset the node's data directory

```sh
# Resets LUKSO's mainnet data directory
$ lukso reset

# Resets LUKSO's testnet data directory
$ lukso reset --testnet

# Resets LUKSO's devnet data directory
$ lukso reset --devnet
```

### `version`

#### How to check the version of the LUKSO CLI

```sh
# Displays the currently installed version of the LUKSO CLI
$ lukso version
```

### `help`

In addition to the help command, the global help flag can be used to generate help pages for commands

#### How to retrieve the help page in the CLI

```sh
# Displays the global help page
$ lukso help

# Displays the help page of the start command
$ lukso start --help

# Displays the help page of the start command
$ lukso -h start
```

## Running a Validator

Deposit keys can be generated using:

- CLI: [tools-key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli)
- GUI: [wagyu-key-gen](https://github.com/lukso-network/tools-wagyu-key-gen)

#### Genesis Amounts

- All genesis validators will be prompted to vote for the initial token supply of LYX
- The initial token supply will determie how much LYX the Foundation will receive
- More details at: https://deposit.mainnet.lukso.network

#### Validator Stake

- Genesis Validators need to have at least 32 LYXe per validator
- Validators in general also need some ETH to pay for gas expenses

### `validator import`

| Option               | Description                                                                |
| -------------------- | -------------------------------------------------------------------------- |
| --validator-keys     | Path to your validator keys directory [required]                           |
| --validator-password | Path to your password file                                                 |
| **NETWORK**          |                                                                            |
| --mainnet            | Will import the keys for mainnet [default] (default: "./mainnet-keystore") |
| --testnet            | Will import the keys for testnet (default: "./testnet-keystore")           |
| --devnet             | Will import the keys for devnet (default: "./devnet-keystore")             |

The validator keys directory is generated with one of these tools:

- CLI: [tools-key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli)
- GUI: [wagyu-key-gen](https://github.com/lukso-network/tools-wagyu-key-gen)

### `validator deposit`

| Option                       | Description                                                                                                                 |
| ---------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --deposit-data-json [string] | Defines the path to the deposit JSON file                                                                                   |
| --gas-price [string]         | Defines the gas price in integers as string                                                                                 |
| --rpc [string]               | Defines the RPC URL used to send deposit transactions. If `--genesis` is used, it should be an RPC URL for ethereum mainnet |
| --genesis                    | Executes deposit to genesis validator contract                                                                              |
| --start-from-index [int]     | Start deposit from specific block index                                                                                     |

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

## Contributing

If you want to contribute to this repository, please check [`CONTRIBUTING.md`](./CONTRIBUTING.md).
