package handlers

import (
	"context"
	"fmt"
	"time"
)

func (h *TelegramHandler) checkRateLimitCommand(userID int64) (allowed bool, remaining time.Duration, err error) {
	ctx := context.Background()
	key := fmt.Sprintf("command_rate_limit:%d", userID)

	set, err := h.RDB.SetNX(ctx, key, "1", h.ComDelay).Result()
	if err != nil {
		return true, 0, err
	}
	if set {
		return true, 0, nil
	}
	ttl, err := h.RDB.TTL(ctx, key).Result()
	if err != nil {
		return true, 0, err
	}
	if ttl < 0 {
		return true, 0, nil
	}
	return false, ttl, nil
}

func (h *TelegramHandler) checkRateLimitVPN(userID int64) (allowed bool, remaining time.Duration, err error) {
	ctx := context.Background()
	key := fmt.Sprintf("vpn_rate_limit:%d", userID)

	set, err := h.RDB.SetNX(ctx, key, "1", h.VPNDelay).Result()
	if err != nil {
		return true, 0, err
	}

	if set {
		return true, 0, nil
	}

	ttl, err := h.RDB.TTL(ctx, key).Result()
	if err != nil {
		return true, 0, err
	}
	if ttl < 0 {
		return true, 0, nil
	}

	return false, ttl, nil
}
