package mailhog_test

import (
	"notif2/config"
	"notif2/smtp"
	"notif2/tools"
	"os"
	"testing"
)

var searchEndPointCases = []struct {
	api    string
	search string
}{
	{api: "http://example.com/api", search: "http://example.com/api/v2/search"},
	{api: "http://example.com/api/", search: "http://example.com/api/v2/search"},
	{api: "https://example.com/api", search: "https://example.com/api/v2/search"},
	{api: "https://example.com/api/", search: "https://example.com/api/v2/search"},
}

var (
	sut        *smtp.MailHogAPI
	mailServer string
	mailPort   int
)

func TestMain(m *testing.M) {
	mailServer = config.StringOption("SMTP_HOST")
	mailPort = config.IntOption("SMTP_PORT")
	var err error
	sut, err = smtp.NewMailHogAPI(config.StringOption("MH_URL"))
	tools.FailOnError("Cound not create API endpoint", err)
	s := m.Run()
	os.Exit(s)
}

func TestDeleteAllEmails(t *testing.T) {
	const (
		to      = "to@example.com"
		from    = "from@example.com"
		subject = "test subject"
		body    = "test body"
	)
	err := smtp.SendEmail(mailServer, mailPort, from, to, subject, body)
	if err != nil {
		t.Fatalf("Unexpected error when sending e-mail: %v", err)
	}
	count, err := sut.MessageCount()
	if err != nil {
		t.Fatalf("Unexpected error counting existing messages: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 message present. Got: %d", count)
	}
	err = sut.DeleteAllEmails()
	if err != nil {
		t.Fatalf("Unexpected error while deleting e-mails: %v", err)
	}
	count, err = sut.MessageCount()
	if err != nil {
		t.Fatalf("Unexpected error counting existing messages: %v", err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 message present after deletion. Got: %d", count)
	}
}

func TestMessageCountOfEmptyIsZero(t *testing.T) {
	err := sut.DeleteAllEmails()
	if err != nil {
		t.Fatalf("Unexpected error deleting emails: %v", err)
	}
	count, err := sut.MessageCount()
	if err != nil {
		t.Fatalf("Unexpected error counting emails: %v", err)
	}
	if count != 0 {
		t.Fatalf("Expected 0 e-mails. Got: %d", count)
	}
}

func TestMessageCountOfNonEmptyIsCorrect(t *testing.T) {
	const (
		from    = "from@example.com"
		to      = "to@example.com"
		subject = "test subject"
		body    = "test body"
	)
	err := sut.DeleteAllEmails()
	if err != nil {
		t.Fatalf("Unexpected error deleting emails: %v", err)
	}
	smtp.SendEmail(mailServer, mailPort, from, to, subject, body)
	count, err := sut.MessageCount()
	if err != nil {
		t.Fatalf("Unexpected error counting emails: %v", err)
	}
	if count != 1 {
		t.Fatalf("Expected 1 e-mails. Got: %d", count)
	}
}

var findsEmailCases = []struct {
	From     string
	To       string
	Subject  string
	Body     string
	MustFind bool
}{
	{"u1@example.com", "u2@example.com", "Test subject", "Test body", true},
	{"xx@example.com", "u2@example.com", "Test subject", "Test body", false},
	{"u1@example.com", "xx@example.com", "Test subject", "Test body", false},
	{"u1@example.com", "u2@example.com", "Other subject", "Test body", false},
	{"u1@example.com", "u2@example.com", "Test subject", "Other body", false},
}

func TestCheckEmailFindsEmail(t *testing.T) {
	const (
		from    = "u1@example.com"
		to      = "u2@example.com"
		subject = "Test subject"
		body    = "Test body"
	)

	for _, tc := range findsEmailCases {
		err := sut.DeleteAllEmails()
		if err != nil {
			t.Errorf("Error when deleting e-mails: %v", err)
		}
		smtp.SendEmail(mailServer, mailPort, from, to, subject, body)
		isThere, err := sut.CheckEmail(tc.To, tc.From, tc.Subject, tc.Body)
		if err != nil {
			t.Errorf("Error when checking email: %v", err)
		}
		if isThere != tc.MustFind {
			t.Errorf("Expected e-mail to %s (%s) to be found: %v", to, subject, tc.MustFind)
		}
	}
}

func TestCheckEmailDoesNotFoundOnEmpty(t *testing.T) {
	const (
		from    = "u1@example.com"
		to      = "u2@example.com"
		subject = "Test subject"
		body    = "test body"
	)

	err := sut.DeleteAllEmails()
	if err != nil {
		t.Errorf("Error when deleting e-mails: %v", err)
	}
	isThere, err := sut.CheckEmail(to, from, subject, body)
	if err != nil {
		t.Errorf("Error when checking email: %v", err)
	}
	if isThere {
		t.Errorf("Expected e-mail to %s (%s) not to be found", to, subject)
	}
}
