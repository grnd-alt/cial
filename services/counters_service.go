package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"backendsetup/m/db/sql/dbgen"

	"github.com/jackc/pgx/v5/pgtype"
)

type CountersService struct {
	queries *dbgen.Queries
}

func InitCountersService(queries *dbgen.Queries) *CountersService {
	return &CountersService{
		queries,
	}
}

func (c *CountersService) CanRead(counterID int, userID string) bool {
	_, err := c.queries.GetUserInCounter(context.Background(), dbgen.GetUserInCounterParams{
		UserID:    userID,
		CounterID: int32(counterID),
	})
	if err != nil {
		log.Printf("failed to share counter %d with %s: %v", counterID, userID, err)
		return false
	}
	return true
}

func (c *CountersService) CanShare(counterID int, userID string) bool {
	row, err := c.queries.GetUserInCounter(context.Background(), dbgen.GetUserInCounterParams{
		UserID:    userID,
		CounterID: int32(counterID),
	})
	if err != nil {
		log.Printf("failed to share counter %d with %s: %v", counterID, userID, err)
		return false
	}
	if row.AccessType.Valid && row.AccessType.String == "owner" {
		return true
	}
	return false
}

func (c *CountersService) GetCountersForUser(userID string) ([]dbgen.GetCountersForUserRow, error) {
	counts, err := c.queries.GetCountersForUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return counts, nil
}

func (c *CountersService) unsecureShareCounter(receivingUserID string, counterID int, accessType string) error {
	random := make([]byte, 32)
	_, err := rand.Read(random)
	if err != nil {
		return err
	}

	token := base64.StdEncoding.EncodeToString(random)

	return c.queries.AddUserToCounter(context.Background(), dbgen.AddUserToCounterParams{
		UserID:     receivingUserID,
		CounterID:  int32(counterID),
		Token:      pgtype.Text{Valid: true, String: token},
		AccessType: pgtype.Text{Valid: true, String: accessType},
	})
}

func (c *CountersService) CreateCounter(name string, icon string, creator string) (*dbgen.Counter, error) {
	counter, err := c.queries.CreateCounter(context.Background(), dbgen.CreateCounterParams{
		Name: pgtype.Text{Valid: true, String: name},
		Icon: pgtype.Text{Valid: true, String: icon},
	})
	if err != nil {
		return nil, err
	}
	err = c.unsecureShareCounter(creator, int(counter.ID), "owner")
	if err != nil {
		return nil, err
	}
	return &counter, nil
}


func (c *CountersService) ShareCounter(receivingUserID string, counterID int, accessType string, sharee string) error {
	if receivingUserID == sharee {
		return errors.New("can't share with yourself")
	}
	if !c.CanShare(counterID, sharee) {
		return errors.New("no permission")
	}

	random := make([]byte, 32)
	_, err := rand.Read(random)
	if err != nil {
		return err
	}

	token := base64.StdEncoding.EncodeToString(random)

	return c.queries.AddUserToCounter(context.Background(), dbgen.AddUserToCounterParams{
		UserID:     receivingUserID,
		CounterID:  int32(counterID),
		Token:      pgtype.Text{Valid: true, String: token},
		AccessType: pgtype.Text{Valid: true, String: accessType},
	})
}

type CounterUsers struct {
	Counter dbgen.Counter                `json:"counter"`
	Users   []dbgen.GetUsersInCounterRow `json:"users"`
}

func (c *CountersService) AddEvent(counterID int, userID string) error {
	// there is no can only read so read/add are the same :)
	if !c.CanRead(counterID, userID) {
		return errors.New("no permission")
	}
	err := c.queries.AddEventToCounter(context.Background(), dbgen.AddEventToCounterParams{
		UserID:    userID,
		CounterID: int32(counterID),
	})
	return err
}

func (c *CountersService) GetEvents(counterID int, userID string) ([]dbgen.CountersUsersEvent, error) {
	if !c.CanRead(counterID, userID) {
		return nil, errors.New("no permission")
	}
	res, err := c.queries.GetEvents(context.Background(), dbgen.GetEventsParams{
		UserID:    userID,
		CounterID: int32(counterID),
	})
	return res, err
}

func (c *CountersService) GetCounter(counterID int, userID string) (*CounterUsers, error) {
	if !c.CanRead(counterID, userID) {
		return nil, errors.New("no permission")
	}
	counter, err := c.queries.GetCounter(context.Background(), int32(counterID))
	if err != nil {
		return nil, err
	}
	users, err := c.queries.GetUsersInCounter(context.Background(), int32(counterID))
	if err != nil {
		return nil, err
	}
	result := CounterUsers{
		Counter: counter,
		Users:   users,
	}
	return &result, err
}
