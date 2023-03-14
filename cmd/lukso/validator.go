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
	"strings"
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
	ether                         = 1_000_000_000_000_000_000
	depositContractAddress        = "0x000000000000000000000000000000000000cafe"
	genesisDepositContractAddress = "0x9C2Ae5bC047Ca794d9388aB7A2Bf37778f9aBA73"
	lyxeContractAddress           = "0x790c4379C82582F569899b3Ca71E78f19AeF82a5"

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
	// Dialing clients and creating bindings
	log.Info("Dialing up blockchain for gas info...")
	eth, err := ethclient.Dial(ctx.String(rpcFlag))
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

	c := context.Background()

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

	depositKeys, err := parseDepositDataFile(selectedDeposit)
	if err != nil {
		return err
	}

	keysNum := len(depositKeys)

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

	startingBlock, err := eth.BlockNumber(c)
	if err != nil {
		return err
	}

	// we take nonce once for 1st transaction and increment it manually
	nonce, err := eth.PendingNonceAt(c, senderAddr)
	if err != nil {
		return err
	}

	// estimate gas based on 1st deposit and inform user before proceeding
	if keysNum < 1 {
		return errDepositNotProvided
	}

	opts.Nonce = big.NewInt(int64(nonce))
	opts.From = senderAddr
	opts.NoSend = true

	key := depositKeys[0]

	var tx *types.Transaction

	switch selectedDeposit {
	case genesisDepositPath:
		depositData, err := parseDepositDataKey(key, supplyAmount)

		if err != nil {
			return err
		}

		tx, err = lyxMock.Send(
			opts,
			common.HexToAddress(genesisDepositContractAddress),
			big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether)),
			depositData,
		)

		gasReadable := estimateGas(tx, int64(keysNum))

		message = fmt.Sprintf("Before proceeding make sure that your private key has sufficient balance:\n"+
			"- %v ETH\n"+
			"- %v LYXe\nDo you wish to continue? [Y/n]: ", gasReadable, keysNum*32)

	case depositPath:
		opts.Value = big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether))
		var depositDataRoot [32]byte

		depositData, err := parseDepositDataKey(key, 0)
		if err != nil {
			return err
		}

		startI := 176
		for i := 0; i < 32; i++ {
			depositDataRoot[i] = depositData[startI+i]
		}

		tx, err = ethDeposit.Deposit(
			opts,
			depositData[:48],
			depositData[48:80],
			depositData[80:176],
			depositDataRoot,
		)
		if err != nil {
			return err
		}

		gasReadable := estimateGas(tx, int64(keysNum))

		message = fmt.Sprintf("Before proceeding make sure that your private key has sufficient balance:\n"+
			"- %v LYX\nDo you wish to continue? [Y/n]: ", float64(keysNum*32)+gasReadable)

	}

	input = registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") {
		log.Info("Aborting...")

		return nil
	}

	txCount := 1 // if txCount reaches 10 then we have to wait for another block

	for i, key := range depositKeys {
		currentBlock, err := eth.BlockNumber(c)
		if err != nil {
			return err
		}

		if currentBlock != startingBlock {
			startingBlock = currentBlock
			txCount = 1
		}

		if txCount > ctx.Int(maxTxsPerBlock) {
			fmt.Printf("Reached %d tx per block - waiting for next block...", ctx.Int(maxTxsPerBlock))
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
		opts.NoSend = false
		opts.Value = big.NewInt(0)

		var tx *types.Transaction
		switch selectedDeposit {
		case genesisDepositPath:
			depositData, err := parseDepositDataKey(key, supplyAmount)
			if err != nil {
				log.Error("Couldn't send transaction - deposit data provided is invalid  - skipping...")
			}

			tx, err = lyxMock.Send(
				opts,
				common.HexToAddress(genesisDepositContractAddress),
				big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether)),
				depositData,
			)

			if err != nil {
				return err
			}

		case depositPath:
			opts.Value = big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether))
			var depositDataRoot [32]byte

			depositData, err := parseDepositDataKey(key, 0)
			if err != nil {
				return err
			}

			startI := 176
			for i := 0; i < 32; i++ {
				depositDataRoot[i] = depositData[startI+i]
			}

			tx, err = ethDeposit.Deposit(
				opts,
				depositData[:48],
				depositData[48:80],
				depositData[80:176],
				depositDataRoot,
			)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Transaction %d/%d sent! Transaction hash: %v\n\n", i+1, keysNum, tx.Hash().String())

		nonce = tx.Nonce() + 1 // we could do nonce += 1, but it's just to make sure we are +1 ahead of previous tx
		txCount++
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

func parseDepositDataFile(depositFilePath string) (keys []DepositDataKey, err error) {
	f, err := os.ReadFile(depositFilePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &keys)

	return
}

func parseDepositDataKey(key DepositDataKey, amount int) (depositData []byte, err error) {
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
4: No vote
Please enter your choice (1-4):
> `
	var option int
	for option < 1 || option > 4 {
		input := registerInputWithMessage(message)
		option, err = strconv.Atoi(input)
		if err != nil {
			log.Warn("Please provide a valid option")

			continue
		}
		if option < 1 || option > 4 {
			log.Warn("Please provide an option between 1-4")
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
		amount = 0
	}

	return
}

// estimateGas estimates a human-readable gas price for txCount transactions summed.
func estimateGas(tx *types.Transaction, txCount int64) float64 {
	txGas := big.NewInt(int64(tx.Gas()))
	txGasFeeCap := tx.GasFeeCap()
	txGas = txGas.Mul(txGas, txGasFeeCap)
	allTxGas := big.NewInt(0).Mul(txGas, big.NewInt(txCount))
	gasInEther, _ := big.NewRat(allTxGas.Int64(), ether).Float64()

	return gasInEther
}
