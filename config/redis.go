package config

import (
	"context"
	"fmt"
	"github.com/Prototype-1/freelanceX_message.notification_service/pkg"
)

func IsUserOnline(ctx context.Context, userID string) (bool, error) {
	if pkg.Rdb == nil {
		return false, fmt.Errorf("redis client not initialized")
	}
	key := fmt.Sprintf("online:%s", userID)
	exists, err := pkg.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
