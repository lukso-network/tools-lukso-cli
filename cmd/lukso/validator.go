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
	"github.com/m8b-dev/lukso-cli/contracts/bindings"
	"github.com/urfave/cli/v2"
)

const (
	ether                         = 1_000_000_000_000_000_000
	gwei                          = 1_000_000_000
	gasMargin                     = 500_000
	gasBump                       = 50_000
	depositContractAddress        = "0x000000000000000000000000000000000000cafe"
	genesisDepositContractAddress = "0x9C2Ae5bC047Ca794d9388aB7A2Bf37778f9aBA73"
	lyxeContractAddress           = "0x790c4379C82582F569899b3Ca71E78f19AeF82a5"
	defaultRpc                    = "https://rpc.execution.3030.devnet.lukso.dev/"

	errUnderpriced = "transaction underpriced" //nolint:all // catches both replacement and normal underpriced

	blockFetchInterval      = 12 // in seconds
	amountOfLyxPerValidator = 32 // LYXe for Genesis, LYX for non genesis
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

type depositController struct {
	c context.Context

	eth            *ethclient.Client
	genesisDeposit *bindings.LYXe
	deposit        *bindings.EthereumDeposit
	senderAddr     common.Address

	depositKeys   []DepositDataKey
	startingIndex int
	keysNum       int
	txOpts        *bind.TransactOpts
}

func newDepositController(rpc string, depositKeys []DepositDataKey, startingIndex int) (dc depositController, err error) {
	c := context.Background()
	keysLen := len(depositKeys)
	if keysLen < 1 {
		err = errDepositNotProvided

		return
	}

	if rpc == "" {
		log.Infof("RPC not provided - falling back to default RPC: %s", defaultRpc)
	}
	log.Info("Dialing up blockchain...")
	eth, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}

	genDep, err := bindings.NewLYXe(common.HexToAddress(lyxeContractAddress), eth)
	if err != nil {
		return
	}

	dep, err := bindings.NewEthereumDeposit(common.HexToAddress(depositContractAddress), eth)
	if err != nil {
		return
	}

	if startingIndex < 0 {
		log.Error("Couldn't send deposits: starting index is smaller than 0")
		err = errIndexOutOfBounds

		return
	}
	if startingIndex >= keysLen {
		log.Error("Couldn't send deposits: starting index is greater than number of deposits")
		err = errIndexOutOfBounds

		return
	}

	depositKeys = depositKeys[startingIndex:]

	message := "Please enter your private key: \n> "

	// TODO: input should be password input
	input := strings.TrimPrefix(registerInputWithMessage(message), "0x")

	privKey, err := crypto.HexToECDSA(input)
	if err != nil {
		return
	}

	senderAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	chainId, err := eth.ChainID(c)
	if err != nil {
		return
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, chainId)
	if err != nil {
		return
	}

	txOpts.From = senderAddr // will stay global, no matter what deposit we are making

	dc = depositController{
		c,
		eth,
		genDep,
		dep,
		senderAddr,
		depositKeys,
		startingIndex,
		len(depositKeys),
		txOpts,
	}

	return
}

// estimateGas estimates gas for sending all deposits, displays this information to user and waits for his confirmation
func (dc depositController) estimateGas(isGenesisDeposit bool) (accepted bool, err error) {
	var (
		message     string
		tx          *types.Transaction
		depositData []byte
	)

	nonce, err := dc.eth.PendingNonceAt(dc.c, dc.senderAddr)
	if err != nil {
		return
	}

	dc.txOpts.NoSend = true
	dc.txOpts.Nonce = big.NewInt(int64(nonce))

	switch isGenesisDeposit {
	case true:
		depositData, err = encodeGenesisDepositDataKey(dc.depositKeys[0], 0)

		if err != nil {
			return
		}

		tx, err = dc.genesisDeposit.Send(
			dc.txOpts,
			common.HexToAddress(genesisDepositContractAddress),
			big.NewInt(0).Mul(big.NewInt(amountOfLyxPerValidator), big.NewInt(ether)),
			depositData,
		)

		gasReadable := estimateGas(tx, int64(dc.keysNum))

		message = fmt.Sprintf("Before proceeding make sure that your private key has sufficient balance:\n"+
			"- %v ETH\n"+
			"- %v LYXe (%v * %v validator[s])\nDo you wish to continue? [Y/n]: ", gasReadable, dc.keysNum*amountOfLyxPerValidator, amountOfLyxPerValidator, dc.keysNum)

	case false:
		dc.txOpts.Value = big.NewInt(0).Mul(big.NewInt(amountOfLyxPerValidator), big.NewInt(ether))

		var depositDataRoot [32]byte

		depositData, err = encodeGenesisDepositDataKey(dc.depositKeys[0], 0)
		if err != nil {
			return
		}

		startI := 176
		for i := 0; i < 32; i++ {
			depositDataRoot[i] = depositData[startI+i]
		}

		tx, err = dc.deposit.Deposit(
			dc.txOpts,
			depositData[:48],
			depositData[48:80],
			depositData[80:176],
			depositDataRoot,
		)
		if err != nil {
			return
		}

		gasReadable := estimateGas(tx, int64(dc.keysNum))

		message = fmt.Sprintf("Before proceeding make sure that your private key has sufficient balance:\n"+
			"- %v LYX\nDo you wish to continue? [Y/n]: ", float64(dc.keysNum*amountOfLyxPerValidator)+gasReadable)

	}

	accepted = true
	input := registerInputWithMessage(message)
	if !strings.EqualFold(input, "y") && input != "" {
		log.Info("Aborting...")

		accepted = false
	}

	return
}

func (dc depositController) sendDeposits(isGenesisDeposit bool, maxTxsPerBatch int) (err error) {
	var (
		txSentCount  = 0 // if txSentCount reaches 10 then we have to wait for another block
		nonce        uint64
		supplyAmount int
		currentBatch = 0
	)

	txsSent := make([]*types.Transaction, 0)

	if isGenesisDeposit {
		supplyAmount, err = chooseSupply()
		if err != nil {
			return
		}
	}

	nonce, err = dc.eth.PendingNonceAt(dc.c, dc.senderAddr)
	if err != nil {
		return
	}

	dc.txOpts.NoSend = false

	for i, key := range dc.depositKeys {
		if txSentCount == maxTxsPerBatch {
			fmt.Printf("Reached %d txs sent - waiting for receipts...\n", maxTxsPerBatch)
			failedBatchedTxIndex, err := dc.waitForReceipts(txsSent)
			if err != nil {
				failedDepositIndex := dc.startingIndex + currentBatch*maxTxsPerBatch + failedBatchedTxIndex
				log.Errorf("Sent transaction has failed with error: %v - aborting...", err)
				log.Errorf("To continue with your deposits please run a deposit command once again, "+
					"but with --start-from-index flag to continue from failed transaction, example:\n"+
					"lukso validator deposit --deposit-data-json *your deposit data file* --start-from-index %d",
					failedDepositIndex,
				)

				return err
			}

			txSentCount = 0
			txsSent = make([]*types.Transaction, 0)
			currentBatch++
		}

		fmt.Printf("Deposit %d/%d\n", i+1, dc.keysNum)
		fmt.Println("Amount:", key.Amount.String())
		fmt.Println("Public Key:", key.PubKey)
		fmt.Println("Withdraw credentials:", key.WithdrawalCredentials)
		fmt.Println("Fork version:", key.ForkVersion)
		fmt.Println("Deposit data root:", key.DepositDataRoot)
		fmt.Println("Signature:", key.Signature)
		fmt.Println("")

		dc.txOpts.Nonce = big.NewInt(int64(nonce))
		dc.txOpts.Value = big.NewInt(0)

		var tx *types.Transaction
		switch isGenesisDeposit {
		case true:
			var depositData []byte

			depositData, err = encodeGenesisDepositDataKey(key, supplyAmount)
			if err != nil {
				log.Error("Couldn't send transaction - deposit data provided is invalid  - skipping...")
			}

			tx, err = dc.genesisDeposit.Send(
				dc.txOpts,
				common.HexToAddress(genesisDepositContractAddress),
				big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether)),
				depositData,
			)

		case false:
			dc.txOpts.Value = big.NewInt(0).Mul(big.NewInt(32), big.NewInt(ether))
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

			tx, err = dc.deposit.Deposit(
				dc.txOpts,
				depositData[:48],
				depositData[48:80],
				depositData[80:176],
				depositDataRoot,
			)
		}

		if err != nil {
			failedDepositIndex := dc.startingIndex + currentBatch*maxTxsPerBatch + txSentCount
			log.Errorf("Sent transaction has failed with error: %v - aborting...", err)
			log.Errorf("To continue with your deposits please run a deposit command once again, "+
				"but with --start-from-index flag to continue from failed transaction, example:\n"+
				"lukso validator deposit --deposit-data-json *your deposit data file* --start-from-index %d",
				failedDepositIndex,
			)

			return err
		}

		txsSent = append(txsSent, tx)

		fmt.Printf("Transaction %d/%d sent! Transaction hash: %v\n\n", i+1, dc.keysNum, tx.Hash().String())

		nonce = tx.Nonce() + 1 // we could do nonce += 1, but it's just to make sure we are +1 ahead of previous tx
		txSentCount++
	}

	return nil
}

func sendDeposit(ctx *cli.Context) (err error) {
	depositPath := ctx.String(depositDataJson)
	isGenesisDeposit := ctx.Bool(genesisDepositFlag)

	depositKeys, err := parseDepositDataFile(depositPath)
	if err != nil {
		return err
	}

	dc, err := newDepositController(ctx.String(rpcFlag), depositKeys, ctx.Int(startFromIndexFlag))
	if err != nil {
		return nil // to avoid duplicating error messages
	}

	accepted, err := dc.estimateGas(isGenesisDeposit)
	if err != nil {
		return
	}
	if !accepted {
		return nil
	}

	err = dc.sendDeposits(isGenesisDeposit, ctx.Int(maxTxsPerBatchFlag))

	return err
}

func importValidator(ctx *cli.Context) error {
	if len(os.Args) < 3 {
		return errNotEnoughArguments
	}

	args := []string{
		"accounts",
		"import",
	}

	// we don't want to pass those flags
	mainnet := fmt.Sprintf("--%s", mainnetFlag)
	testnet := fmt.Sprintf("--%s", testnetFlag)
	devnet := fmt.Sprintf("--%s", devnetFlag)
	walletDir := fmt.Sprintf("--wallet-dir")

	for _, osArg := range os.Args[3:] {
		if osArg == mainnet || osArg == testnet || osArg == devnet {
			continue
		}

		args = append(args, osArg)
	}

	isWalletProvided := false
	walletDefault := ctx.String(validatorWalletDirFlag)

	for _, arg := range args {
		if arg == walletDir {
			isWalletProvided = true
		}
	}

	if !isWalletProvided {
		args = append(args, walletDir, walletDefault)
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

// waitForReceipts waits until sent transactions
func (dc depositController) waitForReceipts(txs []*types.Transaction) (failedIndex int, err error) {
	validatedTxs := make([]bool, len(txs))
	for {
		log.Infof("Waiting %d seconds before fetching receipts...", blockFetchInterval)
		time.Sleep(time.Second * blockFetchInterval)
		for i, tx := range txs {
			var (
				isPending = false //nolint:all
				receipt   *types.Receipt
			)

			if validatedTxs[i] {
				continue
			}

			_, isPending, err = dc.eth.TransactionByHash(dc.c, tx.Hash())
			if err != nil {
				return
			}
			if isPending {
				log.Infof("tx with hash %s is still pending - continuing", tx.Hash().String())
				continue
			}

			log.Infof("getting receipt for tx with hash %s", tx.Hash().String())
			receipt, err = dc.eth.TransactionReceipt(dc.c, tx.Hash())
			if err != nil {
				return
			}

			log.Infof("Got receipt for tx with hash %s, status: %d", tx.Hash().String(), receipt.Status)
			if receipt.Status == 0 {
				err = errTransactionFailed
				failedIndex = i

				return
			}

			validatedTxs[i] = true
		}

		// check if all txs are validated
		allValidated := true
		for _, validated := range validatedTxs {
			if !validated {
				allValidated = false
			}
		}
		if allValidated {
			break
		}
	}

	return
}

func chooseSupply() (amount int, err error) {
	message := `As a Genesis Validator you can provide an indicative voting for the preferred initial token supply of LYX,
which will determine how much the Foundation will receive.
See the https://deposit.mainnet.lukso.network website for details.
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
