package notifier

import (
	"github.com/nikoksr/notify"
)

type NotifierInit func() (notify.Notifier, error)

func InitNotifier(initializers ...NotifierInit) error {
	for _, init := range initializers {
		notifier, err := init()
		if err != nil {
			return err
		}

		notify.UseServices(notifier)
	}

	return nil
}
