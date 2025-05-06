package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tanvir-rifat007/gymBuddy/internal/agents"
	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)

type config struct {
	Port string
	Env  string
	db   struct {
		dsn string
	}
}

type application struct {
	logger *slog.Logger
	sentry *sentry.Hub
	config config
	models data.Models
	webauthn *webauthn.WebAuthn
	ai       *agents.OpenAPI

}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load .env
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", "error", err)
	}

	// Sentry setup
	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		logger.Error("Error initializing Sentry", "error", err)
	}
	defer sentry.Flush(2 * time.Second)

	var cfg config
	port := os.Getenv("PORT")
if port == "" {
	port = "4000" // fallback for local dev
}

	flag.StringVar(&cfg.Port, "port", port, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("DB_URL"), "Database connection string")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Connected to database", "dsn", cfg.db.dsn)



	wconfig := &webauthn.Config{
	RPDisplayName: "GymTrackr",
	RPID:          "gym-trackr-production.up.railway.app",
	RPOrigins:     []string{"https://gym-trackr-production.up.railway.app"},
}

webAuthn, err := webauthn.New(wconfig)
if err != nil {
	logger.Error("Error initializing WebAuthn", err)
	os.Exit(1)
}

	app := &application{
		logger: logger,
		config: cfg,
		sentry: sentry.CurrentHub(),
		models: data.NewModels(db),
		webauthn: webAuthn,
		ai : agents.NewOpenAI(context.Background(), os.Getenv("OPENAI_API_KEY"), "gpt-4o-mini",nil),
	}

	// // Setup cron
	// loc, _ := time.LoadLocation("Asia/Dhaka")
	// c := cron.New(cron.WithLocation(loc))

	// // Log current time in Dhaka
	// logger.Info("Current Asia/Dhaka time", "now", time.Now().In(loc).Format("15:04:05"))

	// _, err = c.AddFunc("* * * * *", func() {
	// 	defer func() {
	// 		if r := recover(); r != nil {
	// 			logger.Error("Cron panic", "error", r)
	// 		}
	// 	}()
	// 	logger.Info("✅ Every minute cron job running")
	// 	app.SendWorkoutReminderEmails(app.models.Users)
	// })

	// if err != nil {
	// 	logger.Error("Error scheduling cron job", "error", err)
	// 	os.Exit(1)
	// }
	// c.Start()
	// logger.Info("Cron job started (every minute)")


	// server

	err = app.serve(":" + port)

	if err != nil {
		logger.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
