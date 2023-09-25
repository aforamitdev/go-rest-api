package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"apexapi.aforamitdev.com/internal/data"
	_ "github.com/lib/pq"
	"golang.org/x/exp/slog"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const version = "1.0"

type Config struct {
	port int
	env  string
	db   struct {
		dns          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
}

type application struct {
	config Config
	logger *slog.Logger
	models data.Models
}

func main() {

	var cfg Config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "ENV (development|staging|production)")
	flag.StringVar(&cfg.db.dns, "db-dns", "postgres://amit:test@localhost/movie_data?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connection ")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSQL Max idle connection ")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "postgres max connections idel time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// open database,
	db, err := openDB(cfg)
	if err != nil {
		fmt.Println(err)
	}
	// defer database ...
	defer db.Close()

	// logger.Info("database connection pool established")

	// // database migrations driver ....
	// migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})

	// if err != nil {
	// 	logger.Error(err.Error())
	// 	os.Exit(1)
	// }
	// migrator, err := migrate.NewWithDatabaseInstance("file:///home/amit/go/src/apexapi/migrations", "postgres", migrationDriver)

	// if err != nil {
	// 	fmt.Println("error in NewWithDatabaseInstance \n ")
	// 	logger.Error(err.Error())
	// 	os.Exit(1)
	// }

	// err = migrator.Up()
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	os.Exit(1)
	// }
	// logger.Info("database migrations applied ")

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dns)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
