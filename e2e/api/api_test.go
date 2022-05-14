package api

import (
	"errors"
	"randgen-game/pkg/client"
	"randgen-game/pkg/resource"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var host = "localhost:5000"
var baseURL = "http://" + host
var websockURL = "ws://" + host + "/ws"

func TestSmokeCreateUser(t *testing.T) {
	c := client.New(baseURL)
	added, err := c.AddRoom()
	require.Nil(t, err)
	room, err := c.GetRoom(added.ID)
	require.Nil(t, err)
	t.Logf("%+v", room)
	assert.Equal(t, *added, room.Room)
	rooms, err := c.GetRooms()
	require.Nil(t, err)
	t.Log("print rooms")
	for _, r := range rooms {
		t.Logf("%+v", r)
	}
	t.Log("--------")
	assert.Greater(t, len(rooms), 0)
	newUser, err := c.AddUser(room.ID, "user01")
	require.Nil(t, err)
	require.NotNil(t, newUser)
	assert.Equal(t, resource.User{Name: "user01", Num: -1}, *newUser)
	_, err = c.GetUsers(room.ID)
	require.Nil(t, err)
	gotUser, err := c.GetUser(room.ID, "user01")
	require.Nil(t, err)
	assert.Equal(t, resource.User{Name: "user01", Num: -1}, *gotUser)
}

func TestSmokeStartGame(t *testing.T) {
	c := client.New(baseURL)
	// Create Room
	room, err := c.AddRoom()
	require.Nil(t, err)
	t.Logf("%+v", room)

	// Create Users
	user01, err := c.AddUser(room.ID, "user01")
	require.Nil(t, err)
	require.NotNil(t, user01)
	assert.Equal(t, resource.User{Name: "user01", Num: -1}, *user01)
	user02, err := c.AddUser(room.ID, "user02")
	require.Nil(t, err)

	// Start Game
	err = c.StartGame(room.ID, "topic01")
	require.Nil(t, err)
	updatedRoom, err := c.GetRoom(room.ID)
	require.Nil(t, err)
	assert.True(t, updatedRoom.Started)
	updatedUser01, err := c.GetUser(room.ID, user01.Name)
	require.Nil(t, err)
	assert.Greater(t, updatedUser01.Num, int64(-1))
	updatedUser02, err := c.GetUser(room.ID, user02.Name)
	require.Nil(t, err)
	assert.Greater(t, updatedUser02.Num, int64(-1))

	// Open card
	err = c.OpenCard(room.ID, user01.Name)
	require.Nil(t, err)
	user01, err = c.GetUser(room.ID, user01.Name)
	require.Nil(t, err)
	assert.True(t, user01.Open)

	// Validate
	user02, err = c.GetUser(room.ID, user02.Name)
	require.Nil(t, err)
	assert.False(t, user02.Open)

	r, err := c.GetRoom(room.ID)
	require.Nil(t, err)
	assert.Equal(t, resource.Room{ID: room.ID, Started: true, Topic: "topic01"}, r.Room)
	u, found := r.Users.Find("user01")
	require.True(t, found)
	assert.Equal(t, resource.User{Name: "user01", Num: u.Num, Open: true}, u)
	u, found = r.Users.Find("user02")
	require.True(t, found)
	assert.Equal(t, resource.User{Name: "user02", Num: u.Num, Open: false}, u)
}

func TestWebsocket(t *testing.T) {
	c := client.New(baseURL)
	added, err := c.AddRoom()
	require.Nil(t, err)
	t.Log(added)
	ws, _, err := websocket.DefaultDialer.Dial(websockURL, nil)
	defer func() {
		ws.WriteMessage(websocket.CloseMessage, nil)
		ws.Close()
	}()
	require.Nil(t, err)
	ws.WriteMessage(websocket.TextMessage, []byte(added.ID.String()))
	go func() {
		time.Sleep(time.Second)
		c.AddUser(added.ID, "user01")
	}()
	go func() {
		time.Sleep(2 * time.Second)
		c.AddUser(added.ID, "user02")
	}()
	_, msg, err := ws.ReadMessage()
	require.Nil(t, err)
	assert.Equal(t, "updated", string(msg))
	_, msg, err = ws.ReadMessage()
	require.Nil(t, err)
	assert.Equal(t, "updated", string(msg))
}

func TestCreateUsersWithSameName(t *testing.T) {
	c := client.New(baseURL)
	room, err := c.AddRoom()
	require.Nil(t, err)
	t.Logf("%+v", room)
	user1, err := c.AddUser(room.ID, "user01")
	require.Nil(t, err)
	require.NotNil(t, user1)
	assert.Equal(t, resource.User{Name: "user01", Num: -1}, *user1)
	user2, err := c.AddUser(room.ID, "user01")
	require.NotNil(t, err)
	t.Log(err)
	restErr := client.Err{}
	require.True(t, errors.As(err, &restErr))
	assert.Equal(t, "23505", restErr.ErrID)
	require.Nil(t, user2)
}
