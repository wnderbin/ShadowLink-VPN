package handlers

import (
	"database/sql"
	"log"
	"shadowlink/config"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/telebot.v4"
)

type TelegramHandler struct {
	Bot      *telebot.Bot
	Logger   *log.Logger
	Conf     *config.Config
	DB       *sql.DB
	RDB      *redis.Client
	ComDelay time.Duration
	VPNDelay time.Duration
}

func NewTelegramHandler(bot *telebot.Bot, logger *log.Logger, cfg *config.Config, db *sql.DB, rdb *redis.Client) *TelegramHandler {
	return &TelegramHandler{
		Bot:      bot,
		Logger:   logger,
		Conf:     cfg,
		DB:       db,
		RDB:      rdb,
		ComDelay: 5 * time.Second,
		VPNDelay: 5 * time.Minute,
	}
}

func (h *TelegramHandler) RegisterHandlers() { // Registers command and message handlers
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle("/help", h.HandleHelp)
	h.Bot.Handle("/vpn", h.HandleVPNConfig)
}
