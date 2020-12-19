package usersvc

import "context"

type Option func(*service)

func WithAuthProvider(name string, authProvider AuthProvider) Option {
	return func(svc *service) {
		if _, ok := svc.authProviders[name]; ok {
			panic("duplicate auth provider registered")
		}

		svc.authProviders[name] = authProvider
	}
}

type AuthProvider interface {
	Validate(ctx context.Context, args map[string]interface{}) (res ValidateResponse, err error)
}

type ValidateResponse interface {
	Validated() bool
	ID() string
}
