package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladimish/talk/db/generated"
	"github.com/vladimish/talk/internal/adapter/in/tg"
	"github.com/vladimish/talk/internal/adapter/out/openai"
	pgAdapter "github.com/vladimish/talk/internal/adapter/out/pg"
	redisAdapter "github.com/vladimish/talk/internal/adapter/out/redis"
	"github.com/vladimish/talk/internal/adapter/out/telegramify"
	tgAdapter "github.com/vladimish/talk/internal/adapter/out/tg"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/pkg/slogctx"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logHandler := slogctx.NewContextHandler(textHandler)
	log := slog.New(logHandler)
	log.InfoContext(ctx, "starting bot")

	pg, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("can't connect to database", "error", err)
		//nolint:gocritic
		os.Exit(1)
	}

	if err = runMigrations(ctx, log, pg.DB); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	queries := generated.New(pg)
	store := pgAdapter.NewPg(queries)

	b, err := bot.New(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Error("OPENAI_API_KEY environment variable is required")
		os.Exit(1)
	}
	completion := openai.NewOpenAICompletion(openAIKey)

	telegramifyURL := os.Getenv("TELEGRAMIFY_URL")
	if telegramifyURL == "" {
		telegramifyURL = "http://localhost:8000"
		log.Warn("TELEGRAMIFY_URL not set, using default", "url", telegramifyURL)
	}
	formatter := telegramify.New(telegramifyURL)

	// Initialize Redis queue
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
		log.Warn("REDIS_URL not set, using default", "url", redisURL)
	}

	redisQueue, err := redisAdapter.NewRedisQueue(redisURL)
	if err != nil {
		log.Error("failed to initialize Redis queue", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := redisQueue.Close(); closeErr != nil {
			log.Error("failed to close Redis connection", "error", closeErr)
		}
	}()

	sender := tgAdapter.NewSender(b, formatter, log)
	updateService := service.NewUpdateService(log, store, sender, completion, redisQueue)
	botAdapter := tg.NewBot(log, updateService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, botAdapter.Handle)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, botAdapter.HandleCallback)

	// Use a general update handler for payment events
	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.PreCheckoutQuery != nil
	}, botAdapter.HandlePreCheckoutQuery)

	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.Message != nil && update.Message.SuccessfulPayment != nil
	}, botAdapter.HandleSuccessfulPayment)

	go b.Start(ctx)

	<-ctx.Done()
	log.Info("shutting down")
}

func runMigrations(ctx context.Context, log *slog.Logger, db *sql.DB) error {
	log.InfoContext(ctx, "running database migrations")

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "db/migrations"); err != nil {
		return err
	}

	log.InfoContext(ctx, "database migrations completed successfully")
	return nil
}
