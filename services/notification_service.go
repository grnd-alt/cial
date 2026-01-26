package services

import (
	"context"
	"encoding/json"
	"log"

	"backendsetup/m/config"
	"backendsetup/m/db/sql/dbgen"

	"github.com/SherClockHolmes/webpush-go"
)

type NotificationType string

const (
	ReminderNotificationType  NotificationType = "reminder"
	NewPostNotificationType   NotificationType = "newPost"
	NewFollowNotificationType NotificationType = "newFollow"
)

type NotificationData struct {
	Type  NotificationType `json:"type"`
	Title string           `json:"title"`
	Body  string           `json:"body"`
	Data  any              `json:"data,omitempty"`
}

type NewPostData struct {
	Author string `json:"author"`
}

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

func (n *NotificationService) SendNotification(data NotificationData, userID string) error {
	subscriptions, err := n.queries.GetSubscriptions(context.Background(), userID)
	if err != nil {
		return err
	}
	err = n.queries.SetLastNotified(context.Background(), userID)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		s := &webpush.Subscription{}
		if err := json.Unmarshal([]byte(subscription), s); err != nil {
			return err
		}

		resp, err := webpush.SendNotification(payload, s, &webpush.Options{
			Subscriber:      "test@belakkaf.net",
			VAPIDPublicKey:  n.vapidPub,
			VAPIDPrivateKey: n.vapidPriv,
			TTL:             30,
		})
		if err != nil || resp.StatusCode != 201 {
			log.Printf("sending notification failed to %s with code %d %v", userID, resp.StatusCode, err)
			if err := n.queries.DeleteSubscription(context.Background(), dbgen.DeleteSubscriptionParams{UserID: userID, Subscription: subscription}); err != nil {
				return err
			}
			return err
		}
	}
	return err
}

func (n *NotificationService) SendFollowersNotification(data NotificationData, userID string) error {
	res, err := n.queries.GetAllFollowers(context.Background(), userID)
	if err != nil {
		return err
	}
	for _, follower := range res {
		if err := n.SendNotification(data, follower.FollowerID); err != nil {
			return err
		}
	}
	return nil
}
