package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"randgen-game/pkg/resource"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

func main() {
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening database: %q", err))
	}
	defer sqlDB.Close()
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening gorm database: %q", err))
	}
	if err := createTables(db); err != nil {
		logger.Fatal(err.Error())
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())

	e.Static("/", "build")
	e.GET("/ws", wsHandler())
	e.POST("/rooms", resource.AddRoom(db))
	e.GET("/rooms", resource.GetRooms(db))
	e.GET("/rooms/:id", resource.GetRoom(db))
	e.POST("/rooms/:room_id/users", resource.AddUser(db))
	e.GET("/rooms/:room_id/users", resource.GetUsers(db))
	e.GET("/rooms/:room_id/users/:name", resource.GetUser(db))
	e.POST("/rooms/:room_id/start", resource.StartGame(db))
	e.POST("/rooms/:room_id/users/:name/open", resource.OpenCard(db))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

func createTables(db *gorm.DB) error {
	return resource.CreateTables(db)
}

var upgrader = websocket.Upgrader{}

func wsHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		l := pq.NewListener(os.Getenv("DATABASE_URL"), 10*time.Second, 2*time.Minute, eventCallback)
		defer l.Close()

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()
		// Read
		_, roomID, err := ws.ReadMessage()
		if err != nil {
			logger.Error(err.Error())
		}
		logger.Info(fmt.Sprintf("received room id: %s", roomID))

		if err := l.Listen(string(roomID)); err != nil {
			return err
		}
		defer l.Unlisten(string(roomID))

		ctx, cancel := context.WithCancel(c.Request().Context())
		go func() {
			for {
				t, _, err := ws.ReadMessage()
				if err != nil || t == websocket.CloseMessage {
					cancel()
					return
				}
			}
		}()

		for {
			select {
			case <-l.Notify:
				logger.Info("notified")
				if err := ws.WriteMessage(websocket.TextMessage, []byte("updated")); err != nil {
					logger.Error(err.Error())
					return err
				}
			case <-ctx.Done():
				logger.Info(fmt.Sprintf("canceled room id: %s", roomID))
				return nil
			case <-time.After(1 * time.Minute):
				logger.Info(fmt.Sprintf("timeout room %s", roomID))
			}
		}
	}
}

func eventCallback(event pq.ListenerEventType, err error) {
	logger.Info(fmt.Sprintf("listen: %d", event))
	logger.Info(fmt.Sprintf("err: %v", err))
}
