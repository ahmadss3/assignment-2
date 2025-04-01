// File: assignment-2/tools/porttools.go
package tools

import "os"

// GetServerPort returns the port from the environment variable PORT, or uses
// the provided fallback if PORT is not set or is empty.

func GetServerPort(fallback string) string {
	port := os.Getenv("PORT")
	if port == "" {
		return fallback
	}
	return port
}
