package notifications

import (
	"errors"
	"net/smtp"
	"strings"
	"testing"
)

func TestSend(t *testing.T) {
	n := Notification{
		MailFunc: func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
			return nil
		},
		SMTPServer:   "server:123",
		SMTPUsername: "user",
		SMTPPassword: "pass",
		Sender:       "me@testuser.com",
		Recipient:    "you@testuser.com",
	}
	err := n.Send("title", "https://mydomain.com")
	if err != nil {
		t.Fatal(err)
	}
}
func TestSendFail(t *testing.T) {
	n := Notification{
		MailFunc: func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
			return errors.New("this failure is expected")
		},
		SMTPServer:   "server:123",
		SMTPUsername: "user",
		SMTPPassword: "pass",
		Sender:       "me@testuser.com",
		Recipient:    "you@testuser.com",
	}
	err := n.Send("title", "https://mydomain.com")
	if err == nil {
		t.Fatal(err)
	}
}

func TestInitializeSuccess(t *testing.T) {
	n := Notification{}
	if Initialize(&n) != nil {
		t.Fatal("should have passed")
	}
	if savedNotification != &n {
		t.Fatal("didn't save notification correctly")
	}
	if savedNotification.MailFunc == nil {
		t.Fatal("didn't set mailfunc")
	}
	savedNotification.MailFunc = func(_ string, _ smtp.Auth, _ string, _ []string, msg []byte) error {
		return nil
	}
	if Send("title1", "https://bogus.com") != nil {
		t.Fatal("send after initialize failed")
	}
}
func TestInitializeFailure(t *testing.T) {
	n := Notification{SMTPUsername: "my-user-name"}
	err := Initialize(&n)
	if !strings.Contains(err.Error(), "No password") {
		t.Fatal("incorrect error message: " + err.Error())
	}
}
