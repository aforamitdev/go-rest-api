package main

import (
	"fmt"
	"os"
	"projects/khomeawayserver/foundation/logger"

	"go.uber.org/zap"
)

func main() {

	log, err := logger.New("HOUSING-API")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

}

func run(log *zap.SugaredLogger) error {

}
