package handlers

import "gopkg.in/telebot.v4"

func (h *TelegramHandler) HandleStart(c telebot.Context) error {
	menu := &telebot.ReplyMarkup{}
	btnGetVPN := menu.Data("Получить VPN", "get_vpn")
	btnInstructions := menu.Data("Инструкция", "instructions")

	menu.Inline(menu.Row(btnGetVPN), menu.Row(btnInstructions))

	return c.Send(h.messageStart(), menu)
}

func (h *TelegramHandler) HandleHelp(c telebot.Context) error {
	return c.Send(h.messageHelp())
}
