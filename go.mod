module eosio-ship-trace-reader

go 1.14

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.10 // indirect
	github.com/eoscanada/eos-go v0.10.2
	github.com/eosswedenorg-go/eos-ship-client v0.0.0-20220202131720-5e908c506be2
	github.com/eosswedenorg-go/pid v1.0.1
	github.com/go-redis/cache/v8 v8.4.3
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/pborman/getopt/v2 v2.1.0
	github.com/streamingfast/logging v0.0.0-20211221170249-09a6ecb200a0 // indirect
	github.com/teris-io/shortid v0.0.0-20201117134242-e59966efd125 // indirect
	github.com/tidwall/gjson v1.14.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.20.0 // indirect
	golang.org/x/crypto v0.0.0-20220131195533-30dcbda58838 // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/tools v0.1.9 // indirect
	internal/abi_cache v0.0.0
)

replace internal/abi_cache => ./internal/abi_cache
