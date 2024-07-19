package main

import (
	"context"
	"time"

	antelopeapi "github.com/shufflingpixels/antelope-go/api"
	log "github.com/sirupsen/logrus"
)

// "Clever" way to make sure we only call the api once.
// Store a info pointer outside the returned closure.
// that pointer will live as long as the closure lives.
// and inside the closure we will reference the pointer and only
// call the api if it is nil.
func chainInfoOnce(api *antelopeapi.Client) func() *antelopeapi.Info {
	var info *antelopeapi.Info
	return func() *antelopeapi.Info {
		if info == nil {
			log.WithField("api", api.Url).Info("Get chain info from api")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			result, err := api.GetInfo(ctx)
			if err != nil {
				log.WithError(err).Fatal("Failed to call eos api")
				return nil
			}

			info = &result
		}
		return info
	}
}
