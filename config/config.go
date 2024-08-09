package config

type Config struct {
	Addr          string
	PoWDifficulty int64
}

func NewConfig(addr string, powDiff int64) *Config {
	return &Config{Addr: addr, PoWDifficulty: powDiff}
}
