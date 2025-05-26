package main

import _ "github.com/jackc/pgx"

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
}

func main() {

}
