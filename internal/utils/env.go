package utils

import "os"

// GetEnvOrDefault returns the value of an environment variable or a default value if the environment variable is not set.
func GetEnvOrDefault(env string, defaultValue string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return defaultValue
}
