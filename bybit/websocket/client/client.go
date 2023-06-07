package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Bybit organize the topics to which
// we can subscribe into Public and Private topics.
// Basically topics to which you can subscribe without authentication
// and the ones you can only subscribe with Authentication.
type ChannelType string

// There are several types of trades on bybit
// we are going to use only spot.
type TradesType string

const (
	PublicChannel  ChannelType = "public"
	PrivateChannel ChannelType = "private"
	Spot           TradesType  = "spot"

	// topic for getting latest prices of
	// TONUSDT pair in spot trading.
	TickersTONUSDTTopic = "tickers.TONUSDT"

	APIVersion           = "v5"
	ByBitWebsocketDomain = "stream.bybit.com"

	// In order to comply with bybit documentation is necessary to
	// send every 20 seconds a ping command to avoid disconnection.
	// A basic heartbeat mechanism.
	PingTimeout = 20

	// Even providing the heartbeat you will suffer disconnections
	// simply because that's how programming works, you should be prepared
	// for every kind of errors.
	MaxRetrialConnections = 10
)

// Client represents connection with ByBit Websocket API.
type Client struct {
	conn *websocket.Conn
}

func NewClient() (*Client, error) {
	conn, err := connect()
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func path(channelType ChannelType, operation TradesType) string {
	return fmt.Sprintf("/%s/%s/%s", APIVersion, channelType, operation)
}

// connect to bybit websocket stream.
func connect() (*websocket.Conn, error) {
	spotPath := path(PublicChannel, Spot)

	u := url.URL{Scheme: "wss", Host: ByBitWebsocketDomain, Path: spotPath}
	fmt.Printf("connecting to %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return conn, err
}

func (c *Client) sendPing() error {
	pingReq := Request{
		ReqID: "100001",
		Op:    "ping",
	}

	return c.conn.WriteJSON(&pingReq)
}

// replace connection with a new one
func (c *Client) reconnect() error {
	conn, err := connect()
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// Run connect to bybit websocket, general idea of what it does.
// 1. Subscribe to tickers
// 2. Read message from websocket.
// 3. Send every 20 seconds a ping(heartbeat), to avoid disconnections.
// 4. In case of abnormal close of connection, performs a reconnection.
// 5. In case the reconnection exceed the max allowed, shut the program.
// 6. Also listen to Ctr+C commands to shutdown gratefully.
func (c *Client) Run(
	ctx context.Context,
	subscriptions []Request,	
) error {
	// interupt signal necessary for catching Ctr+C
	// and shut down gratefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// ticker for sending ping command every 20 seconds
	// for the heartbeat mechanism.
	pingTicker := time.NewTicker(PingTimeout * time.Second)
	// we are going to report our errors in this channel
	errChn := make(chan error)
	done := make(chan struct{})
	var connections int = 0

	// in case we are going to have a disconnection we are going
	// back to this tag to reconnect.
CONNECTION: 
	if connections > 0 {
		err := c.reconnect()
		if err != nil {
			return err
		}
	}
	
	defer func() {
		pingTicker.Stop()
		close(errChn)
		c.conn.Close()
	}()

	// let's increase our connection counter
	// to keep a record of our max retries
	connections++

	// this gourutine will process every message incomming
	// related to the subscriptions passed.
	go func() {
		defer close(done)
		err := c.processRead(ctx, done, subscriptions)
		if err != nil {
			errChn <- err

			return
		}
	}()

	// wait for possible errors and interrupt signals.
	// send ping command every PingTimeout(20) seconds.
	for {
		select {
		case err := <-errChn:
			// not all the errors are retriable
			if waitTime, ok := retriableError(err, connections); ok {
				if connections > MaxRetrialConnections {
					return err
				}

				fmt.Println("RECONNECTING.....")
				fmt.Printf("Sleeping %v milliseconds\n", waitTime)
				time.Sleep(waitTime)
				// go to tag defined above
				goto CONNECTION
			}

			return err

		// in case of Ctr+C let's handle the connection
		case <-interrupt:
			return c.handleInterruptSignal(done)

		// when ticker is triggered send a ping command
		case <-pingTicker.C:
			err := c.sendPing()
			if err != nil {
				return err
			}
		}
	}
}

func (c *Client) processMsg(ctx context.Context, data []byte) error {
	var msg PublicResponse

	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	fmt.Printf("MSG: %s\n", string(data))
	return nil	
}

func (c *Client) processRead(
	ctx context.Context,
	done chan struct{},	
	subscriptions []Request,	
) error {
	// first ping to send.
	c.sendPing()

	for _, subscription := range subscriptions {
		err := c.conn.WriteJSON(subscription)
		if err != nil {
			return fmt.Errorf("sending subscription %w", err)
		}
	}

	for {
		select {
		case <-done:
			return nil
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("reading %w", err)
			}

			err = c.processMsg(ctx, message)
			if err != nil {
				fmt.Println("ERR PROCESSING READ: ", err.Error())
			}
		}
	}
}

func (c *Client) handleInterruptSignal(
	done chan struct{},	
) error {
	fmt.Println("Closing connection...it might take a few seconds until the next tickers comes")
	done <- struct{}{}
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	select {
	case <-time.After(time.Second):
	}

	return nil
}


func retriableError(err error, connections int) (time.Duration, bool) {
	waitTime := time.Duration(connections) * 500 * time.Millisecond
	// 1006(CloseAbnormalClosure) checking if was an abnormal closure.
	isAbnormalClosure := websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure)
	// 1001(CloseGoingAway) indicates that an endpoint is "going away", such as a server
	// going down or a browser having navigated away from a page.
	isGoingAway := websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway)

	isRetriable := isAbnormalClosure || isGoingAway

	// 1013(CloseTryAgainLater) not defined in rfc6455 but can be faced as well.
	if websocket.IsUnexpectedCloseError(err, websocket.CloseTryAgainLater) {
		isRetriable = true
		waitTime = time.Duration(connections) * time.Second
	}

	return waitTime, isRetriable
}
