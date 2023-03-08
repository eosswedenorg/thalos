package shipclient

import (
	"eosio-ship-trace-reader/config"

	"github.com/eoscanada/eos-go"
	shipclient "github.com/eosswedenorg-go/antelope-ship-client"
)

func NewClient(cfg *config.Config, chain *eos.InfoResp) (*shipclient.Client, error) {
	if cfg.StartBlockNum == config.NULL_BLOCK_NUMBER {
		if cfg.IrreversibleOnly {
			cfg.StartBlockNum = uint32(chain.LastIrreversibleBlockNum)
		} else {
			cfg.StartBlockNum = uint32(chain.HeadBlockNum)
		}
	}

	options := func(c *shipclient.Client) {
		c.StartBlock = cfg.StartBlockNum
		c.EndBlock = cfg.EndBlockNum
		c.IrreversibleOnly = cfg.IrreversibleOnly
	}

	return shipclient.NewClient(options), nil
}
