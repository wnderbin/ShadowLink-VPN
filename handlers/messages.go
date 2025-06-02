package handlers

func (h *TelegramHandler) messageStart() string {
	return "<b>Приветствую, пользователь! 👋</b>\n\n"
}

func (h *TelegramHandler) messageHelp() string {
	return "<b>📑 Помощь</b>"
}
