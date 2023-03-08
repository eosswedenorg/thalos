package telegram

import (
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

type Config struct {
	Id      string `json:"id"`
	Channel int64  `json:"channel"`
}

func Notifier(conf Config) func() (notify.Notifier, error) {
	return func() (notify.Notifier, error) {
		tg, err := telegram.New(conf.Id)
		if err == nil {
			tg.AddReceivers(conf.Channel)
		}
		return tg, err
	}
}
