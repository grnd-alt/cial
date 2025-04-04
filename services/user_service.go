package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
)

type UserService struct {
	queries *dbgen.Queries
}

func InitUserService(queries *dbgen.Queries) *UserService {
	return &UserService{
		queries,
	}
}

func (u *UserService) CreateUserIfNotExists(username string, userId string) {
	u.queries.InsertUser(context.Background(), dbgen.InsertUserParams{Username: username, UserID: userId})
}

func (u *UserService) GetUserByUsername(username string) (dbgen.User, error) {
	return u.queries.GetUserByName(context.Background(), username)
}

func (u *UserService) InsertSubscription(userId string, subscription []byte) error {
	return u.queries.InsertSubscription(context.Background(), dbgen.InsertSubscriptionParams{UserID: userId, Subscription: subscription})
}
