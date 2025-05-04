package services

import (
	"backendsetup/m/config"
	"backendsetup/m/db/sql/dbgen"
	"context"
	"encoding/json"
	"fmt"

	"github.com/SherClockHolmes/webpush-go"
)

type NotificationService struct {
	queries   *dbgen.Queries
	vapidPriv string
	vapidPub  string
}

func InitNotificationServe(conf *config.Config, queries *dbgen.Queries) *NotificationService {
	return &NotificationService{
		queries:   queries,
		vapidPriv: conf.VAPIDPriv,
		vapidPub:  conf.VAPIDPub,
	}
}

func (n *NotificationService) SendNotification(message string, userId string) error {

	subscriptions, err := n.queries.GetSubscriptions(context.Background(), userId)
	if err != nil {
		return err
	}
	for _, subscription := range subscriptions {
		s := &webpush.Subscription{}
		if err := json.Unmarshal([]byte(subscription), s); err != nil {
			return err
		}

		resp, err := webpush.SendNotification([]byte(message), s, &webpush.Options{
			Subscriber:      "test@belakkaf.net",
			VAPIDPublicKey:  n.vapidPub,
			VAPIDPrivateKey: n.vapidPriv,
			TTL:             30,
		})
		if resp.StatusCode != 201 || err != nil {
			if err := n.queries.DeleteSubscription(context.Background(), dbgen.DeleteSubscriptionParams{UserID: userId, Subscription: subscription}); err != nil {
				return err
			}
		}
		fmt.Println(resp)
	}
	return err
}

func (n *NotificationService) SendFollowersNotification(message string, userId string) error {
	res, err := n.queries.GetAllFollowers(context.Background(), userId)
	if err != nil {
		return err
	}
	for _, follower := range res {
		if err := n.SendNotification(message, follower.FollowerID); err != nil {
			return err
		}
	}
	return nil
}
