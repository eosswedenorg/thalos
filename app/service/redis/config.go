package redis

type Config struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	CacheID  string `json:"cache_id"`
	Prefix   string `json:"prefix"`
}

var DefaultConfig = Config{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
	Prefix:   "ship",
}
