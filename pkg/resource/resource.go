package resource

import (
	"randgen-game/pkg/resource/internal/orm"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func CreateTables(db *gorm.DB) error {
	logger.Info("Create tables")
	return orm.CreateTables(db)
}
