package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/avery-whitehead/germinal/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	db     internal.DBModel
	logger *slog.Logger
}

func main() {
	db, err := connectWithConnector()
	if err != nil {
		fmt.Println(err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application{
		db:     internal.NewDBModel(db),
		logger: logger,
	}

	port, err := strconv.Atoi(os.Getenv("GM_PORT"))
	if port == 0 || err != nil {
		port = 8080
	}

	if err = app.serve(port); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func (app *application) serve(port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.routes(),
	}

	app.logger.Info("starting server", "addr", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func connectWithConnector() (*sql.DB, error) {
	var (
		dbUser                 = os.Getenv("DB_USER")
		dbPwd                  = os.Getenv("DB_PASS")
		dbName                 = os.Getenv("DB_NAME")
		instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	)

	dsn := fmt.Sprintf("user=%s password=%s database=%s", dbUser, dbPwd, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	opts := []cloudsqlconn.Option{cloudsqlconn.WithLazyRefresh()}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	return dbPool, nil
}
