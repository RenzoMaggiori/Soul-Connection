package migration

import (
	"database/sql"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"soul-connection.com/api/src/database"
	"soul-connection.com/api/src/file-storage"
	"soul-connection.com/api/src/lib"
)

type ApiCredentials struct {
	lib.LoginCredentials
	Jwt string
}

type Ids struct {
	old int
	new int
}

type MigrationFunc func(*sql.DB, *mongo.Database, *ApiCredentials) error

const migrationRate int = 60 * 24

func Start() {
	ticker := time.NewTicker(time.Duration(migrationRate) * time.Minute)
	quit := make(chan struct{})
	credentials := lib.LoginCredentials{
		XGroupAuthentication: os.Getenv("API_KEY"),
		AuthEmail:            os.Getenv("API_EMAIL"),
		AuthPassword:         os.Getenv("API_PASSWORD"),
	}
	database, err := database.Open(database.ConnectionString())
	if err != nil {
		lib.ServerLog("ERROR", err)
		return
	}
	defer database.Close()

	mongoClient, ctx, err := filestorage.Open(filestorage.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)
	fileStorage := mongoClient.Database("soul-connection-files")

	lib.ServerLog("INFO", "Running migration")
	run(database, fileStorage, credentials, quit)
	lib.ServerLog("INFO", "Migration complete")

	go func() {
		for {
			select {
			case <-ticker.C:
				lib.ServerLog("INFO", "Running migration")
				run(database, fileStorage, credentials, quit)
				lib.ServerLog("INFO", "Migration complete")
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	select {}
}

func run(database *sql.DB, fileStorage *mongo.Database, loginCredentials lib.LoginCredentials, quit chan struct{}) {
	jwt, err := lib.Auth(loginCredentials)

	if err != nil {
		lib.ServerLog("ERROR", err)
		quit <- struct{}{}
		return
	}
	apiCredentials := ApiCredentials{
		LoginCredentials: loginCredentials,
		Jwt:              jwt,
	}

	migrations := []MigrationFunc{
		migrateEmployees,
		migrateCustomers,
		migrateEncounters,
		migrateTips,
		migrateEvents,
	}

	for _, migration := range migrations {
		err := migration(database, fileStorage, &apiCredentials)
		if err != nil {
			lib.ServerLog("WARNING", err)
		}
	}
}
