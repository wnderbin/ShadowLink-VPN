package handlers

import (
	"database/sql"
	"log"
	"shadowlink/config"

	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v4"
)

type TelegramHandler struct {
	Bot    *telebot.Bot
	Logger *log.Logger
	Conf   *config.Config
	DB     *sql.DB
	RDB    *redis.Client
}

func NewTelegramHandler(bot *telebot.Bot, logger *log.Logger, cfg *config.Config, db *sql.DB, rdb *redis.Client) *TelegramHandler {
	return &TelegramHandler{
		Bot:    bot,
		Logger: logger,
		Conf:   cfg,
		DB:     db,
		RDB:    rdb,
	}
}

func (h *TelegramHandler) RegisterHandlers() { // Registers command and message handlers
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle("/help", h.HandleHelp)
	h.Bot.Handle("/vpn", h.HandleGenerateVPNConfig)
}
