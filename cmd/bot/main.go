package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"talk/db/generated"
	"talk/internal/adapter/in/tg"
	pgAdapter "talk/internal/adapter/out/pg"
	tg2 "talk/internal/adapter/out/tg"
	"talk/internal/service"
	"talk/pkg/slogctx"

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

	sender := tg2.NewSender(b)
	updateService := service.NewUpdateService(log, store, sender)
	botAdapter := tg.NewBot(log, updateService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains, botAdapter.Handle)

	go b.Start(ctx)

	<-ctx.Done()
	log.Info("shutting down")
}
