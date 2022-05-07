package resource

import (
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Room struct {
	ID string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
}

func CreateRoomTable(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&Room{})
}

func AddRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := Room{}
		db.Create(&r)
		return c.JSON(http.StatusOK, r)
	}
}

func GetRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		r := Room{}
		db.First(&r, "id = ?", id)
		return c.JSON(http.StatusOK, r)
	}
}
