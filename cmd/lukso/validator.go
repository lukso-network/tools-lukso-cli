package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"

	"github.com/m8b-dev/lukso-cli/contracts/bindings"
)

const (
	gasLimit                      = 21_000
	depositContractAddress        = "0x000000000000000000000000000000000000cafe"
	genesisDepositContractAddress = "0x75D1f4695Eb87d60eD4EAE2c0CF05e7428Fa4b5F"
	lyxeContractAddress           = "0x7A2AC110202ebFdBB5dB15Ea994ba6bFbFcFc215"

	maxTxsPerBlock     = 10
	blockFetchInterval = 3 // in seconds
)

type DepositDataKey struct {
	PubKey                string   `json:"pubkey"`
	WithdrawalCredentials string   `json:"withdrawal_credentials"`
	Amount                *big.Int `json:"amount"`
	Signature             string   `json:"signature"`
	DepositMessageRoot    string   `json:"deposit_message_root"`
	DepositDataRoot       string   `json:"deposit_data_root"`
	ForkVersion           string   `json:"fork_version"`
	NetworkName           string   `json:"network_name"`
	DepositCliVersion     string   `json:"deposit_cli_version"`
}

func sendDeposit(ctx *cli.Context) error {
	log.Info("Dialing up blockchain for gas info...")
	eth, err := ethclient.Dial(ctx.String(rpcFlag))
	if err != nil {
		return err
	}

	c := context.Background()
	gasPrice, err := eth.SuggestGasPrice(c)
	if err != nil {
		return err
	}

	log.Infof("Gas Price fetched: %v", gasPrice)

	var (
		selectedDeposit string
		supplyAmount    int
	)

	depositPath := ctx.String(depositFlag)
	genesisDepositPath := ctx.String(genesisDepositFlag)

	if depositPath != "" && genesisDepositPath != "" {
		return errTooManyDepositsProvided
	}

	switch {
	case depositPath != "":
		selectedDeposit = depositPath
	case genesisDepositPath != "":
		selectedDeposit = genesisDepositPath
		supplyAmount, err = processTokenOption()
		if err != nil {
			return err
		}

	default:
		return errDepositNotProvided
	}

	depositKeys, err := parseDepositFile(selectedDeposit)
	if err != nil {
		return err
	}

	keysNum := len(depositKeys)
	singleTxGasPrice := big.NewInt(0).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	overallGasPrice := big.NewInt(0).Mul(singleTxGasPrice, big.NewInt(int64(keysNum)))
	overallGasPriceInt := overallGasPrice.Int64()
	overallGasPriceEth, _ := big.NewRat(overallGasPriceInt, 1_000_000_000_000_000_000).Float64()

	log.Infof("Before proceeding make sure that your private key has sufficient balance:\n"+
		"- %v ETH\n"+
		"- %v LYXe\n\n", overallGasPriceEth, keysNum*32)

	message := "Please enter your private key: \n> "
	input := registerInputWithMessage(message)

	privKey, err := crypto.HexToECDSA(input)
	if err != nil {
		return err
	}

	senderAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	chainId, err := eth.ChainID(c)
	if err != nil {
		return err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privKey, chainId)
	if err != nil {
		return err
	}

	lyxMock, err := bindings.NewLYXe(common.HexToAddress(lyxeContractAddress), eth)
	if err != nil {
		return err
	}

	ethDeposit, err := bindings.NewEthereumDeposit(common.HexToAddress(depositContractAddress), eth)
	if err != nil {
		return err
	}

	startingBlock, err := eth.BlockNumber(c)
	if err != nil {
		return err
	}

	// we take nonce once for 1st transaction and increment it manually
	nonce, err := eth.PendingNonceAt(c, senderAddr)
	if err != nil {
		return err
	}

	txCount := 1 // if txCount reaches 10 then we have to wait for another block
	txHashes := make([]common.Hash, 0)

	for i, key := range depositKeys {
		currentBlock, err := eth.BlockNumber(c)
		if err != nil {
			return err
		}

		if currentBlock != startingBlock {
			startingBlock = currentBlock
			txCount = 1
		}

		if txCount > maxTxsPerBlock {
			fmt.Println("Reached 10 tx per block - waiting for next block...")
			startingBlock, err = waitForNextBlock(c, eth, currentBlock)
			if err != nil {
				return err
			}

			txCount = 1
		}

		fmt.Printf("Deposit %d/%d\n", i+1, keysNum)
		fmt.Println("Amount:", key.Amount.String())
		fmt.Println("Public Key:", key.PubKey)
		fmt.Println("Withdraw credentials:", key.WithdrawalCredentials)
		fmt.Println("Fork version:", key.ForkVersion)
		fmt.Println("Deposit data root:", key.DepositDataRoot)
		fmt.Println("Signature:", key.Signature, "\n")

		opts.Nonce = big.NewInt(int64(nonce))
		opts.From = senderAddr

		var tx *types.Transaction
		switch selectedDeposit {
		case genesisDepositPath:
			depositData, err := prepareDepositData(key, supplyAmount)
			if err != nil {
				log.Error("Couldn't send transaction - deposit data provided is invalid  - skipping...")
			}

			tx, err = lyxMock.Send(
				opts,
				common.HexToAddress(genesisDepositContractAddress),
				big.NewInt(0).Mul(big.NewInt(32), big.NewInt(1000000000000000000)),
				depositData,
			)

			if err != nil {
				return err
			}

		case depositPath:
			var depositDataRoot [32]byte

			depositDataRootBytes := []byte(key.DepositDataRoot)
			if len(depositDataRootBytes) != 32 {
				log.Error("Couldn't send transaction - deposit data root is not 32 bytes long - skipping...")

				continue
			}

			for depI := range depositDataRoot {
				depositDataRoot[depI] = depositDataRootBytes[depI]
			}

			tx, err = ethDeposit.Deposit(
				opts,
				[]byte(key.PubKey),
				[]byte(key.WithdrawalCredentials),
				[]byte(key.Signature),
				depositDataRoot,
			)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Transaction %d/%d successful! Transaction hash: %v\n\n", i+1, keysNum, tx.Hash().String())

		nonce = tx.Nonce() + 1 // we could do nonce += 1, but it's just to make sure we are +1 ahead of previous tx
		txCount++
		txHashes = append(txHashes, tx.Hash())
	}

	return nil
}

func initValidator(ctx *cli.Context) error {
	if ctx.String(validatorKeysDirFlag) == "" {
		return errKeysNotProvided
	}
	args := []string{
		"accounts",
		"import",
		"--keys-dir", ctx.String(validatorWalletDirFlag),
		"--wallet-dir", ctx.String(validatorWalletDirFlag),
	}

	if ctx.String(validatorWalletPasswordFileFlag) != "" {
		args = append(args, "--wallet-password-file", ctx.String(validatorWalletPasswordFileFlag))
	}

	initCommand := exec.Command("validator", args...)

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr
	initCommand.Stdin = os.Stdin

	err := initCommand.Run()
	if err != nil {
		return err
	}

	return nil
}

func parseDepositFile(depositFilePath string) (keys []DepositDataKey, err error) {
	f, err := os.ReadFile(depositFilePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &keys)

	return
}

func prepareDepositData(key DepositDataKey, amount int) (depositData []byte, err error) {
	bytePubKey, err := hex.DecodeString(key.PubKey)
	if err != nil {
		return
	}
	if len(bytePubKey) != 48 {
		return
	}

	byteWithdrawalCredentials, err := hex.DecodeString(key.WithdrawalCredentials)
	if err != nil {
		return
	}
	if len(byteWithdrawalCredentials) != 32 {
		return
	}

	byteSignature, err := hex.DecodeString(key.Signature)
	if err != nil {
		return
	}
	if len(byteSignature) != 96 {
		return
	}

	byteDepositDataRoot, err := hex.DecodeString(key.DepositDataRoot)
	if err != nil {
		return
	}
	if len(byteDepositDataRoot) != 32 {
		return
	}

	depositData = append(depositData, bytePubKey...)
	depositData = append(depositData, byteWithdrawalCredentials...)
	depositData = append(depositData, byteSignature...)
	depositData = append(depositData, byteDepositDataRoot...)
	depositData = append(depositData, byte(amount))
	log.Error("AMOUNT: ", amount)

	return
}

// waitForNextBlock fetches current block in 3 second intervals, and
func waitForNextBlock(c context.Context, eth *ethclient.Client, currentBlock uint64) (blockNumber uint64, err error) {
	for {
		time.Sleep(time.Second * blockFetchInterval)
		blockNumber, err = eth.BlockNumber(c)
		if err != nil {
			return
		}

		if currentBlock != blockNumber {
			break
		}
	}

	return
}

func processTokenOption() (amount int, err error) {
	message := `As a Genesis Validator you can provide an indicative voting for the preferred initial token supply of LYX, which will determine how much the Foundation will receive. See the https://deposit.mainnet.lukso.network website for details.
You can choose between:
1: 35M LYX
2: 42M LYX (This option is the prefered one by the Foundation)
3: 100M LYX
4: An arbitrary amount from 0-100
5: No vote
Please enter your choice (1-5):
> `
	var option int
	for option < 1 || option > 5 {
		input := registerInputWithMessage(message)
		option, err = strconv.Atoi(input)
		if err != nil {
			log.Warn("Please provide a valid option")

			continue
		}
		if option < 1 || option > 5 {
			log.Warn("Please provide an option between 1-5")
		}
	}

	// we only get here when user provides valid option, so no need for catching weird options
	switch option {
	case 1:
		amount = 35
	case 2:
		amount = 42
	case 3:
		amount = 100
	case 4:
		option = -1
		for option < 0 || option > 100 {
			input := registerInputWithMessage("Please enter initial token supply: \n> ")
			option, err = strconv.Atoi(input)
			if err != nil {
				log.Warn("Please provide a valid option")

				continue
			}
			if option < 0 || option > 100 {
				log.Warn("Please provide an option between 0-100")
			}
		}

		amount = option
	case 5:
		amount = 0
	}

	return
}
