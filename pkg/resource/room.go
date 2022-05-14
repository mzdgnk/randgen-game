package resource

import (
	"fmt"
	"math/rand"
	"net/http"
	"randgen-game/pkg/resource/internal/orm"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Room orm.Room

type RoomWithUsers struct {
	Room
	Users Users `json:"users"`
}

func AddRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info("creating room")
		r := Room{}
		if tx := db.Create(&r); tx.Error != nil {
			e := errors.Wrap(tx.Error, "failed to create room")
			logger.Error(e.Error())
			return c.JSON(http.StatusInternalServerError, e.Error())
		}
		logger.Info(fmt.Sprintf("created room: %+v", r))
		return c.JSON(http.StatusOK, r)
	}
}

func updateRoom(db *gorm.DB, r Room) error {
	logger.Info("updating room")
	if tx := db.Save(&r); tx.Error != nil {
		return tx.Error
	}
	return nil
}

func GetRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")
		logger.Info("getting room " + idStr)
		id, err := uuid.Parse(idStr)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		r, err := getRoom(db, id)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		users, err := getUsers(db, r.ID)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, RoomWithUsers{Room: r, Users: users})
	}
}

func getRoom(db *gorm.DB, id uuid.UUID) (Room, error) {
	r := Room{}
	if tx := db.First(&r, "id = ?", id.String()); tx.Error != nil {
		e := errors.Wrap(tx.Error, "failed to get room "+id.String())
		logger.Error(e.Error())
		return Room{}, e
	}
	return r, nil
}

func GetRooms(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info("getting rooms")
		rooms := []Room{}
		if tx := db.Find(&rooms); tx.Error != nil {
			e := errors.Wrap(tx.Error, "failed to get rooms")
			logger.Error(e.Error())
			return c.JSON(http.StatusInternalServerError, e.Error())
		}
		return c.JSON(http.StatusOK, rooms)
	}
}

type StartGameRequest struct {
	Topic string `json:"topic"`
}

func StartGame(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		uuid, err := uuid.Parse(roomID)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid room id %s", roomID))
		}

		reqBody := StartGameRequest{}
		if err := c.Bind(&reqBody); err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, "invalid body")
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			users, err := getUsers(db, uuid)
			if err != nil {
				return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
			}
			if err := updateRoom(db, Room{ID: uuid, Started: true, Topic: reqBody.Topic}); err != nil {
				return newDefaultErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("failed to update room %s", roomID))
			}
			rands := generateRand(len(users))
			for i, _ := range users {
				users[i].Num = rands[i]
				users[i].Open = false
			}
			logger.Info(fmt.Sprintf("users: %+v", users))
			for _, u := range users {
				if err := updateUser(db, uuid, u); err != nil {
					return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		notify(db, uuid)
		return nil
	}
}

func generateRand(userNum int) []int64 {
	r := []int64{}
	for len(r) < userNum {
		i := int64(rand.Intn(101))
		if !contain(r, i) {
			r = append(r, i)
		}
	}
	return r
}

func contain(s []int64, e int64) bool {
	for _, i := range s {
		if i == e {
			return true
		}
	}
	return false
}
