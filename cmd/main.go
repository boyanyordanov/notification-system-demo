package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"notifications-system/notifications"
	"os"
	"runtime"
	"time"
)

func main() {
	// Register all available channels
	notifications.RegisterChannels([]map[string]string{
		{
			"type":       "email-local",
			"name":       "helo (local)",
			"host":       os.Getenv("SMTP_LOCAL_HOST"),
			"port":       os.Getenv("SMTP_LOCAL_PORT"),
			"username":   os.Getenv("SMTP_LOCAL_USERNAME"),
			"password":   "",
			"from_email": os.Getenv("SMTP_LOCAL_FROM_EMAIL"),
		},
		{
			"type":     "email",
			"name":     "sendinblue",
			"host":     os.Getenv("SENDINBLUE_SMTP_HOST"),
			"port":     os.Getenv("SENDINBLUE_SMTP_PORT"),
			"username": os.Getenv("SENDINBLUE_SMTP_USERNAME"),
			"password": os.Getenv("SENDINBLUE_SMTP_PASSWORD"),
		},
		{
			"type":    "sms-log",
			"name":    "log",
			"host":    "127.0.0.1",
			"api_key": "asdfg",
			"sid":     "123",
		},
		{
			"type":        "sms",
			"name":        "twilio",
			"account_sid": os.Getenv("TWILIO_ACCOUNT_SID"),
			"token":       os.Getenv("TWILIO_AUTH_TOKEN"),
			"from":        os.Getenv("TWILIO_FROM_NUMBER"),
		},
		{
			"type":  "slack",
			"name":  "slack",
			"token": os.Getenv("SLACK_TOKEN"),
		},
	})

	log.Println("Listening for notifications on channels:")
	for _, channel := range notifications.Channels {
		log.Println("Channel: ", channel.GetType(), " - ", channel.GetName())
	}

	// Setup the NATS connection and consumer
	nc, _ := nats.Connect(os.Getenv("NATS_SERVER"))
	defer nc.Drain()
	js, _ := nc.JetStream()

	_, jsErr := js.AddConsumer("NOTIFICATIONS", &nats.ConsumerConfig{
		Durable:        "channels",
		DeliverSubject: "channels",
		DeliverGroup:   "channels-group",
		AckPolicy:      nats.AckExplicitPolicy,
		AckWait:        time.Second,
	})

	if jsErr != nil {
		log.Println("Error creating consumer: ", jsErr)
		return
	}

	// Subscribe to the queue using the "durable" consumer
	// This allows to consume messages already in the queue upon start and listening for new messages going forward
	sub, subErr := js.QueueSubscribeSync("notifications.>", "channels-group")
	if subErr != nil {
		log.Println("Error creating subscription: ", subErr)
		return
	}

	for {
		m, err := sub.NextMsg(time.Second * 5)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, processErr := processMsg(m)

		if processErr != nil {
			log.Println("Error processing message: ", processErr)
			continue
		}

		m.Ack()
	}

	runtime.Goexit()
}

func processMsg(msg *nats.Msg) (string, error) {
	// parse the message from the queue
	var notification notifications.Notification
	err := json.Unmarshal([]byte(msg.Data), &notification)

	if err != nil {
		log.Println("Unable to decode notification: ", err)
		return "", errors.New("unable to decode notification")
	}

	// if everything is alright process the notification
	log.Println("Received notification: ", notification)
	result, err := processNotification(notification)

	if err != nil {
		return "", err
	}

	return result, nil
}

func processNotification(n notifications.Notification) (string, error) {
	if _, ok := notifications.Channels[n.Type]; !ok {
		log.Println("Unknown channel type: ", n.Type)
		return "", errors.New("unknown channel type")
	}

	result, err := notifications.Channels[n.Type].Send(n)
	if err != nil {
		log.Println("Error sending notification: ", err)
		return "", err
	}

	log.Println(result)
	return result, nil
}
