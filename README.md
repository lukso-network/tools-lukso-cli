# LUKSO CLI

> ⚠️ DO NOT USE YET, this is WIP!

The LUKSO Command Line Interface (lukso-cli) serves the following purposes:

- easy installation of all node types (full installs into `/bin/` , not docker containers)
- easy starts and stops local nodes (as it runs as a daemon)
- easy access to nodes logs
- running a node as a validator
- making validator deposits

## Repository Structure

- [`./cmd/lukso`](./cmd/lukso): code of LUKSO CLI
- [`./abis`](./abis) - collection of ABIs from smart contracts that are being interacted with
- [`./contracts`](./contracts) - collection of said smart contracts
- [`./contracts/bindings`](./contracts/bindings) - bindings generated from ABIs - to generate new bindings see [Generate bindings](#generate-bindings) section.
- [`./install`](./install/) - collection of things to support the various installation, signing and notarization requirements.
- [`./docs`](./docs) Some small content to be inserted into mac pgk file

## CLI Installation
```sh
curl https://install.lukso.network | sh
```

## Downloading and Installing LUKSO

```bash
# 1. To install the LUKSO Command Line Interface (CLI)
$ curl https://install.lukso.network | sh

# 2. Create a working folder where you want your clients to store their data
$ mkdir myLUKSOFolder && cd ./myLUKSOFolder

# 3. This will initialize your working folder by downloading all network configs from
# https://github.com/lukso-network/network-configs
# NOTE: This will not overwrite any existing config, data or keystore folders
$ lukso init

# 4. Install your clients. It will ask you which ones you want to install
$ lukso install

# If you want to auto accept terms run it with
$ lukso install --agree-terms
```

## Available parameters

`lukso <command>  [--flags]`

| Command   | Description                                                             |
| --------- | ----------------------------------------------------------------------- |
| init      | Initializes with the network configuration files                        |
| install   | Installs clients globally                                               |
| log       | List logs from the different clients                                    |
| reset     | Resets data directories                                                 |
| start     | Starts your installed clients and connects it to the respective network |
| status    | Shows your current process running                                      |
| stop      | Stops all or specific client(s)                                         |
| update    | Sets client(s) to desired version                                       |
| validator | Init and deposits your validator keys                                   |
| version   | Display version of lukso-cli                                            |

## Initializing your working folder

```bash
# Running the init command will initialize your working folder by downloading the network configs
# NOTE: This will not overwrite any existing config, data or keystore folders
$ mkdir myLUKSOFolder && cd ./myLUKSOFolder

# Initalize LUKSO
$ lukso init
```

  Network configs: [github.com/lukso-network/network-configs(https://github.com/lukso-network/network-configs)

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

| Flag                                                                                  | Description                                                                                                                 |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --testnet                                                                             | Starts LUKSO's testnet                                                                                                      |
| --devnet                                                                              | Starts LUKSO's devnet                                                                                                       |
| --geth-\*                                                                             | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon-\*                                                                           | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm-\*                                                                            | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse-\*                                                                       | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --geth-config                                                                         | Path to "./myconfig.toml" file                                                                                              |
| --prysm-config "./myconfig.yaml" --geth-bootnodes "mycustombootnode00000"             | Path to "./myconfig.yaml" file & custom geth boot nodes                                                                     |
| --validator --transaction-fee-recipient                                               | Address that receives block fees (0x12345..abcd).                                                                           |
| --validator --validator-keys "./mainnet-keystore" --validator-password "./myfile.txt" | Passes the validator keys and password from a custom directory                                                              |
| --log -folder "./myCustomLogFolder"                                                   | Sets up a custom log directory when starting lukso-cli                                                                      |

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

| Flag                             | Description                                                                      |
| -------------------------------- | -------------------------------------------------------------------------------- |
| --deposit                        | Path to your deposit file. Makes a deposit to a deposit contract                 |
| --genesis-deposit                | Path to your genesis deposit file; makes a deposit to genesis validator contract |
| --rpc                            | Your RPC provider (URL) - "https//rpc.2022.l16.lukso.network"                    |
| --gas-price                      | Gas price provided by user (int) 1000000000                                      |
| --max-txs-per-block              | Maximum amount of txs sent per single block (int) 10                             |
| --validator-wallet-dir           | Location of a generated wallet "./mainnet/keystore"                              |
| --validator-keys-dir             | Path to your validator keys                                                      |
| --validator-wallet-password-file | Path to your password file                                                       |

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
