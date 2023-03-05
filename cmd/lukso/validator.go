package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	depositMock "github.com/m8b-dev/lukso-cli/contracts/bindings"
	"math/big"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const (
	gasLimit               = 21_000
	depositContractAddress = "0xcd2a3d9f938e13cd947ec0i8um67fe734df8d8861"
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
	eth, err := ethclient.Dial("https://rpc.2022.l16.lukso.network")
	if err != nil {
		return err
	}

	c := context.Background()
	gasPrice, err := eth.SuggestGasPrice(c)
	if err != nil {
		return err
	}
	log.Infof("Gas Price fetched: %v WEI (~%v GWEI)", gasPrice, big.NewInt(0).Div(gasPrice, big.NewInt(1_000_000_000)))

	var selectedDeposit string

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
	default:
		return errDepositNotProvided
	}

	depositKeys, err := parseDepositFile(selectedDeposit)
	if err != nil {
		return err
	}

	keysNum := len(depositKeys)
	overallGasPrice := gasPrice.Mul(gasPrice, big.NewInt(int64(keysNum*gasLimit)))
	overallGasPriceInt := overallGasPrice.Int64()
	overallGasPriceEth, _ := big.NewRat(overallGasPriceInt, 1_000_000_000_000_000_000).Float64()

	log.Infof("Before proceeding make sure that your private key has sufficient balance :\n"+
		"- %v ETH\n"+
		"- %v LYXe\n\n", overallGasPriceEth, keysNum*32)

	fmt.Println("Please enter your private key")
	fmt.Print(">")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	input := scanner.Text()

	privKey, err := crypto.HexToECDSA(input)
	if err != nil {
		return err
	}

	senderAddr := crypto.PubkeyToAddress(privKey.PublicKey)

	for i, key := range depositKeys {
		nonce, err := eth.PendingNonceAt(c, senderAddr)

		validatorPubKey, err := crypto.UnmarshalPubkey([]byte(key.PubKey))
		if err != nil {
			return err
		}

		validatorAddr := crypto.PubkeyToAddress(*validatorPubKey)
		fmt.Println(validatorAddr)

		err = eth.SendTransaction(c, signedTx)
		if err != nil {
			return err
		}
		fmt.Printf("Deposit %d/%d\n", i+1, keysNum)
		fmt.Println("PubKey:", validatorPubKey)
		fmt.Println("Withdraw credentials:", key.WithdrawalCredentials)
		fmt.Println("Amount:", key.Amount.String())
		fmt.Println("Signature:", key.Signature)
		fmt.Println("Deposit message root:", key.DepositMessageRoot)
		fmt.Println("Deposit data root:", key.DepositDataRoot)
		fmt.Println("Fork version:", key.ForkVersion)
		fmt.Println("Network name:", key.NetworkName)
		fmt.Println("Deposit CLI version:", key.DepositCliVersion, "\n")

		dep, err := depositMock.NewDepositMock(common.Address{}, eth)
		if err != nil {
			return err
		}

		var byteSlice [32]byte
		byteData := []byte(key.DepositDataRoot)
		if len(byteData) != 32 {
			return errors.New("asdasdasdasdasdsadasd")
		}

		for i := range byteSlice {
			byteSlice[i] = byteData[i]
		}

		opts, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(2022))
		if err != nil {
			return err
		}

		tx, err := dep.Deposit(opts, []byte(key.PubKey), []byte(key.WithdrawalCredentials), []byte(key.Signature), byteSlice)
		if err != nil {
			return err
		}

		signedTx, err := types.SignTx(tx, opts.Signer, privKey)
		if err != nil {
			return err
		}

	}

	return nil
}

func initValidator(ctx *cli.Context) error {
	initCommand := exec.Command("validator", "accounts", "import", ctx.String(validatorWalletDirFlag))

	initCommand.Stdout = os.Stdout
	initCommand.Stderr = os.Stderr

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
