package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Gealber/fragmentbuyer/utils"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

const (
	MainnetConfig = "https://ton-blockchain.github.io/global.config.json"
	WalletVersion = wallet.V4R2
)

func main() {
	client := liteclient.NewConnectionPool()
	ctx := client.StickyContext(context.Background())
	err := client.AddConnectionsFromConfigUrl(ctx, MainnetConfig)
	if err != nil {
		panic(err)
	}

	api := ton.NewAPIClient(client)
	// getNumberAddress()
	// getOnSaleNumbers()
	// scrapeNumbers()	
	buyNumbers(ctx, api)
	// auctionNumbers()
}

func getNumberAddress(ctx context.Context, api *ton.APIClient) {	
	tokenName := "88805764931"
	nftAddr, err := utils.GetAnonymousNumberAddress(ctx, api, tokenName)
	if err != nil {
		panic(err)
	}

	log.Println(nftAddr)
}

func getOnSaleNumbers() {
	numbers, err := utils.ScrapeFragmentMarket()
	if err != nil {
		panic(err)
	}

	fmt.Printf("NUMBERS: %+v\n", numbers)
}

func scrapeNumbers(ctx context.Context, api *ton.APIClient) {
	numbers, err := utils.ScrapeFragmentMarket()
	if err != nil {
		panic(err)
	}

	f, err := os.Create("numbers.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := 0; i < len(numbers); i++ {
		nftAddr, err := utils.GetAnonymousNumberAddress(ctx, api, numbers[i].Name)
		if err != nil {
			log.Println("Couldn't get the address for ", numbers[i].Name)
			continue
		}
		numbers[i].Address = nftAddr

		f.WriteString(fmt.Sprintf("%s,%d,%s\n", numbers[i].Name, numbers[i].Value, numbers[i].Address))
	}
}

func buyNumbers(ctx context.Context, api *ton.APIClient) {	
	seed := strings.Split("<SEED>", " ")
	w, err := wallet.FromSeed(api, seed, WalletVersion)
	if err != nil {
		panic(err)
	}

	log.Println(w.Address())

	numbers, err := utils.ScrapeFragmentMarket()
	if err != nil {
		panic(err)
	}

	numbersToBuy := 1
	if len(numbers) < numbersToBuy {
		panic(errors.New("insufficient numbers scraped"))
	}

	for i := 0; i < numbersToBuy; i++ {
		nftAddr, err := utils.GetAnonymousNumberAddress(ctx, api, numbers[i].Name)
		if err != nil {
			log.Println("Couldn't get the address for ", numbers[i].Name)
			continue
		}

		numbers[i].Address = nftAddr
	}

	// perform buy of first numbersToBuy = 1 numbers
	log.Println("Buying numbers...")
	err = utils.BuyNumbers(ctx, w, api, numbers[:numbersToBuy])
	if err != nil {
		panic(err)
	}
}
