package main

import (
	"flag"
	"fmt"
	"github.com/go-pg/migrations/v8"
	"log"
	"os"
	"treasury/src/driver/postgres"
)

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
	log.Println("Start migrations...")
	flag.Usage = usage
	flag.Parse()

	pgdb, err := postgres.StartDB()
	if err != nil {
		log.Printf("error: %v", err)
		panic("error starting the database")
	} else {
		log.Println("DB connected")
	}

	oldVersion, newVersion, err := migrations.Run(pgdb, flag.Args()...)
	if err != nil {
		exitf(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}
}

func usage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
