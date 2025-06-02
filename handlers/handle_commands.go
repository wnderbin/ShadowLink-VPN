package handlers

import (
	"fmt"

	"gopkg.in/telebot.v4"
)

func (h *TelegramHandler) HandleStart(c telebot.Context) error {
	user := c.Sender()
	allowed, waitTime, err := h.checkRateLimitCommand(user.ID)
	if err != nil {
		h.Logger.Printf("[ ERROR ] redis error for user (command-delay): %d %s", user.ID, user.Username)
		return c.Send(h.messageStart(), telebot.ModeHTML)
	}
	if !allowed {
		h.Logger.Printf("redis wait time for user: %d %s - %d seconds", user.ID, user.Username, int(waitTime.Seconds()))
		return c.Send(fmt.Sprintf("⏳ Пожалуйста, подождите %d секунд перед отправкой следующего запроса", int(waitTime.Seconds())))
	}
	return c.Send(h.messageStart(), telebot.ModeHTML)
}

func (h *TelegramHandler) HandleHelp(c telebot.Context) error {
	user := c.Sender()
	allowed, waitTime, err := h.checkRateLimitCommand(user.ID)
	if err != nil {
		h.Logger.Printf("[ ERROR ] redis error for user (command-delay): %d %s", user.ID, user.Username)
		return c.Send(h.messageHelp(), telebot.ModeHTML)
	}
	if !allowed {
		h.Logger.Printf("redis waite time for user: %d %s - %d seconds", user.ID, user.Username, int(waitTime.Seconds()))
		return c.Send(fmt.Sprintf("⏳ Пожалуйста, подождите %d секунд перед отправкой следующего запроса", int(waitTime.Seconds())))
	}
	return c.Send(h.messageHelp(), telebot.ModeHTML)
}
