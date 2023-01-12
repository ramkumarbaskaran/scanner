package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/scanner/app/infrastructure"
	"go.uber.org/zap"
)

func main() {
	logger := infrastructure.NewLogger("info")

	infrastructure.Load(logger)

	sqlHandler, err := infrastructure.NewSQLHandler()
	if err != nil {
		logger.Fatal("unable to connect to scylla", zap.Error(err))
	}

	infrastructure.Dispatch(logger, sqlHandler)
}
