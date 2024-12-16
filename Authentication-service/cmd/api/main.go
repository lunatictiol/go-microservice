package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lunatictiol/go-authentication-service/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const WEBPORT = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")
	conn := connectDB()
	if conn == nil {
		log.Panic("cannot connect to db")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	serv := http.Server{
		Addr:    fmt.Sprintf(":%v", WEBPORT),
		Handler: app.routes(),
	}

	err := serv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("postgres is not ready")
			count++
		} else {
			log.Println("connected to db")
			return conn
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("backing connection try to db for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue

	}
}
