package main

import (
	"database/sql"
	"shadowlink/config"
	"shadowlink/handlers"
	"shadowlink/migrator"
	"shadowlink/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/telebot.v4"
)

func main() {
	cfg := config.Load()
	logger := utils.NewLogger(cfg.DebugMode)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "shadowlink.log", // Log file
		MaxSize:    50,               // MB
		MaxBackups: 10,               // Maximum files for storage
		MaxAge:     7,                // Maximum storage time
		Compress:   true,             // Compression of old logs
	})

	db, err := sql.Open("postgres", "postgres://superuser:123@localhost:5432/shadowlink?sslmode=disable")
	if err != nil {
		logger.Printf("postgres error: %s\n", err)
	}
	migrator.ApplyMigrations(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	botSettings := telebot.Settings{
		Token: cfg.BotApiKey,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	}
	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		logger.Fatalf("[ ERROR ] failed to create bot: %v", err)
	}
	tgHandler := handlers.NewTelegramHandler(bot, logger, cfg, db, redisClient)
	tgHandler.RegisterHandlers()
	logger.Println("starting bot...")
	bot.Start()
}
