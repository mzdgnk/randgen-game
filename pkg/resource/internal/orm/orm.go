package orm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTables(db *gorm.DB) error {
	if tx := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`); tx.Error != nil {
		return tx.Error
	}
	if err := db.AutoMigrate(&Room{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

type Room struct {
	ID      uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Started bool      `gorm:"default:false" json:"started"`
	Topic   string    `json:"topic"`
}

type User struct {
	Name      string    `gorm:"primaryKey" json:"name"`
	RoomID    uuid.UUID `gorm:"primaryKey" json:"-"`
	Room      Room      `json:"room"`
	Num       int64     `gorm:"default:-1"`
	Open      bool      `gorm:"default:false"`
	CreatedAt time.Time
}
