package main

import (
	"fmt"
	"os"
	"projects/khomeawayserver/foundation/logger"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func main() {

	log, err := logger.New("HOUSING-API")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
	}

}

func run(log *zap.SugaredLogger) error {

	opt := maxprocs.Logger(log.Infof)

	if _, err := maxprocs.Set(opt); err != nil {
		log.Errorf("maxprocess: %w", err)
	}

	return nil
}
