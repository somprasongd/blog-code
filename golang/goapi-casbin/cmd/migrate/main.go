package main

import (
	"flag"
	"fmt"
	"goapi/pkg/app/database"
	"goapi/pkg/common/logger"
	"goapi/pkg/config"
	"os"

	"github.com/pressly/goose/v3"
)

const dialect = "postgres"

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "./migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	switch command {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			logger.Error(fmt.Sprintf("migrate run: %v", err))
			panic(err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			logger.Error(fmt.Sprintf("migrate run: %v", err))
			panic(err)
		}
		return
	}

	appConf := config.LoadConfig()
	logger.Info("Start migration...")

	// initialize data sources
	sqlDB, err := database.NewDB(appConf)

	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	defer sqlDB.CloseDB()

	if err := goose.SetDialect(dialect); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	if err := goose.Run(command, sqlDB.DB, *dir, args[1:]...); err != nil {
		logger.Error(fmt.Sprintf("migrate run: %v", err))
		panic(err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
