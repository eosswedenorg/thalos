package ship

import (
	"context"
	"fmt"
	"time"

	"github.com/nikoksr/notify"

	log "github.com/sirupsen/logrus"

	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

type Client struct {
	sh  *shipclient.Client
	api string

	done chan interface{}
}

func New(api string, client *shipclient.Client) *Client {
	return &Client{
		api:  api,
		sh:   client,
		done: make(chan interface{}),
	}
}

func (c *Client) connect() {
	var recon_cnt uint = 0

	for {
		recon_cnt++
		log.Infof("Connecting to ship at: %s (Try %d)", c.api, recon_cnt)
		err := c.sh.Connect(c.api)
		if err != nil {
			log.Println(err)

			if recon_cnt >= 3 {
				msg := fmt.Sprintf("Failed to connect to ship at '%s'", c.api)
				if err := notify.Send(context.Background(), "Ship_reader", msg); err != nil {
					log.WithError(err).Error("Failed to send notification")
				}
				recon_cnt = 0
			}

			log.Info("Trying again in 5 seconds ....")
			time.Sleep(5 * time.Second)
			continue
		}

		err = c.sh.SendBlocksRequest()
		if err != nil {
			log.Println(err)
			return
		}

		// Connected
		log.Infof("Connected, Start: %d, End: %d", c.sh.StartBlock, c.sh.EndBlock)
		break
	}
}

func (c *Client) read() {
	err := c.sh.Read()
	if err != nil {
		if shErr, ok := err.(shipclient.ClientError); ok {

			// Bail out if socket is closed
			if shErr.Type == shipclient.ErrSockClosed {
				log.Info("Socket closed, Exiting")
				return
			}

			// Reconnect
			if shErr.Type == shipclient.ErrSockRead || shErr.Type == shipclient.ErrNotConnected {
				c.connect()
			}
		}

		log.WithError(err).Error("Failed to read from ship")
	}
}

func (c *Client) Run() error {
	defer c.Close()

	for {
		select {
		case <-c.done:
			return nil
		default:
			c.read()
		}
	}
}

func (c *Client) Close() {
	err := c.sh.Shutdown()
	if err != nil {
		log.WithError(err).Error("Failed to send close message")
	}

	close(c.done)
}
