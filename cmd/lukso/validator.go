package main

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/m8b-dev/lukso-cli/contracts/bindings"
	"github.com/urfave/cli/v2"
	"math/big"
	"os"
	"os/exec"
)

const (
	gasLimit               = 21_000
	depositContractAddress = "0x75D1f4695Eb87d60eD4EAE2c0CF05e7428Fa4b5F"
	lyxeContractAddress    = "0x7A2AC110202ebFdBB5dB15Ea994ba6bFbFcFc215"
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
	singleTxGasPrice := big.NewInt(0).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	overallGasPrice := big.NewInt(0).Mul(singleTxGasPrice, big.NewInt(int64(keysNum)))
	overallGasPriceInt := overallGasPrice.Int64()
	overallGasPriceEth, _ := big.NewRat(overallGasPriceInt, 1_000_000_000_000_000_000).Float64()

	log.Infof("Before proceeding make sure that your private key has sufficient balance:\n"+
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
	chainId, err := eth.ChainID(c)
	if err != nil {
		return err
	}

	for i, key := range depositKeys {
		nonce, err := eth.PendingNonceAt(c, senderAddr)

		opts, err := bind.NewKeyedTransactorWithChainID(privKey, chainId)
		if err != nil {
			return err
		}

		lyxMock, err := bindings.NewLYXe(common.HexToAddress(lyxeContractAddress), eth)
		if err != nil {
			return err
		}

		fmt.Printf("Deposit %d/%d\n", i+1, keysNum)
		fmt.Println("Withdraw credentials:", key.WithdrawalCredentials)
		fmt.Println("Amount:", key.Amount.String())
		fmt.Println("Signature:", key.Signature)
		fmt.Println("Deposit message root:", key.DepositMessageRoot)
		fmt.Println("Deposit data root:", key.DepositDataRoot)
		fmt.Println("Fork version:", key.ForkVersion)
		fmt.Println("Network name:", key.NetworkName)
		fmt.Println("Deposit CLI version:", key.DepositCliVersion, "\n")

		opts.From = senderAddr
		//opts.GasPrice = singleTxGasPrice
		opts.Value = big.NewInt(0)
		opts.Nonce = big.NewInt(int64(nonce))

		depositData, err := prepareDepositData(key)
		if err != nil {
			return err
		}

		tx, err := lyxMock.Send(opts, common.HexToAddress(depositContractAddress), big.NewInt(0).Mul(big.NewInt(32), big.NewInt(1000000000000000000)), depositData)
		if err != nil {
			return err
		}

		fmt.Println("GAS PRICE:", tx.GasPrice().Int64(), "GAS TIP:", tx.GasTipCap().Int64(), "GAS FEE:", tx.GasFeeCap().Int64())
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

func prepareDepositData(key DepositDataKey) (depositData []byte, err error) {
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
	depositData = append(depositData, byte(32))

	return
}
