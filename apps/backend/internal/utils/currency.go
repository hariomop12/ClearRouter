package utils

import (
	"os"
	"strconv"
)

// GetCurrency returns the configured currency code (defaults to INR)
func GetCurrency() string {
	cur := os.Getenv("CURRENCY")
	if cur == "" {
		return "INR"
	}
	return cur
}

// usdToInrRate returns the configured USD->INR rate (defaults to 83.0)
func usdToInrRate() float64 {
	val := os.Getenv("USD_TO_INR")
	if val == "" {
		return 83.0
	}
	rate, err := strconv.ParseFloat(val, 64)
	if err != nil || rate <= 0 {
		return 83.0
	}
	return rate
}

// ConvertUSDToConfigured converts a USD amount into the configured currency.
// Currently supports INR; if currency is not INR, returns input unchanged.
func ConvertUSDToConfigured(amountUSD float64) float64 {
	if GetCurrency() == "INR" {
		return amountUSD * usdToInrRate()
	}
	return amountUSD
}
