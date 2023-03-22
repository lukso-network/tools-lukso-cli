package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	gasMargin                     = 500_000
	gasBump                       = 50_000
	depositContractAddress        = "0x000000000000000000000000000000000000cafe"
	genesisDepositContractAddress = "0x9C2Ae5bC047Ca794d9388aB7A2Bf37778f9aBA73"
	lyxeContractAddress           = "0x790c4379C82582F569899b3Ca71E78f19AeF82a5"

	errUnderpriced = "transaction underpriced" // catches both replacement and normal underpriced

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

	var supplyAmount int

	depositPath := ctx.String(depositDataJson)
	isGenesisDeposit := ctx.Bool(genesisDepositFlag)

	if isGenesisDeposit {
		supplyAmount, err = chooseSupply()
		if err != nil {
			return err
		}
	}

	depositKeys, err := parseDepositDataFile(depositPath)
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

	switch isGenesisDeposit {
	case true:
		depositData, err := encodeGenesisDepositDataKey(key, supplyAmount)

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

	case false:
		opts.Value = big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether))
		var depositDataRoot [32]byte

		depositData, err := encodeGenesisDepositDataKey(key, 0)
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
		var currentBlock uint64

		currentBlock, err = eth.BlockNumber(c)
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
		switch isGenesisDeposit {
		case true:
			var depositData []byte

			depositData, err = encodeGenesisDepositDataKey(key, supplyAmount)
			if err != nil {
				log.Error("Couldn't send transaction - deposit data provided is invalid  - skipping...")
			}

			tx, err = lyxMock.Send(
				opts,
				common.HexToAddress(genesisDepositContractAddress),
				big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether)),
				depositData,
			)

		case false:
			opts.Value = big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether))
			var (
				depositData     []byte
				depositDataRoot [32]byte
			)

			depositData, err = encodeGenesisDepositDataKey(key, 0)
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
		}

		if err != nil && strings.Contains(err.Error(), errUnderpriced) {
			tx, err = bumpGasAndSend(eth, tx, opts)
			if err != nil {
				log.Fatalf("Couldn't bump transation gas: %v", err)

				return err
			}
		}

		fmt.Printf("Transaction %d/%d sent! Transaction hash: %v\n\n", i+1, keysNum, tx.Hash().String())

		nonce = tx.Nonce() + 1 // we could do nonce += 1, but it's just to make sure we are +1 ahead of previous tx
		txCount++
	}

	return nil
}

func bumpGasAndSend(eth *ethclient.Client, tx *types.Transaction, signer *bind.TransactOpts) (*types.Transaction, error) {
	log.Warn("Transaction failed with underpriced error - resending...")
	gas := tx.Gas() + uint64(gasBump)

	log.Debugf("Bumping gas, had %d, added %d", tx.Gas(), gas)
	gas += tx.Gas()

	fee := tx.GasFeeCap()
	feeFloat, ok := big.NewFloat(0).SetString(fee.String())
	if ok {
		feeFloat = feeFloat.Mul(feeFloat, big.NewFloat(1.15))
		feeFloat.Int(fee)
		log.Debugf("Bumping fee %s -> %s WEI", tx.GasFeeCap().String(), fee.String())
	} else {
		log.Warnf("Failed to bump base fee: not ok fir bigFloat construction")
	}

	timeout := time.Now().Add(time.Second * 10)
	nonce := tx.Nonce()
	for {
		if time.Now().After(timeout) {
			return nil, errors.New("failed to send transaction - timed out with handleable noncing errors")
		}

		signed, err := signer.Signer(signer.From, types.NewTx(&types.DynamicFeeTx{
			ChainID:   tx.ChainId(),
			Nonce:     nonce,
			GasTipCap: tx.GasTipCap(),
			GasFeeCap: fee,
			Gas:       gas,
			To:        tx.To(),
			Value:     tx.Value(),
			Data:      tx.Data(),
		}))

		if err != nil {
			return nil, err
		}

		log.Debugf("Signed bumped gas tx, gas=%d hash=%s", signed.Gas(), signed.Hash().String())

		err = eth.SendTransaction(context.Background(), signed)
		if err != nil {
			switch strings.Contains(err.Error(), errUnderpriced) {
			case true: // ie. data race for nonce getter
				nonce++
				continue
			default:
				return nil, err
			}
		}

		return signed, nil
	}
}

func initValidator(ctx *cli.Context) error {
	args := []string{
		"accounts",
		"import",
		"--keys-dir", ctx.String(validatorKeysFlag),
		"--wallet-dir", ctx.String(validatorKeysFlag),
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

func encodeGenesisDepositDataKey(key DepositDataKey, amount int) (depositData []byte, err error) {
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

func chooseSupply() (amount int, err error) {
	message := `As a Genesis Validator you can provide an indicative voting for the preferred initial token supply of LYX, which will determine how much the Foundation will receive. See the https://deposit.mainnet.lukso.network website for details.
You can choose between:
1: 35M LYX
2: 42M LYX (This option is the preferred one by the Foundation)
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
	allTxGas.Add(allTxGas, big.NewInt(gasMargin))
	gasInEther, _ := big.NewRat(allTxGas.Int64(), ether).Float64()

	return gasInEther
}
