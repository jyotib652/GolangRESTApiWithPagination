package main

import (
	"database/sql"
	"fmt"

	// "log"
	"myRestAPIWithPagination/data"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	// Set up zerolog multiwriter
	// First, set log level for zerolog
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// Then, mention the file or path to the file to store logs
	logFile, _ := os.OpenFile(
		"restApiWithPagination.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	// Then, set up the multi logger with multiwriter. Here, one looger for file another for std output(console)
	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)
	// Then, include timestamp with the loggers
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()

	log.Info().Msg("Application is starting...")
	// log.Println("Starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic().Msg("Can't connect to Postgres!")
	}
	// Set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}

	log.Info().Msgf("Application is listenning on:http://localhost:%s", webPort)

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic().Msg(err.Error())
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

// this authentication service may start up before the database does, so
// we're creating connectToDB so that we can connect to the database even
// though database starts later than the authentication service. Here, in
// connectDB we'll try to establish connection after every 2 seconds
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Info().Msg("Postgres not yet ready...")
			counts++
		} else {
			log.Info().Msg("Connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Info().Msg(err.Error())
		}

		log.Info().Msg("Backing off for 2 seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
