package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type FollowService struct {
	queries *dbgen.Queries
}

func (f *FollowService) GetFollowersCount(username string) (int64, error) {
	user, err := f.queries.GetUserByName(context.Background(), username)
	if err != nil {
		return 0, err
	}
	res, err := f.queries.GetFollowersCount(context.Background(), user.UserID)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (f *FollowService) GetFollowingCount(username string) (int64, error) {
	user, err := f.queries.GetUserByName(context.Background(), username)
	if err != nil {
		return 0, err
	}
	res, err := f.queries.GetFollowingCount(context.Background(), user.UserID)
	if err != nil {
		return 0, err
	}
	return res, nil

}

func InitFollowService(queries *dbgen.Queries) *FollowService {
	return &FollowService{
		queries: queries,
	}
}

func (f *FollowService) GetFollowers(username string) ([]dbgen.UserFollow, error) {
	user, err := f.queries.GetUserByName(context.Background(), username)
	if err != nil {
		return nil, err
	}
	res, err := f.queries.GetFollowers(context.Background(), dbgen.GetFollowersParams{FollowedID: user.UserID})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (f *FollowService) Follow(followerName string, followingName string, subscription []byte) error {
	follower, err := f.queries.GetUserByName(context.Background(), followerName)
	if err != nil {
		return err
	}

	following, err := f.queries.GetUserByName(context.Background(), followingName)
	if err != nil {
		return err
	}

	f.queries.InsertSubscription(context.Background(), dbgen.InsertSubscriptionParams{
		UserID:       follower.UserID,
		Subscription: subscription,
	})

	return f.queries.InsertFollower(context.Background(), dbgen.InsertFollowerParams{
		FollowerID:       follower.UserID,
		FollowedID:       following.UserID,
		NotificationType: pgtype.Text{String: "posts", Valid: true},
	})
}
