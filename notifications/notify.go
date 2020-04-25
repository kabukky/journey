package notifications

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"net/url"
	"strings"
)

// NotifyMessage is an interface to send notification messages for new posts
// this gets sent when a post is set to 'published' for the first time
type NotifyMessage interface {
	send(title string, url string) error
}

// Notification contains configuration for how to notify people of new
// posts via smtp
type Notification struct {
	MailFunc     func(string, smtp.Auth, string, []string, []byte) error `json:"-"`
	SMTPServer   string
	SMTPUsername string
	SMTPPassword string
	Sender       string
	Recipient    string
}

// TemplateText contains the HTML that is filled in and emailed to someone
// This could get moved to either a file or inline in the configuration.
const TemplateText = `To: {{.Recipient}}
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";
Subject: New Post: {{.Title}}


<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
        "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
<head></head>
<body>
<p>Hello,</p>
<p>There's a new post on <b>{{.Domain}}</b> entitled "<a href="{{.URL}}">{{.Title}}"</a></p>
<p>Please visit {{.URL}} to view it.</p>
<p style="font-size:small">Please don't reply, as this account is not monitored.</p>
</body>

</html>`

// Send implements sending a post to someone
func (n *Notification) Send(title string, posturl string) error {
	components := strings.SplitN(n.SMTPServer, ":", 2)
	auth := smtp.PlainAuth("", n.SMTPUsername, n.SMTPPassword, components[0])
	templateData := struct {
		Title     string
		URL       string
		Domain    string
		Recipient string
	}{
		Title:     title,
		URL:       posturl,
		Recipient: n.Recipient,
	}

	if parsed, err := url.Parse(posturl); err == nil {
		templateData.Domain = parsed.Host
	}

	buf := new(bytes.Buffer)
	t, err := template.New("email").Parse(TemplateText)
	if err != nil {
		return err
	}

	if err = t.Execute(buf, templateData); err != nil {
		return err
	}

	body := []byte(buf.String())

	if err := n.MailFunc(n.SMTPServer, auth, n.Sender, []string{n.Recipient}, body); err != nil {
		log.Println("Send mail failed: " + err.Error())
		return err
	}

	return nil
}

var savedNotification *Notification

// Initialize remembers the notification configuration for later use
func Initialize(n *Notification) error {
	n.MailFunc = smtp.SendMail
	if n.SMTPUsername != "" && n.SMTPPassword == "" {
		return fmt.Errorf("No password specified for SMTP user %q", n.SMTPUsername)
	}
	savedNotification = n
	return nil
}

// Send sends an email using the globally saved configuration
func Send(title string, posturl string) error {
	return savedNotification.Send(title, posturl)
}
