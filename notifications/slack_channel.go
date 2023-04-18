package notifications

import (
	"bytes"
	"log"
	"net/http"
	"time"
)

type SlackChannel struct {
	Type          string
	Name          string
	Configuration map[string]string
}

func (c *SlackChannel) GetType() string {
	return c.Type
}

func (c *SlackChannel) GetName() string {
	return c.Name
}

func (c *SlackChannel) Send(n Notification) (string, error) {
	log.Println("Sending Slack message to: ", n.To)
	log.Println("Message: ", n.Message)

	jsonBody := []byte(`{"channel": "` + n.To + `", "text": "` + n.Message + `"}`)

	log.Println(string(jsonBody))

	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.postMessage", bodyReader)

	log.Println(req)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.Configuration["token"])

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, resErr := client.Do(req)

	if resErr != nil {
		return "", resErr
	}

	log.Println(res)

	var resp []byte
	_, err = res.Body.Read(resp)

	if err != nil {
		return "", err
	}

	return string(resp), nil
}
