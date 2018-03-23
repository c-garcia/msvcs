package users_test

import (
	"notif2/config"
	"notif2/svcs"
	"os"
	"strings"
	"testing"
)

var (
	userSrvURL string
)

func TestMain(m *testing.M) {
	userSrvURL = config.StringOption("USER_SVC_URL")
	exitStatus := m.Run()
	os.Exit(exitStatus)
}

func TestFindAllListsUsers(t *testing.T) {
	sut, err := svcs.NewBasicUserService(userSrvURL)
	if err != nil {
		t.Errorf("Could not instantiate UserService: %s", err)
	}

	users, err := sut.FindAll()
	if err != nil {
		t.Errorf("FindAll returned an error: %s", err)
	}
	if len(users) < 1 {
		t.Errorf("FindAll should have found some users but found: %d", len(users))
	}
	for _, u := range users {
		if len(u.Name) < 1 {
			t.Errorf("FindAll found an empty Name")
		}
		if strings.Index(u.EMail, "@") < 1 {
			t.Errorf("FindAll found an incorrect e-mail: %s", u.EMail)
		}
	}
}
