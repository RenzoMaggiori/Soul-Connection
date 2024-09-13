package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"soul-connection.com/migration/src/migration"
	"soul-connection.com/migration/src/parser"
)

func main() {
	params, err := parser.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(*params.EnvPath)
	if err != nil {
		log.Fatal(err)
	}

	migration.Start()
}
