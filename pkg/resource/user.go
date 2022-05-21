package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"randgen-game/pkg/resource/internal/orm"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/lib/pq"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
)

type User struct {
	// Name is player name
	Name string `json:"name"`

	// Num is ITO Number (0-100) for a user.
	// -1 means that Num is not defined.
	Num int64 `json:"num"`

	// Open means that the player opens its number.
	Open bool `json:"open"`
}

func newUser(ormUser orm.User) User {
	return User{
		Name: ormUser.Name,
		Num:  ormUser.Num,
		Open: ormUser.Open,
	}
}

type Users []User

func (users Users) Find(name string) (user User, found bool) {
	for _, u := range users {
		if u.Name == name {
			return u, true
		}
	}
	return User{}, false
}

func AddUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info("creating user")
		roomID := c.Param("room_id")
		logger.Info(fmt.Sprintf("room id %s", roomID))
		if c.Request() == nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, "body is nil")
		}
		reqBodyStream := c.Request().Body
		defer reqBodyStream.Close()
		reqBodyBytes, _ := io.ReadAll(reqBodyStream)
		u := orm.User{}
		if err := json.Unmarshal(reqBodyBytes, &u); err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		logger.Info(fmt.Sprintf("user %+v", u))
		if uuid, err := uuid.Parse(roomID); err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, xerrors.Errorf("invalid room id format : %w", err).Error())
		} else {
			u.RoomID = uuid
		}
		if tx := db.Create(&u); tx.Error != nil {
			var pqErr *pq.Error
			if xerrors.As(tx.Error, &pqErr) {
				if pqErr.Code.Name() == "unique_violation" {
					return newErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("%s already exists", u.Name), string(pqErr.Code))
				}
				return newErrorResponse(c, http.StatusInternalServerError, xerrors.Errorf("failed to create new user : %w", tx.Error).Error(), string(pqErr.Code))
			}
			return newDefaultErrorResponse(c, http.StatusInternalServerError, xerrors.Errorf("failed to create new user : %w", tx.Error).Error())
		}
		logger.Info(fmt.Sprintf("created user %+v", u))
		db.Exec(fmt.Sprintf("NOTIFY \"%s\";", roomID))
		return c.JSON(http.StatusOK, User{Name: u.Name, Num: u.Num, Open: u.Open})
	}
}

func DeleteUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		name := c.Param("name")
		logger.Info("delete user " + name + " in room " + roomID)
		uuid, err := uuid.Parse(roomID)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, xerrors.Errorf("invalid room id format : %w", err).Error())
		}

		if tx := db.Delete(&orm.User{RoomID: uuid, Name: name}); tx.Error != nil {
			e := xerrors.Errorf("failed to delete users %s in room %s : %w", name, roomID, tx.Error)
			logger.Error(e.Error())
			return c.JSON(http.StatusInternalServerError, e.Error())
		}
		notify(db, uuid)
		return c.JSON(http.StatusOK, nil)
	}
}

func GetUsers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		logger.Info("getting users in room " + roomID)
		uuid, err := uuid.Parse(roomID)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, xerrors.Errorf("invalid room id format : %w", err).Error())
		}
		users, err := getUsers(db, uuid)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	}
}

func GetUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		name := c.Param("name")
		logger.Info("getting user " + name + " in room " + roomID)
		userORM := orm.User{Name: name}
		if uuid, err := uuid.Parse(roomID); err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, xerrors.Errorf("invalid room id format : %w", err).Error())
		} else {
			userORM.RoomID = uuid
		}
		if tx := db.First(&userORM); tx.Error != nil {
			e := xerrors.Errorf("failed to get users in room %s : %w", roomID, tx.Error)
			logger.Error(e.Error())
			return c.JSON(http.StatusInternalServerError, e.Error())
		}
		logger.Info(fmt.Sprintf("got user %+v", userORM))
		return c.JSON(http.StatusOK, User{Name: userORM.Name, Num: userORM.Num, Open: userORM.Open})
	}
}

func getUser(db *gorm.DB, roomID uuid.UUID, username string) (User, error) {
	userORM := orm.User{RoomID: roomID, Name: username}
	if tx := db.First(&userORM); tx.Error != nil {
		e := xerrors.Errorf("failed to get users in room %s : %w", roomID, tx.Error)
		logger.Error(e.Error())
		return User{}, e
	}
	logger.Info(fmt.Sprintf("got user %+v", userORM))
	return newUser(userORM), nil
}

func getUsers(db *gorm.DB, roomID uuid.UUID) (Users, error) {
	usersORM := []orm.User{}
	if tx := db.Order("created_at").Find(&usersORM, "room_id = ?", roomID); tx.Error != nil {
		e := xerrors.Errorf("failed to get users in room %s : %w", roomID, tx.Error)
		logger.Error(e.Error())
		return nil, e
	}
	users := Users{}
	for _, orm := range usersORM {
		users = append(users, newUser(orm))
	}
	return users, nil
}

func OpenCard(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomID := c.Param("room_id")
		username := c.Param("name")
		uuid, err := uuid.Parse(roomID)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusBadRequest, xerrors.Errorf("invalid room id format : %w", err).Error())
		}
		user, err := getUser(db, uuid, username)
		if err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		user.Open = true
		if err := db.Model(orm.User{RoomID: uuid, Name: username}).Select("open").Updates(orm.User{Open: true}).Error; err != nil {
			return newDefaultErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		notify(db, uuid)
		return c.JSON(http.StatusOK, nil)
	}
}

func notify(db *gorm.DB, roomID uuid.UUID) {
	db.Exec(fmt.Sprintf("NOTIFY \"%s\"", roomID))
}
