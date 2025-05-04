package services_test

import (
	"backendsetup/m/services"
	"backendsetup/m/tests"
	"testing"
)

func TestFollow(t *testing.T) {
	var followService = services.InitFollowService(tests.Queries)
	err := followService.Follow("testuser", "testuser2", []byte("test"))
	if err != nil {
		t.Error(err)
		return
	}

	followers, err := followService.GetFollowers("testuser2")
	if err != nil {
		t.Error(err)
		return
	}
	if len(followers) != 1 {
		t.Errorf("Expected 1 follower, got %d", len(followers))
		return
	}
	if followers[0].FollowerID != "1" {
		t.Errorf("Expected follower username 'testuser', got '%s'", followers[0].FollowerID)
		return
	}
}
