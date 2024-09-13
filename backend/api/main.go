package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"soul-connection.com/api/src/database"
	"soul-connection.com/api/src/endpoints"
	filestorage "soul-connection.com/api/src/file-storage"
	"soul-connection.com/api/src/lib"
	"soul-connection.com/api/src/middleware"
	"soul-connection.com/api/src/parser"
	"soul-connection.com/api/src/server"
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

	database, err := database.Open(database.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	mongoClient, ctx, err := filestorage.Open(filestorage.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(ctx)
	fileStorage := mongoClient.Database("soul-connection-files")

	router, err := endpoints.CreateRouter(database, fileStorage)
	if err != nil {
		log.Fatal(err)
	}

	corsRouter := mux.NewRouter()
	corsRouter.Use(middleware.CorsMiddleware)
	corsRouter.PathPrefix("/").Handler(router)

	apiServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *params.Port),
		Handler: corsRouter,
	}

	running := make(chan struct{})
	go server.Start(apiServer, running)
	defer server.Stop(apiServer)
	<-running
	initalLog(apiServer, corsRouter)
}

func initalLog(server *http.Server, router *mux.Router) {
	fmt.Println(`
   _____  __________ .___
  /  _  \ \______   \|   |
 /  /_\  \ |     ___/|   |
/    |    \|    |    |   |
\____|__  /|____|    |___|
        \/
        `)

	lib.ServerLog("INFO", fmt.Sprintf("Server is available at %s", server.Addr))
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		if pathTemplate == "" {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			methods = []string{"ANY"}
		}
		lib.ServerLog("INFO", fmt.Sprintf("%-39s %s", pathTemplate, strings.Join(methods, ", ")))
		return nil
	})
}
