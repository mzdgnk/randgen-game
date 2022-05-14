package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"randgen-game/pkg/resource"
)

type Err struct {
	resource.ErrBody
	StatusCode int
}

func (e Err) Error() string {
	return fmt.Sprintf("http status %d %s", e.StatusCode, e.ErrBody)
}

func newError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	msg := resource.ErrBody{}
	json.Unmarshal(body, &msg)
	return Err{ErrBody: msg, StatusCode: resp.StatusCode}
}
