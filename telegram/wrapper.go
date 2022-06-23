
package telegram

import (
    "fmt"
    _api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var _bot *_api.BotAPI
var _channel int64
var _prefix string

func Init(prefix string, id string, channel int64) (error) {

    var err error

    _bot, err = _api.NewBotAPI(id)
    if err == nil {
        _prefix = prefix
        _channel = channel
    }
    return err
}

func Send(message string) (error) {
    msg := _api.NewMessage(_channel, fmt.Sprintf("%s: %s", _prefix, message))
    _, err := _bot.Send(msg);
    return err
}
