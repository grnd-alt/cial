package services_test

import (
	"backendsetup/m/services"
	"backendsetup/m/tests"
	"log"
	"os"
	"testing"
)

var userService *services.UserService

func TestMain(m *testing.M) {
	os.Chdir("../..")
	tests.SetUp()
	userService = services.InitUserService(tests.Queries)
	m.Run()
}

func TestCreateUserIfNotExists(t *testing.T) {
	username := "testuser3"
	userId := "3"

	userService.CreateUserIfNotExists(username, userId)

	user, err := userService.GetUserByUsername(username)
	if err != nil {
		log.Fatal(err)
		t.Fail()
		return
	}
	if user.UserID != userId {
		log.Fatal(err)
		t.Fail()
		return
	}
}
