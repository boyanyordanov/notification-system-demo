package notifications

import (
	"fmt"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMSTwilioChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *SMSTwilioChannel) GetType() string {
	return c.Type
}

func (c *SMSTwilioChannel) GetName() string {
	return c.Name
}

func (c *SMSTwilioChannel) Send(n Notification) (string, error) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: c.Configuration["account_sid"],
		Password: c.Configuration["token"],
	})

	params := &api.CreateMessageParams{}
	params.SetBody(n.Message)
	params.SetFrom(c.Configuration["from"])
	params.SetTo(n.To)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}

	return "OK", nil
}
