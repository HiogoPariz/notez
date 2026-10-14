package config

import "os"

// Get returns the value of the environment variable key,
// or fallback when it is unset or empty.
func Get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
