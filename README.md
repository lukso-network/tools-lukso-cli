# LUKSO CLI

The LUKSO CLI is a command line tool to install, manage and set up validators of different clients for the LUKSO Blockchain. For more information and tutorial, please check the [LUKSO Tech Docs](https://docs.lukso.tech/).

## Features

- 🧰 Installation of Execution, Consensus, and Validator Clients
- 📀 Running a node as a validator
- 📑 Accessing various client logs

## Supported EVM Clients

The LUKSO CLI is able to install multiple clients for running the node.

- Execution Clients: [Geth](https://geth.ethereum.org/), [Erigon](https://github.com/ledgerwatch/erigon)
- Consensus Clients: [Prysm](https://github.com/prysmaticlabs/prysm), [Lighthouse](https://github.com/sigp/lighthouse), [Teku](https://github.com/Consensys/teku)
- Validator Staking Clients: [Prysm](https://docs.prylabs.network/docs/how-prysm-works/prysm-validator-client), [Lighthouse](https://github.com/sigp/lighthouse), [Teku](https://github.com/Consensys/teku)

### Client versions

| Client     | Version  | Release                                                      |
| ---------- |----------|--------------------------------------------------------------|
| Geth       | v1.13.5  | https://github.com/ethereum/go-ethereum/releases/tag/v1.13.5 |
| Erigon     | v2.55.0  | https://github.com/ledgerwatch/erigon/releases/tag/v2.55.0   |
| Prysm      | v4.0.8   | https://github.com/prysmaticlabs/prysm/releases/tag/v4.0.8   |
| Lighthouse | v4.5.0   | https://github.com/sigp/lighthouse/releases/tag/v4.5.0       |
| Teku       | v23.11.0 | https://github.com/Consensys/teku/releases/tag/23.11.0       |

> More clients will be added in the future.

## Supported Platforms

The LUKSO CLI is officially supported on Mac, Ubuntu, and Debian with the following architectures:

- `x86`/`x86_64`: Intel and AMD Processors
- `ARM`/`aarch64`: Single-Board Computers as M1 or Raspberry

> The experience might differ with other setups or versions of these operating systems.

## Setting up the Node

Process of setting up the node using the LUKSO CLI.

### Installing the LUKSO CLI

- Download and execute the LUKSO CLI installation script
- Running this script will install the full LUKSO CLI on Mac and Linux
- Installation directory: `/usr/local/bin/lukso`

```sh
# Might need admin access by typing `sudo` in front of the command
$ curl https://install.lukso.network | sh
```

### Initialise the Working Directory

1. Create and move into a working directory for your node's data

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

3. Install choosen LUKSO clients into the working directory

```sh
# Installing Execution, Consensus, and Validator Client
# Might need admin access by typing `sudo` in front of the command
$ lukso install
```

### Starting the Clients

Please refer to the `start` command below for more information.

## Working Directories' Structure

As the LUKSO CLI is able to manage multiple clients for multiple blockchain networks in one folder, the structure of the node is set up in a generic way.

- When initializing the node (with `lukso init`), a global configuration folder is created, which holds shared and unique client information for each type of network.
- When executing commands, directories for the associated network type will be created accordingly.

Network Types: `mainnet`, `testnet`

> Even if multiple networks are set up, only one can be active simultaneously, as the LUKSO CLI runs natively within your system environments.

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
├───[network_type]-keystore                 // Network's Validator Wallet
│   ├───keys                                // Encrypted Private Keys
│   ├───...                                 // Files for Signature Creation
|   ├───pubkeys.json                        // Validator Public Keys
|   ├───deposit_data.json                   // Deposit JSON for Validators
|   └───node_config.yaml                    // Node Configuration File
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

| Command                                     | Description                                                                                                     |
|---------------------------------------------|-----------------------------------------------------------------------------------------------------------------|
| [`install` ](#install)                      | Installs choosen clients (Execution, Consensus, Validator) and their binary dependencies                        |
| [`update` ](#update)                        | Update all currently selected clients to their newest versions                                                  |
| [`update configs` ](#update-configs)        | Update chain configuration files. This commands overwrites your oldchain configs, but keeps your client configs |
| [`init`](#initialise-the-working-directory) | Initializes the working directory, it's structure, and network configuration                                    |
| [`start`](#start)                           | Starts all or specific clients and connects to the specified network                                            |
| [`stop`](#stop)                             | Stops all or specific clients that are currently running                                                        |
| [`logs`](#logs)                             | Listens to and logs all events from a specific client in the current terminal window                            |
| [`status`](#status)                         | Shows the client processes that are currently running                                                           |
| [`status peers`](#status-peers)             | Shows the peer count of your node                                                                               |
| [`reset`](#reset)                           | Resets all or specific client data directories and logs excluding the validator keys                            |
| [`validator import`](#validator-import)     | Import the validator keys in the wallet                                                                         |
| [`validator list`](#validator-list)         | Display the imported validator keys                                                                             |
| [`validator exit`](#validator-exit)         | Issue an exit for your validator                                                                                |
| [`version`](#version)                       | Display the version of the LUKSO CLI that is currently installed                                                |
| [`help`, `h`](#help)                        | Shows the full list of commands, global options, and their usage                                                |

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

#### How to install the clients

```sh
# User is able to select its Consensus and Execution clients.
# Detects pre-installed clients and will ask for overrides
$ lukso install

# Installs clients and agrees with Terms & Conditions automatically
$ lukso install --agree-terms
```

#### Options for `install`

| Option        | Description                               |
| ------------- | ----------------------------------------- |
| --agree-terms | Automatically accept Terms and Conditions |

### `update`
`update` will install the newest verions of the clients that you selected duing installation process.  

#### How to update clients
```sh
# starts an update of your selected clients - note that your node must be STOPPED before proceeding with update
$ lukso update

```

### `update configs`
`update configs` Update chain configuration files. This commands overwrites your oldchain configs, but keeps your client configs.  
In order to update your configs, you need to: 
1) Have your node stopped
2) Be in the LUKSO CLI initialized directory

#### How to update chain config files
```sh
# starts an update of chain config files. 
$ lukso update configs

```

### `start`

`start` will automatically start clients selected during the installation process and connect them to the network.  
Please note that Teku needs a JDK installed in order to operate - this can be installed during the client installation process.

#### How to start the clients

```sh
# Starts your node and connects to LUKSO mainnet
# Uses the default config files from configs/mainnet folder
$ lukso start

# Starts your node and connects to mainnet as a validator
$ lukso start --validator --transaction-fee-recipient "0x1234..."

# Starts your node and connects to the LUKSO testnet
$ lukso start --testnet

# Starts your node and connects to testnet as a validator
$ lukso start --testnet --validator --transaction-fee-recipient "0x1234..."
```

#### How to start a node using config files

```sh
# Geth Configutation
# Change [config] to the name of your configuration file
$ lukso start --geth-config "./[config].toml"

# Prysm Configutation
# Change [config] to the name of your configuration file
# Change [custom_bootnode] to the bootnode's name
$ lukso start --prysm-config "./[config].yaml" --geth-bootnodes "[custom_bootnode]"

# Lighthouse Configutation
# Change [config] to the name of your configuration file
# Change [custom_bootnode] to the bootnode's name
$ lukso start --lighthouse-config "./[config].yaml" --lighthouse-bootnodes "[custom_bootnode]"

```

#### How to set up and customize a logs folder

```sh
# Setting up a custom logs directory
# Change [folder path] to a static or dynamic directory path
$ lukso start --logs-folder "[folder_path]"
```

#### Using Checkpoint Syncing

Checkpoint synchronization is a feature that significantly speeds up the initial sync time of the consensus client. If enabled, your node will begin syncing from a recently finalized consensus checkpoint instead of genesis. It will then download the rest of the blockchain data while your consensus is already running.

> After the synchronization is finalized, you will end up with the equal blockchain data. You can use the flag on every startup. However, it shows the most significant effect when synchronizing from scratch or after an extended downtime. The shortcut is ideal for fresh installations, validator migration, or recovery.

##### Checkpoints with LUKSO CLI version 0.8

```sh
# Mainnet Checkpoint Sync for Consensus Client
$ lukso start --checkpoint-sync

# Testnet Checkpoint Sync for Consensus Client
$ lukso start --testnet --checkpoint-sync
```

LUKSO CLI takes advantage of a weak subjectivity checkpoint flag (varies across different clients) that allows you to specify a weak subjectivity checkpoint.
With this flag specified, your beacon node will ensure that it reconstructs a historical chain that matches the checkpoint root at the given epoch.
This can offer the same level of weak subjectivity protection that checkpoint sync offers.
The CLI will automatically retrieve the latest finalized values to use with this feature.

##### Checkpoints with LUKSO CLI version 0.7 or below

Visit the [Mainnet Checkpoint Explorer](https://checkpoints.mainnet.lukso.network/) and get the latest block root and epoch. Then input both values into the flags below.

```sh
# Mainnet Checkpoint for Prysm Consensus Client
$ lukso start --prysm-checkpoint-sync-url=https://checkpoints.mainnet.lukso.network \
--prysm-genesis-beacon-api-url=https://checkpoints.mainnet.lukso.network/ \
--prysm-weak-subjectivity-checkpoint=$<BLOCK_ROOT>:$<EPOCH>

# Mainnet Checkpoint for Lighthouse Consensus Client
$ lukso start --lighthouse-checkpoint-sync-url=https://checkpoints.mainnet.lukso.network \
--lighthouse-genesis-beacon-api-url=https://checkpoints.mainnet.lukso.network/ \
--lighthouse-weak-subjectivity-checkpoint=$<BLOCK_ROOT>:$<EPOCH>

# Testnet Checkpoint for Prysm Consensus Client
$ lukso start --testnet --prysm-checkpoint-sync-url=https://checkpoints.testnet.lukso.network

# Testnet Checkpoint for Lighthouse Consensus Client
$ lukso start --testnet --lighthouse-checkpoint-sync-url=https://checkpoints.testnet.lukso.network
```

#### Options for `start`

| Option                               | Description                                                                                                                                           |
| ------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| **NETWORK**                          |                                                                                                                                                       |
| --mainnet                            | Starts the LUKSO node with mainnet data [default] (./configs/mainnet)                                                                                 |
| --testnet                            | Starts the LUKSO node with testnet data (./configs/tesnet)                                                                                            |
| --devnet                             | Starts the LUKSO node with devnet data (./configs/devnet)                                                                                             |
| **VALIDATOR**                        |                                                                                                                                                       |
| --validator                          | Starts the validator client                                                                                                                           |
| --validator-keys [string]            | Directory of the validator keys (default: "./\[network_type\]-keystore")                                                                              |
| --validator-wallet-password [string] | Location of password file that you used for generated validator keys                                                                                  |
| --validator-config [string]          | Path to prysms validator.yaml config file                                                                                                             |
| --transaction-fee-recipient [string] | The address that receives block reward from transactions (required for --validator flag)                                                              |
| --genesis-json [string]              | The path to genesis JSON file                                                                                                                         |
| --genesis-ssz [string]               | The path to genesis SSZ file                                                                                                                          |
| --no-slasher                         | Disables slasher                                                                                                                                      |
| **CLIENT OPTIONS**                   |                                                                                                                                                       |
| --logs-folder [string]               | Sets up a custom logs directory (default: "./\[network_type\]-logs")                                                                                  |
| --geth-config [string]               | Defines the path to geth TOML config file                                                                                                             |
| --prysm-config [string]              | Defines the path to prysm YAML config file                                                                                                            |
| --erigon-config [string]             | Defines the path to erigon TOML config file                                                                                                           |
| --teku-config [string]               | Defines the path to teku YAML config file                                                                                                             |
| --validator-config [string]          | Defines the path to teku validator YAML config file                                                                                                   |
| --geth-[command]                     | The `command` will be passed to the Geth client. [See the client docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --prysm-[command]                    | The `command` will be passed to the Prysm client. [See the client docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighhouse-[command]                | The `command` will be passed to the Lighthouse client. [See the client docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --erigon-[command]                   | The `command` will be passed to the Erigon client. [See the client docs for details](https://github.com/ledgerwatch/erigon)                           |
| --teku-[command]                     | The `command` will be passed to the Teku client. [See the client docs for details](https://github.com/ledgerwatch/erigon)                             |
| --checkpoint-sync                    | Run a node with checkpoint sync feature                                                                                                               |

### `stop`

#### How to stop all or specific clients

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

#### Options for `stop`

| Option      | Description                |
| ----------- | -------------------------- |
| --validator | Stops the validator client |
| --execution | Stops the execution client |
| --consensus | Stops the consensus client |

### `logs`

#### How to view logs of the clients

```sh
# Displays and saves the logs of the mainnet's consensus client
$ lukso logs consensus

# Displays and saves the logs of the devnet's execution client
$ lukso logs execution --devnet

# Displays and saves the logs of the testnet's validator
$ lukso logs validator --testnet
```

#### Options for `logs`

| Option    | Description                                                       |
| --------- | ----------------------------------------------------------------- |
| --mainnet | Logs the mainnet client [default] (./mainnet-logs/[client_type]/) |
| --testnet | Logs the testnet client (./testnet-logs/[client_type]/)           |
| --devnet  | Logs the devnet client (./devnet-logs/[client_type]/)             |

### `status`

#### How to check the status of the node

```sh
# Shows which client processes are currently running
$ lukso status
```

### `status peers`

#### How to get your peer count

```sh
# Shows the peer count of your node
$ lukso status peers
```

Ensure that the appropriate API is enabled when starting the node, as not all clients enable peer querying by default. For specific information about each client's peers queries, visit their documentation:

- [Geth Peer Interaction](https://geth.ethereum.org/docs/interacting-with-geth/rpc/ns-admin)
- [Erigon Peer Commands](https://github.com/ledgerwatch/erigon/blob/devel/cmd/rpcdaemon/README.md#rpc-implementation-status)

All supported consensus clients follow the [Beacon API](https://ethereum.github.io/beacon-APIs/#/Node/getPeers) standardization.

### `reset`

#### How to reset the node's data directory

```sh
# Resets LUKSO's mainnet data directory
$ lukso reset

# Resets LUKSO's testnet data directory
$ lukso reset --testnet

# Resets LUKSO's devnet data directory
$ lukso reset --devnet
```

#### Options for `reset`

| Option    | Description                           |
| --------- | ------------------------------------- |
| --mainnet | Resets LUKSO's mainnet data [default] |
| --testnet | Resets LUKSO's testnet data           |
| --devnet  | Resets LUKSO's devnet data            |

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
$ lukso start -h
```

## Running a Validator

Validator keys can be generated using:

- CLI: [tools-key-gen-cli](https://github.com/lukso-network/tools-key-gen-cli)
- GUI: [wagyu-key-gen](https://github.com/lukso-network/tools-wagyu-key-gen)

> Both of them will generate the validator keys directory.

After generating the validator keys, they can be imported into the LUKSO CLI. To fill the validator keys with funds to participate on the LUKSO Blockchain, you must use the [LUKSO Launchpad](https://deposit.mainnet.lukso.network) to send LYXe from your wallet to the generated keys.

#### Genesis Amounts

- All genesis validators will be prompted to vote for the initial token supply of LYX
- The initial token supply will determine how much LYX the Foundation will receive
- More details at: https://deposit.mainnet.lukso.network

#### Validator Stake

- Genesis Validators need to have at least 32 LYXe per validator
- Validators also need some ETH on their wallet to pay for gas expenses

### `validator import`

Import existing EIP-2335 keystore files (such as those generated by the [Wagyu Keygen](https://github.com/lukso-network/tools-wagyu-key-gen)) into Prysm.

#### How to import validator keys

For `validator import` command the `--validator-keys` flag is required.

```sh
# Regular import process
# You will be asked for the passwords - of your validator keys and newly generated wallet

lukso validator import --validator-keys "./myKeysDir"

# You will be asked only for wallet password
lukso validator import --validator-keys "./myKeysDir" --validator-password "./myKeysPasswordFile.txt"
```

#### Options for `validator import`

| Option                        | Description                                                                 |
| ----------------------------- | --------------------------------------------------------------------------- |
| --validator-keys [string]     | Directory of the validator keys (default: "./\[network_type\]-keystore")    |
| --validator-password [string] | Path to validator keys' password file                                       |
| **NETWORK**                   |                                                                             |
| --mainnet                     | Will import the keys for mainnet [default] (Will use: "./mainnet-keystore") |
| --testnet                     | Will import the keys for testnet (Will use: "./testnet-keystore")           |
| --devnet                      | Will import the keys for devnet (Will use: "./devnet-keystore")             |

### `validator list`

List all imported keys from the validators keystore.

### `validator exit`

Issue an exit for your validator keys that have set withdrawal credentials.

Running this command will take you to interactive version of exit command implemented by the respective consensus client developers (Running validator exit command while running Prysm consensus layer will take you to Prysm exit process etc.)

> ETH1 withdrawal addresses are required for the exit to work. These are automatically generated by the [Wagyu Key-Gen Tool](https://github.com/lukso-network/tools-wagyu-key-gen). If you used BLS keys from the [Key Gen CLI](https://github.com/lukso-network/tools-key-gen-cli) without specifying a withdrawal address, please update your withdrawal address credentials first.

**IMPORTANT:** The Validator exit is an **irreversible** action. Before exiting your validator you need make sure you are mindful of what the exit process carries with itself. #

Because different clients have different processes of exiting please make sure that **you are familiar with the flags** provided by the `validator exit` command.

#### How to exit validator keys

Make sure your validator node is running before starting the exit command.

```sh
# Exit validator keys with Prysm/Teku for Mainnet
sudo lukso validator exit

# Exit validator keys with Prysm/Teku for Testnet
sudo lukso validator exit --testnet

# Exit a validator key with Lighthouse for Mainnet - please note that you may issue an exit for only a single validator at a time
sudo lukso validator exit --keystore "./mainnet-keystore/keystore-xxx.json"

# Exit a validator key with Lighthouse for Testnet
sudo lukso validator exit --testnet --keystore "./testnet-keystore/keystore-xxx.json"
```

Note that each client that you use may have different exit process - you can read more about those on client's official documentation:

- Prysm: https://docs.prylabs.network/docs/wallet/exiting-a-validator
- Lighthouse: https://lighthouse-book.sigmaprime.io/voluntary-exit.html
- Teku: https://docs.teku.consensys.net/how-to/voluntarily-exit

> The Lighthouse client only allows to exit one validator key at a time. In case you have plenty of keys, please generate a separate keyfolder in a different working directory using Prysm. The new keystore folder can then be used for the exit command, even if your node is running the Lighthouse consensus client.

| Option                          | Description                                                                                                                        |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| --keystore [string]             | Path to keystore-xxx.json file containing public key that you want to exit - this flag is required when using Lighthouse validator |
| --validator-wallet-dir [string] | Path to wallet containing validators that you want to exit (default: "./mainnet-keystore" - is affected by the network flag)       |
| --rpc-address [string]          | Address of node that is used to make an exit (defaults to the default RPC address provided by your selected client)                |
| --testnet-dir value             | Path to network configuration folder (default: "./configs/mainnet/shared" - is affected by the network flag)                       |
| **NETWORK**                     |                                                                                                                                    |
| --mainnet                       | Will import the keys for mainnet [default] (Will use: "./mainnet-keystore")                                                        |
| --testnet                       | Will import the keys for testnet (Will use: "./testnet-keystore")                                                                  |
| --devnet                        | Will import the keys for devnet (Will use: "./devnet-keystore")                                                                    |

#### How to import validator keys

```sh
# List imported keys for mainnet keystore
# You will be asked for the password
lukso validator list

# List imported keys for devnet keystore
# You will be asked for the password
lukso validator list --devnet
```

#### Options for `validator list`

| Option      | Description                                                              |
| ----------- | ------------------------------------------------------------------------ |
| **NETWORK** |                                                                          |
| --mainnet   | Will list the keys for mainnet [default] (default: "./mainnet-keystore") |
| --testnet   | Will list the keys for testnet (default: "./testnet-keystore")           |
| --devnet    | Will list the keys for devnet (default: "./devnet-keystore")             |

For specific validator options, please visit the [Prysm Validator Specification](https://docs.prylabs.network/docs/wallet/nondeterministic). All flags and their parameters will be passed to the client. This can be useful to configure additional features like the validator's graffiti, extended logging, or connectivity options.

### Starting the Validator

#### How to start your validator (keys & tx fee recipient)

When you use `--validator`, the `--transaction-fee-recipient` flag is required.

```sh
# Specify the transaction fee recipient, also known as coinbase
# It is the address where the transactions fees are sent to
$ lukso start --validator --transaction-fee-recipient "0x12345678..."
```

If you want to provide a specific keystore directory, you can use the `--validator-keys` flag. If no `--validator-keys` flag was specified, the CLI will look in the default directory: `./[network_type]-keystore`.:

```sh
# Validator keys
# Command split across multiple lines for readability
# Change [file_name] with the your password text file's name
$ lukso start --validator \
--transaction-fee-recipient "0x12345678..." \
--validator-keys "./custom-keystore-dir-path"
```

#### How to start a genesis node

```sh
# Command split across multiple lines for readability
# Make sure that both SSZ and JSON files are placed correctly
$ lukso start \
--genesis-ssz "./config/mainnet/shared/genesis.ssz" \
--genesis-json "./config/mainnet/shared/genesis.json"
```

## Uninstalling

The LUKSO CLI and downloaded clients are located within the binary folder of the user's system directory.
It can be removed at any time. All node data is directly located within the working directory.

```sh
# Make sure to stop the node
$ lukso stop

# Uninstall the LUKSO CLI
$ rm -rf /usr/local/bin/lukso

# Uninstall Geth Execution Client
$ rm -rf /usr/local/bin/geth

# Uninstall Prysm Consensus Client
$ rm -rf /usr/local/bin/prysm

# Uninstall Lighthouse Consensus Client
$ rm -rf /usr/local/bin/lighthouse

# Uninstall Erigon Execution Client
$ rm -rf /usr/local/bin/erigon

# Remove the node data
# Make sure to backup your keys first
$ rm -rf ~/myNodeFolder
```

## Contributing

If you want to contribute to this repository, please check [`CONTRIBUTING.md`](./CONTRIBUTING.md).
