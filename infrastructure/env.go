package infrastructure

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Load app configuration variables from an env file.
func Bootstrap(/** logger should be parameter here */) {
	filePath := ".env"

	envFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open env file %s", err)
	}
	defer envFile.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(envFile)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading env file %s", err)
	}

	for _, line := range lines {
		pair := strings.Split(line, "=")
		os.Setenv(pair[0], pair[1])
	}
}