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

# build
go build -o lukso

# run
./lukso <command>
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
