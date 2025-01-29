package db

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/hrdemo/internal/config"

	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) *gorm.DB {
	dsn := cfg.DataBaseURL
	if strings.ToUpper(cfg.Env) == "LOCAL" {
		dsn += "?sslmode=disable&application_name=hr_api"
	} else {
		dsn += "?application_name=hr_api"
	}

	dbName, err := getDbName(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	oteldb, err := otelsql.Open("postgres", dsn, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
		semconv.DBNamespace(dbName),
	))
	if err != nil {
		log.Fatalln(err)
	}

	dbClient, err := gorm.Open(postgres.New(postgres.Config{
		Conn: oteldb,
	}))
	if err != nil {
		log.Fatalln(err)
	}

	sqlDb, err := dbClient.DB()
	if err != nil {
		log.Fatalln(err)
	}

	sqlDb.SetConnMaxIdleTime(time.Hour)
	sqlDb.SetMaxOpenConns(5)
	sqlDb.SetMaxIdleConns(5)

	if err := sqlDb.Ping(); err != nil {
		log.Fatalln(err)
	}

	return dbClient
}

func getDbName(dsn string) (string, error) {
	parsedURL, err := url.Parse(dsn)
	if err != nil {
		return "", err
	}

	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	return dbName, nil
}
