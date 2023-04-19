package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"notifications-system/notifications"
	"os"
)

func main() {

	// Set up the NATS connection and stream
	nc, _ := nats.Connect(os.Getenv("NATS_SERVER"))
	js, _ := nc.JetStream()

	// Create a Stream
	_, err := js.AddStream(&nats.StreamConfig{
		Name:     "NOTIFICATIONS",
		Subjects: []string{"notifications.>"},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new handler for the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Try creating a notification with POST /notifications."))
		if err != nil {
			return
		}
	})

	// Attach a function as handler for the more complicated endpoint
	http.HandleFunc("/notifications", handleNotificationCreation(js))

	// Start the server
	log.Println("Listening on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleNotificationCreation(js nats.JetStream) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var notification notifications.Notification
			err := json.NewDecoder(r.Body).Decode(&notification)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error parsing JSON"))
				return
			}

			resp, _ := json.Marshal(notification)

			// Publish the notification to the queue
			pubAck, pubErr := js.Publish("notifications.1", resp)

			if pubErr != nil {
				log.Println(pubErr)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Error sending notification"))
				return
			}

			log.Println(pubAck)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(resp))
		} else {
			// We're only interested in creating notifications for this project so everything else will be marked as not implemented
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("{\"error\": \"Not implemented. Only available operation is POST.\"}"))
		}
	}
}
