package main

import (
	"context"
	"database/sql"
	"os"
	"sync"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mhakash/greenlight/internal/data"
	"github.com/mhakash/greenlight/internal/jsonlog"
	"github.com/mhakash/greenlight/internal/mailer"
)

const version = "1.0.0"

// visibility of fields needs to be public for env.Parse to work
type Config struct {
	Port int    `env:"PORT" envDefault:"4000"`
	Env  string `env:"ENV" envDefault:"development"`
	DB   struct {
		DSN          string `env:"DB_DSN"`
		MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
		MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
		MaxIdleTime  string `env:"DB_MAX_IDLE_TIME" envDefault:"15m"`
	}
	Limiter struct {
		Rps     float64 `env:"LIMITER_RPS" envDefault:"2"`
		Burst   int     `env:"LIMITER_BURST" envDefault:"4"`
		Enabled bool    `env:"LIMITER_ENABLED" envDefault:"true"`
	}
	SMTP struct {
		Host     string `env:"SMTP_HOST"`
		Port     int    `env:"SMTP_PORT"`
		Username string `env:"SMTP_USERNAME"`
		Password string `env:"SMTP_PASSWORD"`
		Sender   string `env:"SMTP_SENDER"`
	}
}

type application struct {
	config Config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelError)

	err := godotenv.Load()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDb(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDb(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
