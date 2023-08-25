package main

import (
	"fmt"
	"os/signal"
	"strconv"
	"syscall"

	"net/http"
	"os"
	"time"
	"virtual-orb/pkg/platform"
	"virtual-orb/pkg/service"

	"github.com/bwmarrin/snowflake"
	"github.com/joho/godotenv"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

// main is the entry point for the application. It initializes the logger,
// loads environment variables, sets up various services including an HTTP
// client and a circuit breaker, and periodically reports system status and
// simulates a signup process.
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error creating logger")
		os.Exit(1)
	}
	defer logger.Sync()

	err = godotenv.Load()
	if err != nil {
		logger.Error("Loading env file failed",
			zap.Error(err))
		os.Exit(1)
	}

	orbIDStr := GetEnvWithDefault("ORB_ID", "1")
	orbID, _ := strconv.ParseInt(orbIDStr, 10, 64)
	signKey := GetEnvWithDefault("SIGN_KEY", "default-secret-key")
	cbTimeoutStr := GetEnvWithDefault("CB_TIMEOUT", "60s")
	cbTimeout, _ := time.ParseDuration(cbTimeoutStr)
	cbMaxRequestsStr := GetEnvWithDefault("CB_MAX_REQUESTS", "5")
	cbMaxRequests, _ := strconv.Atoi(cbMaxRequestsStr)
	cbIntervalStr := GetEnvWithDefault("CB_INTERVAL", "60s")
	cbInterval, _ := time.ParseDuration(cbIntervalStr)
	tickerIntervalStr := GetEnvWithDefault("TICKER_INTERVAL", "5s")
	tickerInterval, _ := time.ParseDuration(tickerIntervalStr)
	baseURL := GetEnvWithDefault("BASE_URL", "http://mock-uniqueness-service:8001")
	snowflakeNode, err := snowflake.NewNode(orbID)
	if err != nil {
		logger.Error("Creating snowflake node failed",
			zap.Error(err))
		os.Exit(1)
	}
	httpClient := http.DefaultClient
	cbSettings := gobreaker.Settings{
		Name:        "HTTP Request Circuit Breaker",
		Timeout:     cbTimeout,
		MaxRequests: uint32(cbMaxRequests),
		Interval:    cbInterval,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)

			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	cb := gobreaker.NewCircuitBreaker(cbSettings)
	requestSvc := service.NewRequestSvc(baseURL, httpClient, cb)
	signUp := service.NewSignUpSvc(signKey, snowflakeNode, requestSvc)
	systemInfo := platform.NewSystemInfo()
	status := service.NewStatusSvc(requestSvc, systemInfo)

	ticker := time.NewTicker(tickerInterval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				//1. Periodically report status
				err := status.Report()
				if err != nil {
					logger.Error("Reporting status failed",
						zap.Error(err))
				} else {
					logger.Info("Reporting status succeeded")
				}

				//2. Periodically simulate a signup and submit images ..
				imgData, err := platform.GenerateRandomImageData()
				if err != nil {
					logger.Error("Scanning iris image failed",
						zap.Error(err))
				}

				err = signUp.SignUp(imgData)
				if err != nil {
					logger.Error("Signing up failed",
						zap.Error(err))
				} else {
					logger.Info("Signing up succeeded")
				}

			}

		}
	}()

	// Implement graceful shutdown incase jobs were doing work at time of stoppage
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	logger.Info("Gracefully shutting down server...")
	ticker.Stop()
	done <- true
	os.Exit(0)

}

// GetEnvWithDefault fetches the value of an environment variable.
// If the variable isn't set, it returns a provided default value.
//
// Parameters:
//
//	key: Name of the environment variable.
//	defaultValue: Value to return if the environment variable is not set.
//
// Returns:
//
//	Value of the environment variable or defaultValue if it's not set.
func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
