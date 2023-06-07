package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	startAuction       = 0x487a8e81
	nanoTons     int64 = 1e9
)

func buyNumber(
	ctx context.Context,
	w *wallet.Wallet,
	tokenAddress string,
	value int,
) error {
	nftAddr, err := address.ParseAddr(tokenAddress)
	if err != nil {
		return err
	}

	transfer := wallet.SimpleMessage(nftAddr, tlb.MustFromTON(strconv.Itoa(value)), nil)

	return w.Send(ctx, transfer, true)
}

func auctionNumber(
	ctx context.Context,
	w *wallet.Wallet,
	tokenAddress string,
	value int,
) error {
	nftAddr, err := address.ParseAddr(tokenAddress)
	if err != nil {
		return err
	}
	
	currentTime := time.Now().UTC().Unix()
	valueNanoTons := uint64(int64(value) * nanoTons)
	auctionConfig := cell.BeginCell().
		// beneficiary_address
		MustStoreAddr(w.Address()).
		// initial_min_bid, value in nanotons
		MustStoreCoins(valueNanoTons).
		// max_bid
		MustStoreCoins(valueNanoTons).
		// min_bid_step = 5
		MustStoreUInt(5, 8).
		// min_extend_time = 1 hour
		MustStoreUInt(3600, 32).
		// duration = 7 days
		MustStoreUInt(604800, 32).
		EndCell()

	payload := cell.BeginCell().
		// op = op::teleitem_start_auction
		MustStoreUInt(startAuction, 32).
		// query_id = what ever
		MustStoreUInt(uint64(currentTime), 64).
		// new_auction_config
		MustStoreRef(auctionConfig).
		EndCell()

	transfer := wallet.SimpleMessage(nftAddr, tlb.MustFromTON("0.1"), payload)

	return w.Send(ctx, transfer, true)
}

func BuyNumbers(
	ctx context.Context,
	w *wallet.Wallet,
	api *ton.APIClient,
	numbers []AnonNumbers,
) error {
	// check available balance
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return fmt.Errorf("CurrentMasterchainInfo::%v\n", err)
	}

	balanceCoins, err := w.GetBalance(ctx, block)
	if err != nil {
		return fmt.Errorf("GetBalance::%v+\n", err)
	}

	balanceStr := balanceCoins.TON()
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return err
	}

	// check if we have enough balance
	if enoughBudget(numbers, int(balance)) {
		for _, number := range numbers {
			if !number.OnSell {
				err = buyNumber(ctx, w, number.Address, number.Value)
				if err != nil {
					return err
				}

				// add this number into db maybe
				log.Printf("Number %s bought\n", number.Name)
			}
		}
	}

	return nil
}

func AuctionNumbers(
	ctx context.Context,
	w *wallet.Wallet,
	api *ton.APIClient,
	numbers []AnonNumbers,
) error {
	for i := 0; i < len(numbers); i++ {
		err := auctionNumber(ctx, w, numbers[i].Address, numbers[i].valueToSell())
		if err != nil {
			return err
		}

		numbers[i].OnSell = true

		// mark in db this number as on auction
		log.Printf("Number %s on auction\n", numbers[i].Name)
	}

	return nil
}
