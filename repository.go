package usersvc

import "context"

type Repository interface {
	Create(ctx context.Context, entity Entity) (err error)
	Find(ctx context.Context, uuid string) (res *Entity, err error)
	FindByAuthMethodAndID(ctx context.Context, method, id string) (res *Entity, err error)
	Update(ctx context.Context, uuid string, entity Entity) (err error)
}
