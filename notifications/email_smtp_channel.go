package notifications

import (
	"net/smtp"
)

type SMTPEmailChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *SMTPEmailChannel) GetType() string {
	return c.Type
}

func (c *SMTPEmailChannel) GetName() string {
	return c.Name
}

func (c *SMTPEmailChannel) Send(n Notification) (string, error) {
	auth := smtp.Auth(smtp.PlainAuth("", c.Configuration["username"], c.Configuration["password"], c.Configuration["host"]))
	//msg := []byte("FSubject: Notification!\r\n" + "\r\n" + n.Message + ".\r\n")
	msg := []byte("From: " + c.Configuration["from_email"] + "\r\n" + "To: " + n.To + "\r\n" + "Subject: Notification!\r\n" + "\r\n" + n.Message + ".\r\n")
	err := smtp.SendMail(c.Configuration["host"]+":"+c.Configuration["port"], auth, c.Configuration["from_email"], []string{n.To}, msg)

	if err != nil {
		return "", err
	}
	return "OK", nil
}
