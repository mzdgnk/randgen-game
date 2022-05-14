package client

import (
	"encoding/json"
	"io"
	"net/http"
	"randgen-game/pkg/resource"

	"bytes"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewDevelopment()
}

type Interface interface {
	AddRoom() (*resource.Room, error)
	GetRoom(id uuid.UUID) (*resource.RoomWithUsers, error)
	GetRooms() ([]*resource.Room, error)
	StartGame(roomID uuid.UUID, topic string) error
	AddUser(roomID uuid.UUID, name string) (*resource.User, error)
	GetUsers(roomID uuid.UUID) (resource.Users, error)
	GetUser(roomID uuid.UUID, name string) (*resource.User, error)
	OpenCard(roomID uuid.UUID, username string) error
}

type client struct {
	baseURL string
}

func New(baseURL string) Interface {
	return &client{baseURL: baseURL}
}

func (c *client) AddRoom() (*resource.Room, error) {
	url := c.baseURL + "/" + "rooms"
	httpClient := new(http.Client)
	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to create room")
		logger.Error(e.Error())
		return nil, e
	}
	body, _ := io.ReadAll(resp.Body)
	room := &resource.Room{}
	if err := json.Unmarshal(body, room); err != nil {
		return nil, err
	}
	return room, nil
}

func (c *client) GetRoom(id uuid.UUID) (*resource.RoomWithUsers, error) {
	url := c.baseURL + "/" + "rooms" + "/" + id.String()
	httpClient := new(http.Client)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, _ := io.ReadAll(resp.Body)
	room := &resource.RoomWithUsers{}
	if err := json.Unmarshal(body, room); err != nil {
		e := errors.Wrap(newError(resp), "failed to create room")
		logger.Error(e.Error())
		return nil, e
	}
	return room, nil
}

func (c *client) GetRooms() ([]*resource.Room, error) {
	url := c.baseURL + "/rooms"
	httpClient := new(http.Client)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to create room")
		logger.Error(e.Error())
		return nil, e
	}
	body, _ := io.ReadAll(resp.Body)
	room := []*resource.Room{}
	if err := json.Unmarshal(body, &room); err != nil {
		return nil, err
	}
	return room, nil
}

func (c *client) StartGame(roomID uuid.UUID, topic string) error {
	url := c.baseURL + "/rooms/" + roomID.String() + "/start"
	b, _ := json.Marshal(resource.StartGameRequest{Topic: topic})
	httpClient := new(http.Client)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to start game")
		logger.Error(e.Error())
		return e
	}
	return nil
}
