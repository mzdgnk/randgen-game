package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"randgen-game/pkg/resource"

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
	createTables(db)

	e := echo.New()

	e.GET("/", helloworld)
	e.POST("/rooms", resource.AddRoom(db))
	e.GET("/rooms/:id", resource.GetRoom(db))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func helloworld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

func createTables(db *gorm.DB) {
	resource.CreateRoomTable(db)
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
