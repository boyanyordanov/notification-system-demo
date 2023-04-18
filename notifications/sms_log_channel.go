package notifications

import (
	"log"
	"os"
)

type SMSLogChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *SMSLogChannel) GetType() string {
	return c.Type
}

func (c *SMSLogChannel) GetName() string {
	return c.Name
}

func (c *SMSLogChannel) Send(n Notification) (string, error) {
	log.Println("Sending SMS to: ", n.To)
	log.Println("Message: ", n.Message)
	log.Println("Processed by: ", os.Getpid())
	return "OK", nil
}
