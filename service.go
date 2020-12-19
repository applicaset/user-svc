package usersvc

import (
	"context"
	"time"
)

type Service interface {
	GetUser(ctx context.Context, userUUID string) (res *Entity, err error)
	Login(ctx context.Context, method string, args map[string]interface{}) (res *LoginResponse, err error)
	Link(ctx context.Context, userUUID, method string, args map[string]interface{}) (err error)
	Unlink(ctx context.Context, userUUID, method string) (err error)
}

type LoginResponse struct {
	UserUUID     string
	AccessToken  string
	ExpiresAt    *time.Time
	RefreshToken *string
}
