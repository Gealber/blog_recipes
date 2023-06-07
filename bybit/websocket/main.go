package main

import (
	"context"

	"github.com/Gealber/blog_recipes/bybit/websocket/client"
)

func main() {
	clt, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tickerSubsciption := client.Request{
		Op: "subscribe",
		Args: []interface{}{
			client.TickersTONUSDTTopic,
		},
	}
	subscriptions := []client.Request{
		tickerSubsciption,
	}

	clt.Run(ctx, subscriptions)
}
