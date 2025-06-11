package main

import (
	"context"
	"github.com/vladimish/talk/db/generated"
	"github.com/vladimish/talk/internal/adapter/in/tg"
	"github.com/vladimish/talk/internal/adapter/out/openai"
	pgAdapter "github.com/vladimish/talk/internal/adapter/out/pg"
	"github.com/vladimish/talk/internal/adapter/out/telegramify"
	tgAdapter "github.com/vladimish/talk/internal/adapter/out/tg"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/pkg/slogctx"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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

	sender := tgAdapter.NewSender(b, formatter, log)
	updateService := service.NewUpdateService(log, store, sender, completion)
	botAdapter := tg.NewBot(log, updateService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, botAdapter.Handle)

	go b.Start(ctx)

	<-ctx.Done()
	log.Info("shutting down")
}
