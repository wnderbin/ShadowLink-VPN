package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/telebot.v4"
)

func (h *TelegramHandler) HandleGenerateVPNConfig(c telebot.Context) error {
	user := c.Sender()
	conf := h.Conf
	h.Logger.Printf("Generate vpn config from user: %d - %s\n", user.ID, user.Username)

	if err := os.MkdirAll("wg-configs", 0755); err != nil {
		h.Logger.Printf("[ ERROR ] create config directory error: %s", err)
		return c.Send("❌ Ошибка создания директории для конфигов")
	}

	wgConfigName := fmt.Sprintf("client_%d.conf", user.ID)
	wgConfigPath := filepath.Join("wg-configs", wgConfigName)

	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] private key generation error: %s\n", err)
		return c.Send("❌ Ошибка генерации приватного ключа")
	}
	privateKeyStr := strings.TrimSpace(string(privateKey))

	publicKey, err := exec.Command("wg", "pubkey", privateKeyStr).Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] public key generation error: %s\n", err)
		return c.Send("❌ Ошибка генерации публичного ключа")
	}
	publicKeyStr := strings.TrimSpace(string(publicKey))

	clientIP := fmt.Sprintf("10.8.0.%d", (user.ID%253)+2)

	wgConfig := fmt.Sprintf(`[Interface]
		PrivateKey = %s
		Address = %s/24
		DNS = %s

		[Peer]
		PublicKey = %s
		Endpoint = %s:%d
		AllowedIPs = %s
		PersistentKeepalive = 25`,

		privateKeyStr,
		clientIP,
		conf.DNS,
		conf.ServerPublicKey,
		conf.ServerPublicIP,
		conf.ServerPort,
		conf.AllowedIPs,
	)

	if err := os.WriteFile(wgConfigPath, []byte(wgConfig), 0600); err != nil {
		h.Logger.Printf("[ ERROR ] saving configuration error: %s", err)
		return c.Send("❌ Ошибка сохранения конфигурации")
	}

	cmd := exec.Command("wg", "set", conf.WGInterface,
		"peer", publicKeyStr,
		"allowed-ips", clientIP+"/32")

	if output, err := cmd.CombinedOutput(); err != nil {
		h.Logger.Printf("[ ERROR ] adding client error: %s, output: %s", err, string(output))
		return c.Send("❌ Ошибка добавления клиента в сервер")
	}

	file := &telebot.Document{
		File:     telebot.FromDisk(wgConfigPath),
		FileName: wgConfigName,
		MIME:     "text/plain",
	}

	h.Logger.Printf("New client added: IP %s, PublicKey %s\n", clientIP, publicKeyStr)

	return c.Send(file)
}
