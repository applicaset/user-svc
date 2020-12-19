package usersvc

import "time"

type Entity struct {
	UUID         string
	RegisteredAt time.Time
	AuthData     map[string]interface{}
}
