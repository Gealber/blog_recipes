package utils

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton"
)

const (
	MainnetConfig = "https://ton-blockchain.github.io/global.config.json"
	MaxRetries    = 3
	Timeout       = 50 * time.Millisecond
)

var (
	AnonymouAddress = address.MustParseAddr("EQAOQdwdw8kGftJCSFgOErM1mBjYPe4DBPq8-AhF6vr9si5N")
)

func GetAnonymousNumberAddress(ctx context.Context, api *ton.APIClient, tokenName string) (string, error) {
	index := itemIndex(tokenName)

	nftAddr, err := numberAddress(ctx, api, index)
	if err != nil {
		return "", err
	}

	return nftAddr, nil
}

func itemIndex(tokenName string) *big.Int {
	h := sha256.New()
	h.Write([]byte(tokenName))
	hash := h.Sum(nil)
	index := new(big.Int)
	index.SetBytes(hash)

	return index
}

func numberAddress(ctx context.Context, api *ton.APIClient, itemIndex *big.Int) (string, error) {
	tries := 0
START:
	block, err := api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return "", fmt.Errorf("CurrentMasterchainInfo::%v\n", err)
	}
	tries++

	res, err := api.RunGetMethod(ctx, block, AnonymouAddress, "get_nft_address_by_index", itemIndex)
	if err != nil {
		if lsErr, ok := err.(ton.LSError); ok && lsErr.Code == 651 && tries <= MaxRetries {
			fmt.Println("Sleeping ", Timeout)
			time.Sleep(Timeout)
			goto START
		}

		return "", fmt.Errorf("RunGetMethod::%v\n", err)
	}

	cs := res.MustSlice(0)
	addr := cs.MustLoadAddr()	

	return addr.String(), nil
}
