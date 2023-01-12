package infrastructure

import (
	"bufio"
	"os"
	"strings"

	"go.uber.org/zap"
)

// Load is load configs from a env file.
func Load(logger *zap.Logger) {
	filePath := ".env"

	f, err := os.Open(filePath)
	if err != nil {
		logger.Error("%s", zap.Error(err))
	}
	defer f.Close()

	lines := make([]string, 0, 100)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Error("%s", zap.Error(err))
	}

	for _, l := range lines {
		pair := strings.Split(l, "=")
		os.Setenv(pair[0], pair[1])
	}
}
