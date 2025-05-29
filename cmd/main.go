package main

import (
	"shadowlink/config"
	"shadowlink/handlers"
	"shadowlink/utils"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/telebot.v4"
)

func main() {
	cfg := config.Load()
	logger := utils.NewLogger(cfg.DebugMode)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "shadowlink.log", // Log file
		MaxSize:    100,              // MB
		MaxBackups: 10,               // Maximum files for storage
		MaxAge:     7,                // Maximum storage time
		Compress:   true,             // Compression of old logs
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
	tgHandler := handlers.NewTelegramHandler(bot, logger)
	tgHandler.RegisterHandlers()
	logger.Println("starting bot...")
	bot.Start()
}
