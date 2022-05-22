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

	"randgen-game/pkg/env"
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
	if !env.IsProd {
		e.Use(middleware.CORS())
	}

	e.Static("/", "build")
	e.Static("/rooms/:id", "build")
	e.GET("/api/v1/ws", wsHandler())
	e.POST("/api/v1/rooms", resource.AddRoom(db))
	e.GET("/api/v1/rooms", resource.GetRooms(db))
	e.GET("/api/v1/rooms/:id", resource.GetRoom(db))
	e.POST("/api/v1/rooms/:room_id/users", resource.AddUser(db))
	e.DELETE("/api/v1/rooms/:room_id/users/:name", resource.DeleteUser(db))
	e.GET("/api/v1/rooms/:room_id/users", resource.GetUsers(db))
	e.GET("/api/v1/rooms/:room_id/users/:name", resource.GetUser(db))
	e.POST("/api/v1/rooms/:room_id/start", resource.StartGame(db))
	e.POST("/api/v1/rooms/:room_id/end", resource.EndGame(db))
	e.POST("/api/v1/rooms/:room_id/users/:name/open", resource.OpenCard(db))

	ctx, cancel := context.WithCancel(context.Background())
	go resource.CleanUpInBackground(ctx, db)
	if err := e.Start(":" + os.Getenv("PORT")); err != nil {
		cancel()
		e.Logger.Fatal(err)
	}
}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

func createTables(db *gorm.DB) error {
	return resource.CreateTables(db)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
		ws.SetCloseHandler(func(code int, text string) error {
			logger.Info("websocket closed", zap.Int("code", code), zap.String("message", text))
			cancel()
			return nil
		})

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
			case <-time.After(5 * time.Minute):
				logger.Info(fmt.Sprintf("timeout room %s", roomID))
				return nil
			}
		}
	}
}

func eventCallback(event pq.ListenerEventType, err error) {
	logger.Info(fmt.Sprintf("listen: %d", event))
	logger.Info(fmt.Sprintf("err: %v", err))
}
