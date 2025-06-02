package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

func (h *TelegramHandler) HandleVPNConfig(c telebot.Context) error {
	user := c.Sender()

	allowed, waitTime, err := h.checkRateLimitVPN(user.ID)
	if err != nil {
		h.Logger.Printf("[ ERROR ] redis error for user (vpn-delay): %d %s", user.ID, user.Username)
		return h.processVPN(c, user)
	}
	if !allowed {
		minutes := int(waitTime.Seconds()) / 60
		seconds := int(waitTime.Seconds()) % 60
		h.Logger.Printf("redis wait time for user: %d %s - %d minutes %d seconds", user.ID, user.Username, minutes, seconds)
		return c.Send("⏳ Пожалуйста, подождите %d минут и %d секунд", minutes, seconds)
	}
	return h.processVPN(c, user)
}

func (h *TelegramHandler) processVPN(c telebot.Context, user *telebot.User) error {
	var filename string
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := c.Notify(telebot.Typing); err != nil {
		h.Logger.Printf("[ ERROR ] Failed to send typing action %d %s: %v\n", user.ID, user.Username, err)
	}
	err := h.DB.QueryRowContext(ctx, "SELECT filename FROM wgconfigs WHERE user_id = $1", user.ID).Scan(&filename)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Printf("Creating a configuration for a user: %d %s", user.ID, user.Username)
			c.Send("🔧 Пользователь не зарегистрирован, создание конфигурации...")
			return h.generateVPNConfig(c, user)
		}
		h.Logger.Printf("[ ERROR ] QueryRowContext error: \"%s\" | User %d %s", err, user.ID, user.Username)
		return c.Send("❌ Ошибка базы данных")
	}
	h.Logger.Printf("Sending an existing configuration file to a user: %d %s", user.ID, user.Username)
	c.Send("✅ Пользователь зарегистрирован, отправка конфигурации...")
	return h.getVPNConfig(c, filename)
}

func (h *TelegramHandler) saveConfig(ctx context.Context, id int64, filepath string) error {
	_, err := h.DB.ExecContext(ctx,
		"INSERT INTO wgconfigs (user_id, filename) VALUES ($1, $2)", id, filepath)
	return err
}

func (h *TelegramHandler) getVPNConfig(c telebot.Context, configName string) error {
	file := &telebot.Document{
		File:     telebot.FromDisk(filepath.Join("wg-configs", configName)),
		FileName: configName,
		MIME:     "text/plain",
	}
	return c.Send(file)
}

func (h *TelegramHandler) generateVPNConfig(c telebot.Context, user *telebot.User) error {
	conf := h.Conf
	h.Logger.Printf("Generate vpn config from user: %d - %s\n", user.ID, user.Username)

	if err := os.MkdirAll("wg-configs", 0755); err != nil {
		h.Logger.Printf("[ ERROR ] create config directory error: %s", err)
		return c.Send("❌ Ошибка создания директории для конфигов")
	}

	wgConfigName := fmt.Sprintf("wg_%d.conf", user.ID)
	wgConfigPath := filepath.Join("wg-configs", wgConfigName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	h.saveConfig(ctx, user.ID, wgConfigName)

	privateKey, err := exec.Command("wg", "genkey").Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] private key generation error: %s\n", err)
		return c.Send("❌ Ошибка генерации приватного ключа")
	}
	privateKeyStr := strings.TrimSpace(string(privateKey))

	cmd := exec.Command("wg", "pubkey")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		h.Logger.Printf("[ ERROR ] creating stdin pipe: %s\n", err)
		return c.Send("❌ Ошибка создания потока данных")
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, privateKeyStr)
	}()
	publicKey, err := cmd.Output()
	if err != nil {
		h.Logger.Printf("[ ERROR ] error getting public key: %s\n", err)
		return c.Send("❌ Ошибка получения публичного ключа")
	}
	publicKeyStr := strings.TrimSpace(string(publicKey))

	clientIP := fmt.Sprintf("10.8.0.%d", (user.ID%253)+2)

	wgConfig := fmt.Sprintf("[Interface]\nPrivateKey = %s\nAddress = %s/24\nDNS = %s\n[Peer]\nPublicKey = %s\nEndpoint = %s:%d\nAllowedIPs = %s\nPersistentKeepalive = 25",
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

	cmd1 := exec.Command("wg", "set", conf.WGInterface,
		"peer", publicKeyStr,
		"allowed-ips", clientIP+"/32")

	if output, err := cmd1.CombinedOutput(); err != nil {
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
