package main

import (
	"github.com/avery-whitehead/germinal/internal"
	_ "github.com/jackc/pgx"
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

}
