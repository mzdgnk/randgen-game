package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"randgen-game/pkg/resource"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"
)

func (c *client) AddUser(roomID uuid.UUID, name string) (*resource.User, error) {
	reqBody, _ := json.Marshal(map[string]string{"name": name})

	url := c.baseURL + "/rooms/" + roomID.String() + "/users/"
	httpClient := new(http.Client)
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("fialed to create room : %w", newError(resp))
	}
	body, _ := io.ReadAll(resp.Body)
	user := &resource.User{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *client) PutUser(roomID uuid.UUID, name string) (*resource.User, error) {
	reqBody, _ := json.Marshal(map[string]int{"num": -1})
	url := c.baseURL + "/rooms/" + roomID.String() + "/users/" + name
	httpClient := new(http.Client)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("fialed to create room : %w", newError(resp))
	}
	body, _ := io.ReadAll(resp.Body)
	user := &resource.User{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *client) GetUsers(roomID uuid.UUID) (resource.Users, error) {
	url := c.baseURL + "/rooms/" + roomID.String() + "/users"
	httpClient := new(http.Client)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to get users")
		logger.Error(e.Error())
		return nil, e
	}
	body, _ := io.ReadAll(resp.Body)
	users := resource.Users{}
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *client) GetUser(roomID uuid.UUID, name string) (*resource.User, error) {
	url := c.baseURL + "/rooms/" + roomID.String() + "/users/" + name
	httpClient := new(http.Client)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to get user "+name)
		logger.Error(e.Error())
		return nil, e
	}
	body, _ := io.ReadAll(resp.Body)
	user := &resource.User{}
	logger.Info(fmt.Sprintf("got user %s", body))
	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *client) OpenCard(roomID uuid.UUID, username string) error {
	url := c.baseURL + "/rooms/" + roomID.String() + "/users/" + username + "/open"
	httpClient := new(http.Client)
	resp, err := httpClient.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		e := errors.Wrap(newError(resp), "failed to open card")
		logger.Error(e.Error())
		return e
	}
	return nil
}
