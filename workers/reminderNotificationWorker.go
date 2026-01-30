package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"backendsetup/m/db/sql/dbgen"
	"backendsetup/m/services"

	"github.com/jackc/pgx/v5/pgtype"
)

type ReminderNotificationWorker struct {
	queries             *dbgen.Queries
	notificationService *services.NotificationService
}

func NewReminderNotificationWorker(queries *dbgen.Queries, notificationService *services.NotificationService) ReminderNotificationWorker {
	return ReminderNotificationWorker{
		queries,
		notificationService,
	}
}

func (r ReminderNotificationWorker) SendNotifications() {
	users, err := r.queries.GetNoLoggedInSince(context.Background(), dbgen.GetNoLoggedInSinceParams{
		LastLogin: pgtype.Timestamptz{Valid: true, Time: time.Now().Add(-48 * time.Hour)},
		Limit:     1000,
	})
	if err != nil {
		log.Printf("failed to receive Users to notify %v", err)
		return
	}

	for _, user := range users {
		if user.LastNotified.Time.After(time.Now().Add(-48 * time.Hour)) {
			return
		}
		notificationData := services.NotificationData{
			Type:  services.ReminderNotificationType,
			Title: fmt.Sprintf("Hey %s, jump back in", user.Username),
			Body:  "There might be new posts waiting for you",
		}
		err = r.notificationService.SendNotification(notificationData, user.UserID)
		if err != nil {
			log.Printf("failed to notify %s: %v", user.UserID, err)
		}
	}
}

func (r ReminderNotificationWorker) StartWorker() {
	for {
		r.SendNotifications()
		time.Sleep(time.Hour * 8)
	}
}
