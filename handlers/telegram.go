package handlers

import (
	"log"
	"shadowlink/config"

	"gopkg.in/telebot.v4"
)

type TelegramHandler struct {
	Bot    *telebot.Bot
	Logger *log.Logger
	Conf   *config.Config
}

func NewTelegramHandler(bot *telebot.Bot, logger *log.Logger, cfg *config.Config) *TelegramHandler {
	return &TelegramHandler{
		Bot:    bot,
		Logger: logger,
		Conf:   cfg,
	}
}

func (h *TelegramHandler) RegisterHandlers() { // Registers command and message handlers
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle("/help", h.HandleHelp)
	h.Bot.Handle("/vpn", h.HandleGenerateVPNConfig)
}
