package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/avery-whitehead/germinal/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rickb777/date"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	db     internal.DBModel
}

func main() {
	db, err := connectWithConnector()
	if err != nil {
		fmt.Println(err)
		return
	}
	app := application{db: internal.NewDBModel(db)}
	repub, _ := app.toRepublican(date.New(2025, time.May, 26))
	fmt.Printf("%+v\n", repub)
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
