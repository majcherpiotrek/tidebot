package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"tidebot/pkg/environment"
	"tidebot/pkg/jobs"
	notificationRepos "tidebot/pkg/notifications/repositories"
	"tidebot/pkg/ui/home"
	"tidebot/pkg/users/repositories"
	"tidebot/pkg/users/services"
	"tidebot/pkg/whatsapp"
	"tidebot/pkg/worldtides"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)

	e.Static("/assets", "assets")

	envFlagValue := flag.String("env", "", fmt.Sprintf("Environment ('%s' or '%s')", environment.EnvDevelopment, environment.EnvProduction))
	flag.Parse()

	env, err := environment.ParseEnvironment(*envFlagValue)

	if err != nil {
		envValue := os.Getenv("GO_ENV")
		env, err = environment.ParseEnvironment(envValue)
	}

	if err != nil {
		env = environment.EnvDevelopment
	}

	e.Logger.Infof("Environment set: %s", env)

	if env == environment.EnvDevelopment {
		e.Logger.SetLevel(log.DEBUG)
		e.Logger.Debug("Debug logging enabled for development environment")
	}

	envVars, err := env.LoadEnv()
	if err != nil {
		e.Logger.Fatalf("%v", err)
	}

	db, err := sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", envVars.TursoDbUrl, envVars.TursoDbAuthToken))
	if err != nil {
		e.Logger.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	err = runMigrations(db, e.Logger)
	if err != nil {
		e.Logger.Fatalf("Faied to run migrations: %v\n", err)
	}

	// Initialize repositories
	userRepository := repositories.NewUserRepository(db, e.Logger)
	notificationSubscriptionRepository := notificationRepos.NewNotificationSubscriptionRepository(db, e.Logger)

	// Initialize clients
	whatsappClient := whatsapp.NewWhatsappClient(envVars.TwilioWhatsAppFrom, e.Logger)
	worldTidesClient := worldtides.NewWorldTidesClient(envVars.WorldTidesApiKey, e.Logger)

	// Initialize services
	userService := services.NewUserService(userRepository, db, e.Logger)
	whatsappService := whatsapp.NewWhatsAppService(userService, notificationSubscriptionRepository, worldTidesClient, whatsappClient, e.Logger)
	jobsService := jobs.NewJobsService(userService, notificationSubscriptionRepository, whatsappService, worldTidesClient, e.Logger)

	// Initialize controllers
	jobsController := jobs.NewJobsController(jobsService, envVars.ApiKey, e.Logger)

	// Register routes
	whatsapp.RegisterWhatsappWebhook(e, whatsappService)
	whatsapp.RegisterComponents(e, envVars.TwilioWhatsAppFrom)
	jobsController.RegisterRoutes(e)

	home.RegisterHomeRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", envVars.ServerPort)))
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
