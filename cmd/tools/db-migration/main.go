package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/joho/godotenv"
	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/pkg/config"
	"github.com/woningfinder/woningfinder/pkg/logging"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../../.env"); err != nil {
		_ = config.MustGetString("APP_NAME")
	}
}

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run *.go <command> [args]
`

func main() {
	flag.Usage = usage
	flag.Parse()

	logger := logging.NewZapLogger(config.GetBoolOrDefault("APP_DEBUG", false), config.MustGetString("SENTRY_DSN"))
	dbClient := bootstrap.CreateDBClient(logger)
	db := dbClient.Conn()

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		logger.Sugar().Fatal(err)
	}
	if newVersion != oldVersion {
		logger.Sugar().Infof("database schema migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		logger.Sugar().Infof("database schema version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}
