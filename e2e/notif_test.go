//+build e2e

package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"notif2/config"
	"notif2/smtp"
	"notif2/svcs"
	"notif2/tools"
)

var (
	users       svcs.UserService
	mailChecker smtp.EMailChecker
	notifSrvURL string
)

func TestDummy(t *testing.T) {
	t.Logf("Test infra works")
}

func TestMain(m *testing.M) {

	notifSrvURL = config.StringOption("NOTIF_SVC_URL")
	userSrvURL := config.StringOption("USER_SVC_URL")
	var err error
	users, err = svcs.NewBasicUserService(userSrvURL)
	tools.FailOnError("Could not instantiate User Svc", err)
	mailChecker, err = smtp.NewMailHogAPI(config.StringOption("MH_URL"))
	tools.FailOnError("Could not instantiate MailHog API", err)

	exitStatus := m.Run()

	os.Exit(exitStatus)
}

func TestNotifyAllUsersByEmail(t *testing.T) {
	const (
		TEST_ORIGIN  = "from@example.com"
		TEST_SUBJECT = "Test subject"
		TEST_BODY    = "Test body"
	)

	givenWeHaveSomeUsers := func() {
	}

	whenINotifyThemByEmailWithTheMessage := func(orig, subject, body string) {
		type notifRequest struct {
			From    string `json:"from"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
		}
		makeNotifRequest := func(from, subject, body string) (*http.Request, error) {
			reqData := &notifRequest{from, subject, body}
			reqBody, err := json.Marshal(reqData)
			if err != nil {
				return nil, err
			}
			req, err := http.NewRequest(
				http.MethodPost,
				notifSrvURL+"notif",
				bytes.NewBuffer(reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
			t.Logf("Sending request: %v", req)
			if err != nil {
				return nil, err
			}
			return req, nil
		}
		req, err := makeNotifRequest(TEST_ORIGIN, TEST_SUBJECT, TEST_BODY)
		if err != nil {
			t.Fatalf("Unexpected error when building request: %v", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Unexpected error when raising request: %v", err)
		}
		if resp.StatusCode/100 != 200 {
			t.Fatalf("Bad response from server: %d", resp.StatusCode)
		}
	}

	ourUsersReceiveTheSentEmails := func(subject, body string) {
		users, err := users.FindAll()
		if err != nil {
			t.Fatalf("Could not retrieve user list: %v", err)
		}
		for _, u := range users {
			isThere, err := mailChecker.CheckEmail(u.EMail, TEST_ORIGIN, TEST_SUBJECT, TEST_BODY)
			if err != nil {
				t.Errorf("Error when checking e-mail presence: %v", err)
			}
			if !isThere {
				t.Errorf("E-mail to: %s expected to be found", u.EMail)
			}
		}
	}

	givenWeHaveSomeUsers()
	whenINotifyThemByEmailWithTheMessage(TEST_ORIGIN, TEST_SUBJECT, TEST_BODY)
	ourUsersReceiveTheSentEmails(TEST_SUBJECT, TEST_BODY)
}
