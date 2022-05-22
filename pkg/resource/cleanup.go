package resource

import (
	"context"
	"randgen-game/pkg/resource/internal/orm"
	"time"

	"gorm.io/gorm"
)

func CleanUpInBackground(ctx context.Context, db *gorm.DB) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(24 * time.Hour):
			db.Where("updated_at < CURRENT_TIMESTAMP - interval '24:00'").Delete(&orm.User{})
			db.Where("updated_at < CURRENT_TIMESTAMP - interval '24:00'").Delete(&orm.Room{})
			db.Where("updated_at IS NULL").Delete(&orm.User{})
			db.Where("updated_at IS NULL").Delete(&orm.Room{})
		}
	}
}
