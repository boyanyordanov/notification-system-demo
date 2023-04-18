package notifications

import (
	"log"
	"os"
)

type EmailLogChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *EmailLogChannel) GetType() string {
	return c.Type
}

func (c *EmailLogChannel) GetName() string {
	return c.Name
}

func (c *EmailLogChannel) Send(n Notification) (string, error) {
	log.Println("Sending Email to: ", n.To)
	log.Println("Message: ", n.Message)
	log.Println("Processed by: ", os.Getpid())
	return "OK", nil
}
