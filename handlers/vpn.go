package handlers

import (
	"fmt"
	"os/exec"
	"strings"

	"gopkg.in/telebot.v4"
)

func (h *TelegramHandler) HandleGenerateVPNConfig(c telebot.Context) error {
	user := c.Sender()
	conf := h.Conf
	h.Logger.Printf("Generate vpn config from user: %d - %s\n", user.ID, user.Username)

	configName := fmt.Sprintf("client_%d.conf", user.ID)
	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] private key generation error: %s\n", err)
		return c.Send("❌ Ошибка генерации приватного ключа")
	}

	publicKey, err := exec.Command("wg", "pubkey", string(privateKey)).Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] public key generation error: %s\n", err)
		return c.Send("❌ Ошибка генерации публичного ключа")
	}

	wgConfig := fmt.Sprintf("[Interface]\nPrivateKey = %s\nAddress = 10.0.0.%d/24\nDNS = %s\n\n[Peer]\nPublicKey = %s\nEndpoint = %s:%d\nAllowedIPs = %s\nPersistentKeepalive = 25",
		strings.TrimSpace(string(privateKey)),
		user.ID%254+1,
		conf.DNS,
		conf.ServerPublicKey,
		conf.ServerPublicIP,
		conf.ServerPort,
		conf.AllowedIPs,
	)

	wgConfigPath := fmt.Sprintf("%s/%s", "wg-configs", configName)
	err = exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > %s", wgConfig, wgConfigPath)).Run()
	if err != nil {
		h.Logger.Printf("[ ERROR ] saving configuration error: %s", err)
		return c.Send("❌ Ошибка сохранения конфигурации")
	}
	err = exec.Command("wg", "set", conf.WGInterface, "peer", strings.TrimSpace(string(publicKey)), "allowed-ips", fmt.Sprintf("10.0.0.%d/32", user.ID%254+1)).Run()
	if err != nil {
		h.Logger.Printf("[ ERROR ] adding client error: %s", err)
		return c.Send("❌ Ошибка доббавления клиента")
	}
	file := &telebot.Document{File: telebot.FromDisk(wgConfigPath), FileName: configName}
	return c.Send(file)
}
