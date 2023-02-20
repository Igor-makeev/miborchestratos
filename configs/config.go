package config

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort string
	DBAddress  string
	Shards     map[int]string
	AuthConfig
}

type AuthConfig struct {
	Salt       string
	SigningKey string
	TokenTTL   time.Duration
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Salt:       os.Getenv("SALT"),
		SigningKey: os.Getenv("SIGNING_KEY"),
		TokenTTL:   6 * time.Hour,
	}
}

func NewConfig() *Config {
	cfg := &Config{
		AuthConfig: *NewAuthConfig(),
		ServerPort: os.Getenv("ORCHESTRATOR_SERVER_PORT"),
		DBAddress:  os.Getenv("DATABASE_URI"),
		Shards:     make(map[int]string),
	}
	cfg.Shards[0] = os.Getenv("SHARD_ADRESS_1")
	cfg.Shards[1] = os.Getenv("SHARD_ADRESS_2")
	logrus.Printf("config:serverport:%v", cfg.ServerPort)
	logrus.Printf("config:shard number:%v,address:%v", 1, cfg.Shards[0])
	logrus.Printf("config:shard number:%v,address:%v", 2, cfg.Shards[1])
	logrus.Printf("config:number of shards:%v", len(cfg.Shards))
	return cfg
}
