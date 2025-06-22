package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"tidebot/pkg/whatsapp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var (
	TURSO_DB_URL         = ""
	TURSO_DB_AUTH_TOKEN  = ""
	TWILIO_WHATSAPP_FROM = ""
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Static("/assets", "assets")

	envFlagValue := flag.String("env", "", fmt.Sprintf("Environment ('%s' or '%s')", EnvDevelopment, EnvProduction))
	flag.Parse()

	env, err := ParseEnvironment(*envFlagValue)

	if err != nil {
		envValue := os.Getenv("GO_ENV")
		env, err = ParseEnvironment(envValue)
	}

	if err != nil {
		env = EnvDevelopment
	}

	err = env.LoadEnv()
	if err != nil {
		e.Logger.Fatalf("%v", err)
	}

	// Assign environment variables to the corresponding variables
	TURSO_DB_URL = os.Getenv("TURSO_DB_URL")
	TURSO_DB_AUTH_TOKEN = os.Getenv("TURSO_DB_AUTH_TOKEN")
	TWILIO_WHATSAPP_FROM = os.Getenv("TWILIO_WHATSAPP_FROM")

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", TURSO_DB_URL, TURSO_DB_AUTH_TOKEN))
	if err != nil {
		e.Logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	err = runMigrations(db, e.Logger)
	if err != nil {
		e.Logger.Fatalf("Faied to run migrations: %v\n", err)
	}

	whatsappClient := whatsapp.NewWhatsappClient(TWILIO_WHATSAPP_FROM, e.Logger)

	whatsappClient.SendMessage("Hello from code", "+34608368242")

	e.Logger.Fatal(e.Start(":42069"))
}

func runMigrations(db *sql.DB, log echo.Logger) error {
	log.Infof("Database migration started")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		log.Fatal("Failed to create sqlite3 driver instance: %v\n", err)
	}

	migrationTool, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3", driver)

	if err != nil {
		log.Fatal("Failed to create migration tool instance: %v\n", err)
	}

	log.Infof("Migration client created")

	// Run migrations
	err = migrationTool.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("An error occurred while running the migrations: %v", err)
	}

	if err == migrate.ErrNoChange {
		log.Infof("Database is up to date")
	}
	if err == nil {
		log.Infof("Migrations successfully applied")
	}

	return nil
}
