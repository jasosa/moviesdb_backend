package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jasosa/backend/models"
	"github.com/jasosa/backend/pkg/api"

	_ "github.com/lib/pq"
)

// type config struct {
// 	Port int
// 	Env  string
// 	Db   struct {
// 		dsn string
// 	}
// 	Jwt struct {
// 		secret string
// 	}
// }

// type Application struct {
// 	Config config
// 	Logger *log.Logger
// 	Models models.Models
// }

func main() {

	var apiPort int
	var apiEnv, dbDSN, secret string

	//tmp move secret to an env variable
	flag.IntVar(&apiPort, "port", 4000, "Server port to listen on.")
	flag.StringVar(&apiEnv, "env", "development", "Application environment (development | production) ")
	flag.StringVar(&dbDSN, "dsn", "postgres://postgres:t8sq1rF5@localhost/go_movies?sslmode=disable", "connection string")
	flag.StringVar(&secret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "jwt secret")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(dbDSN)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &api.ApiConfig{
		Port:   apiPort,
		Env:    apiEnv,
		Db:     &api.DBConfig{DSN: dbDSN},
		Jwt:    secret,
		Logger: logger,
		Models: models.NewModels(db),
	}

	api.New(app)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      api.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port ", app.Port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
