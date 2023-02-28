package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

const gasLimit = 21_000

type DepositDataKey struct {
	Pubkey                string   `json:"pubkey"`
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
	eth, err := ethclient.Dial("https://mainnet.infura.io/v3/ab367e2f12804177b29bade35c399475")
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

	t := &types.LegacyTx{
		Nonce:    0,
		GasPrice: nil,
		Gas:      0,
		To:       nil,
		Value:    nil,
		Data:     nil,
		V:        nil,
		R:        nil,
		S:        nil,
	}

	tx, err := types.SignTx(types.NewTx(t), nil, privKey)

	err = eth.SendTransaction(c, tx)
	if err != nil {
		return err
	}

	for i, key := range depositKeys {
		fmt.Printf("Deposit %d/%d\n", i+1, keysNum)
		fmt.Println("Pubkey:", key.Pubkey)
		fmt.Println("Withdraw credentials:", key.WithdrawalCredentials)
		fmt.Println("Amount:", key.Amount.String())
		fmt.Println("Signature:", key.Signature)
		fmt.Println("Deposit message root:", key.DepositMessageRoot)
		fmt.Println("Deposit data root:", key.DepositDataRoot)
		fmt.Println("Fork version:", key.ForkVersion)
		fmt.Println("Network name:", key.NetworkName)
		fmt.Println("Deposit CLI version:", key.DepositCliVersion, "\n")
	}

	return err
}

func parseDepositFile(depositFilePath string) (keys []DepositDataKey, err error) {
	f, err := os.ReadFile(depositFilePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &keys)

	return
}
