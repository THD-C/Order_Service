package config

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	Addr                       string
	Port                       string
	ServiceName                string
	PrometheusPort             string
	DBManagerAddress           string
	CoingeckoServiceAddress    string
	DBManagerTimeout           time.Duration
	CoingeckoServiceTimeout    time.Duration
	CoingeckoPollingFrequency  time.Duration
	PendingOrderCheckFrequency time.Duration
}

var (
	instance *Config
	once     sync.Once
)

func LoadConfig() (*Config, error) {
	var err error
	once.Do(
		func() {
			applicationAddr := os.Getenv("APPLICATION_ADDR")
			if applicationAddr == "" {
				applicationAddr = "0.0.0.0"
			}

			applicationPort := os.Getenv("APPLICATION_PORT")
			if applicationPort == "" {
				applicationPort = "50051"
			}

			prometheusPort := os.Getenv("PROMETHEUS_PORT")
			if prometheusPort == "" {
				prometheusPort = "8111"
			}

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

			pendingOrderCheckFrequency, err := strconv.Atoi(os.Getenv("PENDING_ORDER_CHECK_FREQUENCY"))
			if err != nil {
				pendingOrderCheckFrequency = 60
			}

			instance = &Config{
				Addr:                       applicationAddr,
				Port:                       applicationPort,
				ServiceName:                "order_service",
				PrometheusPort:             prometheusPort,
				DBManagerAddress:           os.Getenv("DB_MANAGER_ADDRESS"),
				CoingeckoServiceAddress:    os.Getenv("COINGECKO_SERVICE_ADDRESS"),
				DBManagerTimeout:           time.Duration(dbManagerTimeout) * time.Second,
				CoingeckoServiceTimeout:    time.Duration(coingeckoServiceTimeout) * time.Second,
				CoingeckoPollingFrequency:  time.Duration(coingeckoPollingFrequency) * time.Second,
				PendingOrderCheckFrequency: time.Duration(pendingOrderCheckFrequency) * time.Second,
			}
		},
	)
	return instance, err
}

func GetConfig() *Config {
	return instance
}
