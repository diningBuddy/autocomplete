package env

import (
	"os"
)

// GetString get string value from os environment variables by key parameter.
// It should give fallback value when the value of the key does not exist in environment variables.
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}
