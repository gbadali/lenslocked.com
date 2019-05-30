package context

import (
	"context"

	"github.com/gbadali/lenslocked.com/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// WithUser takes a context and a User and puts the user in the context
// it uses the userKey private type to prevent anyone else from adding
// invalid data to the context with the "user" type
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User retrieves data user data from the context, this func is needed
// because the data that we put into the context with WithUser can
// only be accessed within this context package
func User(ctx context.Context) *models.User {
	// first check to make sure it exists
	if temp := ctx.Value(userKey); temp != nil {
		// then check to make sure it is the correct type
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
