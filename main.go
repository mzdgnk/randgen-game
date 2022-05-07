package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {
	sqlDB, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening gorm database: %q", err)
	}

	e := echo.New()

	createTables(c, db)

	e.GET("/", helloworld)
	e.POST("/rooms", addRoom(db))
	e.GET("/rooms/:id", getRoom(db))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

type room struct {
	ID string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
}

func addRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := room{}
		db.Create(&r)
		return c.JSON(http.StatusOK, r)
	}
}

func getRoom(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		r := room{}
		db.First(&r, "id = ?", id)
		return c.JSON(http.StatusOK, r)
	}
}

func createTables(c echo.Context, db *gorm.DB) {
	db.AutoMigrate(&room{})
	// if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS room (
	// 	id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
	// 	)`); err != nil {
	// 	c.String(http.StatusInternalServerError,
	// 		fmt.Sprintf("Error creating database table: %q", err))
	// 	return
	// }

	// if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS session (
	// 	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	// 	FOREIGN KEY (room_id) REFERENCES room(id)
	// 	)`); err != nil {
	// 	c.String(http.StatusInternalServerError,
	// 		fmt.Sprintf("Error creating database table: %q", err))
	// 	return
	// }

	// if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS user (
	// 	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	// 	FOREIGN KEY (session_id) REFERENCES session(id),
	// 	name varchar
	// 	)`); err != nil {
	// 	c.String(http.StatusInternalServerError,
	// 		fmt.Sprintf("Error creating database table: %q", err))
	// 	return
	// }

	// if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS number (
	// 	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	// 	FOREIGN KEY (session_id) REFERENCES session(id),
	// 	FOREIGN KEY (user_id) REFERENCES user(id),
	// 	num int
	// 	)`); err != nil {
	// 	c.String(http.StatusInternalServerError,
	// 		fmt.Sprintf("Error creating database table: %q", err))
	// 	return
	// }
}
