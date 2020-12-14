package user

import (
	"context"
	"github.com/applicaset/auth-svc"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type service struct {
	repo          Repository
	authProviders map[string]AuthProvider
	authSvc       auth.Service
}

func (svc *service) GetUser(ctx context.Context, userUUID string) (*Entity, error) {
	entity, err := svc.repo.Find(ctx, userUUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on find user by uuid")
	}

	if entity == nil {
		return nil, ErrUserNotFound{UUID: userUUID}
	}

	return entity, nil
}

func (svc *service) Login(ctx context.Context, method string, args map[string]interface{}) (*LoginResponse, error) {
	authProvider, ok := svc.authProviders[method]
	if !ok {
		return nil, ErrInvalidAuthMethod{Name: method}
	}

	validate, err := authProvider.Validate(ctx, args)
	if err != nil {
		return nil, errors.Wrap(err, "error on validate auth data")
	}

	if validate.Validated() {
		return nil, ErrInvalidAuthData{}
	}

	entity, err := svc.repo.FindByAuthMethodAndID(ctx, method, validate.ID())
	if err != nil {
		return nil, errors.Wrap(err, "error on find user by auth method and id")
	}

	if entity == nil {
		entity = &Entity{
			UUID:         uuid.New().String(),
			RegisteredAt: time.Now(),
			AuthData: map[string]interface{}{
				method: validate,
			},
		}

		err = svc.repo.Create(ctx, *entity)
		if err != nil {
			return nil, errors.Wrap(err, "error on create new user")
		}
	} else {
		entity.AuthData[method] = validate
		err = svc.repo.Update(ctx, entity.UUID, *entity)
		if err != nil {
			return nil, errors.Wrap(err, "error on update user by uuid")
		}
	}

	atr, err := svc.authSvc.GenerateToken(ctx, entity.UUID)
	if err != nil {
		return nil, errors.Wrap(err, "error on generate token")
	}

	rsp := LoginResponse{
		UserUUID:     entity.UUID,
		AccessToken:  atr.AccessToken,
		ExpiresAt:    atr.ExpiresAt,
		RefreshToken: atr.RefreshToken,
	}

	return &rsp, nil
}

func (svc *service) Link(ctx context.Context, userUUID, method string, args map[string]interface{}) error {
	authProvider, ok := svc.authProviders[method]
	if !ok {
		return ErrInvalidAuthMethod{Name: method}
	}

	validate, err := authProvider.Validate(ctx, args)
	if err != nil {
		return errors.Wrap(err, "error on validate auth data")
	}

	if validate.Validated() {
		return ErrInvalidAuthData{}
	}

	entity, err := svc.repo.Find(ctx, userUUID)
	if err != nil {
		return errors.Wrap(err, "error on find user by uuid")
	}

	if entity == nil {
		return ErrUserNotFound{UUID: userUUID}
	}

	entity.AuthData[method] = validate

	err = svc.repo.Update(ctx, entity.UUID, *entity)
	if err != nil {
		return errors.Wrap(err, "error on update user by uuid")
	}

	return nil
}

func (svc *service) Unlink(ctx context.Context, userUUID, method string) error {
	entity, err := svc.repo.Find(ctx, userUUID)
	if err != nil {
		return errors.Wrap(err, "error on find user by uuid")
	}

	if entity == nil {
		return ErrUserNotFound{UUID: userUUID}
	}

	delete(entity.AuthData, method)

	err = svc.repo.Update(ctx, entity.UUID, *entity)
	if err != nil {
		return errors.Wrap(err, "error on update user by uuid")
	}

	return nil
}

func New(repo Repository, authSvc auth.Service, options ...Option) Service {
	svc := service{
		repo:          repo,
		authSvc:       authSvc,
		authProviders: make(map[string]AuthProvider),
	}

	for i := range options {
		options[i](&svc)
	}

	return &svc
}
