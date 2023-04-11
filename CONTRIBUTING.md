# Contributing

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
│   └───macos_packages    // macOS Codesigning Scripts
|
└───pid                   // Process ID Management
```

## Development & Testing

### Install CLI from Codebase

#### Clone the LUKSO CLI Repository

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

### Smart Contracts Bindings

#### Prerequisites:

- Solidity Compiler [solc.js](https://github.com/ethereum/solidity)
- ABI Generator [abigen](https://geth.ethereum.org/docs/tools/abigen)

#### Generate Bindings

The contract can be found in: [lukso-network/network-genesis-deposit-contract](https://github.com/lukso-network/network-genesis-deposit-contract).

1. Paste the interacting smart contract into the [`contracts`](./contracts) directory
2. Generate the ABIs from your smart contracts

Note: the files under: `contracts/*.sol` are ignored from the git repo.

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
