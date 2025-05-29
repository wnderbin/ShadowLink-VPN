package handlers

import (
	"log"

	"gopkg.in/telebot.v4"
)

type TelegramHandler struct {
	Bot    *telebot.Bot
	Logger *log.Logger
}

func NewTelegramHandler(bot *telebot.Bot, logger *log.Logger) *TelegramHandler {
	return &TelegramHandler{
		Bot:    bot,
		Logger: logger,
	}
}

func (h *TelegramHandler) RegisterHandlers() { // Registers command and message handlers
	h.Bot.Handle("/start", h.HandleStart)
	h.Bot.Handle("/help", h.HandleHelp)
}
