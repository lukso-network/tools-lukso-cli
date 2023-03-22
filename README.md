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

## Downloading and Installing LUKSO

```bash
# 1. To install the LUKSO Command Line Interface  (CLI)
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

| Command   | Description                             |
| --------- | --------------------------------------- |
| init      | Initializes configuration files         |
| install   | Downloads all default client(s)         |
| log       | Show logs                               |
| reset     | Resets data directories                 |
| start     | Starts up all or specific client(s)     |
| status    | Shows status of all or specified client |
| stop      | Stops all or specific client(s)         |
| update    | Sets client(s) to desired version       |
| validator | Manages validator-related commands      |
| version   | Display version of LUKSO CLI            |

## Initializing your working folder

Running the init command will initialize your working folder by downloading the [network configs](https://github.com/lukso-network/network-configs)
NOTE: This will not overwrite any existing config, data or keystore folders

```bash
$ cd myLUKSOFolder

# inside the myLUKSOFolder run
$ lukso init
```

## How to install LUKSO clients

The install command prompts you to select which Execution (Geth or Erigon) and which Consensus (Prysm or Lighthouse) client you want to install.

Install also detects if you have any pre-installed client and confirms an override to a newer version in case needed.

```bash
$ lukso install
```

```bash
# Or simply type to accept all Terms & Conditions of LUKSO's clients.
# You don't need to run this if you already installed LUKSO.
$ lukso install --agree-terms
```

| Flag          | Description                                              |
| ------------- | -------------------------------------------------------- |
| --agree-terms | installs LUKSO client and agrees with Terms & Conditions |

## How to view logs

Displays the logs of LUKSO's execution/consensus/validator clients. Here are the common flags:

```bash
# displays the logs of LUKSO's execution client
$ lukso log execution
# displays the LUKSO's consensus client's logs
$ lukso log consensus
# displays the LUKSO's validator client's logs
$ lukso log validator
```

## How to reset your data directory

LUKSO reset will reset the mainnet data directory, not the keys

```bash
# resets LUKSO data
$ lukso reset
# resets LUKSO's testnet data
$ lukso reset --testnet
```

## How to start a node

Starts your currently installed execution and consensus clients and connects, by default, to LUKSO's mainnet. LUKSO start takes the default config files from the default path "./config/mainnet/geth/config.toml" for you.

```bash

$ lukso start
```

```bash
# starts and connects to the testnet
$ lukso start --testnet
# starts and connects to the devnet
$ lukso start --devnet
# starts and connects to mainnet as a validator, using the default keystore folder (/mainnet-keystore)
$ lukso start --validator
```

#### How Genesis Validators should start their nodes

The LUKSO start command for Genesis validators should be run as the following:

```bash

$ lukso start --genesis-ssz "./config/mainnet/shared/genesis.ssz" --genesis-json "./config/mainnet/geth/genesis.json"
```

#### How to start a node with your own config files

As an experienced validator; you might want to pass your own config files "./myconfig.toml". In this case, the flags available for Geth/Erigon, Prysm/Lighthouse are available. Here's an example:

```bash
# Geth Configs
$ lukso start --geth-config "./myconfig.toml"
# Prysm Configs
$ lukso start --prysm-config "./myconfig.yaml" --geth-bootnodes "mycustombootnode0000"

# An experienced user might also want to start custom clients
$ lukso start --lighthouse --erigon
```

#### How to start your validator (keys & tx fee recipient)

Starts your node as a validator node

```bash
$ lukso start --validator
```

The transaction fee recipient; aka coinbase, is the address where the transactions fees are sent to. To start your validator, you will pass this transaction fee recipient address and LUKSO's CLI will perform a few checks regarding your keys.

```bash
$ lukso start --validator --transaction-fee-recipient  "0x12345678..."
```

```bash
# validator keys and password
$ lukso start --validator --validator-keys "./mainnet-keystore" --validator-password "./myfile.txt"
```

#### How to set & customize a log folder

You can setup a custom log directory when starting LUKSO client by indicating the location of your folder.

```bash
# Path to geth log file that you want to log
$ lukso start --log-folder "./myCustomLogFolder"
```

Log files are by default un-capped. Be aware that the # of these files can grow very large. You can cap the file size in n MB.
Example: `lukso start --log-size 3` to cap the file to 3MB of data

```bash
$ lukso start --log-size n
```

| Flag                                                                                  | Description                                                                                                                 |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet                                                                             | starts LUKSO's mainnet. User is connected to mainnet by default.                                                            |
| --testnet                                                                             | starts LUKSO's testnet                                                                                                      |
| --devnet                                                                              | starts LUKSO's devnet                                                                                                       |
| --geth-\*                                                                             | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon-\*                                                                           | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm-\*                                                                            | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse-\*                                                                       | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --geth-config                                                                         | path to "./myconfig.toml" file                                                                                              |
| --prysm-config "./myconfig.yaml" --geth-bootnodes "mycustombootnode00000"             | path to "./myconfig.yaml" file & custom geth boot nodes                                                                     |
| --validator --transaction-fee-recipient                                               | address that receives block fees (0x12345..abcd).                                                                           |
| --validator --validator-keys "./mainnet-keystore" --validator-password "./myfile.txt" | passes the validator keys and password from a custom directory                                                              |
| --validator -std-output                                                               | set output to console                                                                                                       |
| --log -folder "./myCustomLogFolder"                                                   | user can setup a custom log directory when starting LUKSO client                                                            |
| --log -size n                                                                         | Log files capped to the size in MB                                                                                          |

## How to check the status of LUKSO node

Displays the most recent status of LUKSO's node

```bash
$ lukso status
```

## How to stop LUKSO node

Stops all client's activities. usually used when upgrading the client or running maintenance tasks.

```bash
# This command stops all your client's activities
$ lukso stop
# only stops the validator client
$ lukso stop --validator
# only stops the execution client
$ lukso stop --execution
# only stops the consensus client
$ lukso stop --consensus
```

## How to update LUKSO

Updates LUKSO client to the latest available version

```bash
$ lukso update
```

| Flag            | Description                                                                                                                 |
| --------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet       | updates LUKSO's mainnet                                                                                                     |
| --testnet       | updates LUKSO's testnet                                                                                                     |
| --devnet        | updates LUKSO's devnet                                                                                                      |
| --geth-\*       | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon-\*     | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm-\*      | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse-\* | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
| --validator     | updates a specific version of the validator                                                                                 |

#### How to update a specific client to a specific version

```bash
# updates to the specific version of (geth/prysm/erigon/lighthouse) client - Example Geth v1.11.4
$ lukso update --geth-tag

#  updates a specific version of the validator
$ lukso update --validator-tag
```

## Running your validator

The main activity you can perform as a validator is depositing your keys.

## How Genesis Validators proceed with their deposits

As a validator you can deposit your keys using an RPC connection of your choice.

```bash
$ lukso validator deposit --deposit-data-json "./validator-deposit-data.json" [--gasPrice "1000000000" --rpc "https://infura.io./apiv1"]
```

```bash
$ lukso validator deposit --genesis --deposit-data-json "./validator-deposit-data.json" --rpc "https://infura.io./apiv1" [--gas-price "1000000000" --start-from-index N]
```

As a Genesis Validator you can provide an indicative voting for the prefered initial token supply of LYX, which will determine how much the Foundation will receive. See the https://deposit.mainnet.lukso.network website for details.
You can choose between:
1: 35M LYX
2: 42M LYX (This option is the prefered one by the Foundation)
3: 100M LYX
4: No vote

The private key you are using requires enough the LYXe to fund pay for validator keys deposits.
Make sure your privatekey has sufficient balances:
320 LYXe
1.6 ETH for GAS
Enter private key:

> \_

| Flag                             | Description                                                                                                                 |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| --mainnet                        | runs a validator on LUKSO's mainnet (default)                                                                               |
| --testnet                        | runs a validator on LUKSO's testnet                                                                                         |
| --devnet                         | runs a validator on LUKSO's devnet                                                                                          |
| --geth-\*                        | \* Pass any flag to the Geth node [See docs for details](https://geth.ethereum.org/docs/fundamentals/command-line-options)  |
| --erigon-\*                      | \* Pass any flag to the Erigon node [See docs for details](https://github.com/ledgerwatch/erigon)                           |
| --prysm-\*                       | \* Pass any flag to the Prysm node [See docs for details](https://docs.prylabs.network/docs/prysm-usage/parameters)         |
| --lighthouse-\*                  | \* Pass any flag to the Lighthouse node [See docs for details](https://lighthouse-book.sigmaprime.io/advanced-datadir.html) |
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
