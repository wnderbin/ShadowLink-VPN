package handlers

import "gopkg.in/telebot.v4"

func (h *TelegramHandler) HandleStart(c telebot.Context) error {
	return c.Send(h.messageStart())
}

func (h *TelegramHandler) HandleHelp(c telebot.Context) error {
	return c.Send(h.messageHelp())
}
