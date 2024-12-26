package config

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Config struct {
	DBManagerAddress          string        `json:"db_manager_address"`
	CoingeckoServiceAddress   string        `json:"coingecko_service_address"`
	DBManagerTimeout          time.Duration `json:"db_manager_timeout"`
	CoingeckoServiceTimeout   time.Duration `json:"coingecko_service_timeout"`
	CoingeckoPollingFrequency time.Duration `json:"coingecko_polling_frequency"`
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig(filePath string) (*Config, error) {
	var err error
	once.Do(
		func() {
			file, err := os.Open(filePath)
			if err != nil {
				return
			}
			defer file.Close()

			decoder := json.NewDecoder(file)
			instance = &Config{}
			err = decoder.Decode(instance)
			if err != nil {
				instance = nil
				return
			}

			instance.DBManagerTimeout *= time.Second
			instance.CoingeckoServiceTimeout *= time.Second
			instance.CoingeckoPollingFrequency *= time.Second
		},
	)
	return instance, err
}

func GetConfig() *Config {
	return instance
}
