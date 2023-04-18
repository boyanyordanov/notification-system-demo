package notifications

import (
	"log"
	"os"
)

type SlackLogChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *SlackLogChannel) GetType() string {
	return c.Type
}

func (c *SlackLogChannel) GetName() string {
	return c.Name
}

func (c *SlackLogChannel) Send(n Notification) (string, error) {
	log.Println("Sending slack message to channel: ", n.To)
	log.Println("Message: ", n.Message)
	log.Println("Processed by: ", os.Getpid())
	return "OK", nil
}
