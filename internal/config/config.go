package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	DBManagerAddress          string
	CoingeckoServiceAddress   string
	DBManagerTimeout          time.Duration
	CoingeckoServiceTimeout   time.Duration
	CoingeckoPollingFrequency time.Duration
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() (*Config, error) {
	var err error
	once.Do(
		func() {
			dbManagerTimeout, err := strconv.Atoi(os.Getenv("DB_MANAGER_TIMEOUT"))
			if err != nil {
				dbManagerTimeout = 30
			}
			coingeckoServiceTimeout, err := strconv.Atoi(os.Getenv("COINGECKO_SERVICE_TIMEOUT"))
			if err != nil {
				coingeckoServiceTimeout = 30
			}
			coingeckoPollingFrequency, err := strconv.Atoi(os.Getenv("COINGECKO_POLLING_FREQUENCY"))
			if err != nil {
				coingeckoPollingFrequency = 60
			}

			instance = &Config{
				DBManagerAddress:          os.Getenv("DB_MANAGER_ADDRESS"),
				CoingeckoServiceAddress:   os.Getenv("COINGECKO_SERVICE_ADDRESS"),
				DBManagerTimeout:          time.Duration(dbManagerTimeout) * time.Second,
				CoingeckoServiceTimeout:   time.Duration(coingeckoServiceTimeout) * time.Second,
				CoingeckoPollingFrequency: time.Duration(coingeckoPollingFrequency) * time.Second,
			}
		},
	)
	return instance, err
}

func GetConfig() *Config {
	return instance
}
