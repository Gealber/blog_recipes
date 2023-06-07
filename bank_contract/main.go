package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/Gealber/blog_recipes/bank_contract/config"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	MainnetConfig = "https://ton-blockchain.github.io/global.config.json"
	WalletVersion = wallet.V4R2
)

func main() {
	cfg := config.Config()
	client := liteclient.NewConnectionPool()
	ctx := client.StickyContext(context.Background())

	defaultPath := path.Join("bin", "contract.cell")
	cellFile := flag.String("file", defaultPath, "provide the path to the file containing the cell")

	err := client.AddConnectionsFromConfigUrl(ctx, MainnetConfig)
	if err != nil {
		panic(err)
	}

	api := ton.NewAPIClient(client)

	w, err := wallet.FromSeed(api, cfg.TON.Wallet.Seed, WalletVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("WALLET ADDRESS", w.Address().String())

	addr, err := deployContract(ctx, w, *cellFile)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deployed contract addr:", addr.String())
}

func fileToHex(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(data)
}

func getContractData() (*cell.Cell, error) {	
	users := cell.NewDict(256)
	data := cell.BeginCell().
		MustStoreCoins(0). // bank_balance
		MustStoreCoins(0). // bank_borrowed
		MustStoreDict(users). // hasmap with users of the bank
		EndCell()

	return data, nil
}

func deployContract(
	ctx context.Context,
	w *wallet.Wallet,
	file string,
) (*address.Address, error) {
	msgBody := cell.BeginCell().EndCell()

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	contractCode, err := cell.FromBOC(data)
	if err != nil {
		return nil, err
	}

	contractData, err := getContractData()
	if err != nil {
		return nil, err
	}
	fmt.Println("Deploying contract to mainnet...")

	return w.DeployContract(
		ctx,
		tlb.MustFromTON("2"),
		msgBody,
		contractCode,
		contractData,
		true,
	)
}
