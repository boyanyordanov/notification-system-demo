# Notification Service Test Project

## Task 

Please create a notification sending system.
 - The system needs to be able to send notifications via several different channels (email, sms, slack) and be easily extensible to support more channels in the future.
 - The system needs to be horizontally scalable.
 - The system must guarantee an "at least once" SLA for sending the message.
 - The interface for accepting notifications to be sent can be chosen on your own discretion.

## Solution

The current solution consists of 3 main components:
 - Notification Service API
 - Notification Service Worker
 - Message Queue (in this case NATS)

Each component is supposed to solve one or more of the above objectives.

### Notification Service API

As an interface for accepting notifications to be sent, I've made a very simple API. 
It has just one endpoint, which accepts a POST request with a JSON body containing the following fields:
 - `type` - the type of notification (email, sms, slack and more)
 - `message` - the message to send
 - `to` - the recipient of the notification (email address, phone number, slack channel)

For the sake of simplicity, the API is not secured and does not have any validation.
It also uses only the components provided by the net/http package.

In a real-world scenario, I would use a framework like Gin or Gorilla Mux for the API.
I would also provide an OpenAPI definition and use it to generate documentation and maybe even help with basic validation. 
Validation will have to be dependant on the channel type because different message types have different restrictions.

The code for the api can be found in the /api directory.
It doesnt deal with state and only publishes messages to the message queue, so with minor modifications it can be scaled horizontally.

### Notifications Worker

The worker is responsible for receiving messages from the message queue and sending the notifications.
It first registers all the available notification channels. 
Then subscribes to a persistent queue and starts listening for messages.

The implementation of the connection with NATS is also not particularly robust but it could be improved enough for production use. 

The code for the worker can be found in the /cmd directory.
It is possible to run multiple instances of it concurrently and NATS will distribute the messages between them.

### Shared Code

The parts which can be used in both the API and the Worker are in the /notifications package.

The structure is pretty simple. 
We have a Notification type, which represents a notification to be sent.
It is used both for input and output in the api and as message in the message queue, because it already has all we need. 

The NotificationChannel interface is implemented by all the notification channels.
This allows extending the system with new channels without changing the code of the API or the Worker.
The only thing that needs to be done is to add a new implementation of the NotificationChannel interface.

Currently a channel is represented by a struct which contains its type, a name and a map which can hold any channel specific configurations. 
In other languages there would be a base channel which could be extended and reduce some of the duplication. 

The main difference between the channels is their Send function. 
Each channel has different setup, configuration and APIs. 

I am sure that this structure can be implemented more elegantly by someone with more experience with the language.

As demonstrated in the channels file, where the channels are registered and in the separate files for each channel, 
it's pretty straight forward to add a new one. 

Again, the design could be better because we are exposing the type field to the api directly and that may not be a good thing in a production environment.

### Message Queue

I went trough several different ideas in regards to the "send at least once requirement. 
First I was thinking about persisting the messages in a database or even a file and then sending them and checking for new ones upon restart of the process. 
However this will not be ideal and could create problems in a distributed environment. 

So in the end I decided to outsource it to a messaging system. 
I looked at RabbitMQ and NATS. I really liked how NATS handles interacting with many publishers and consumers and how it's build to scale while still maintaining a relatively simple interface. 
That said, it got a little bit more complicated in the code when I switched from th Core system, which doesn't include persistence
to the JetStream system, which does. 

Because this is a test project I've worked only with an installation on my local machine.
The only thing the application needs is a nats server running on the default port with jetstream enabled.
For example like this `nats-server -m 8222 -js` to enable the web interface as well.

For a production deployment it will require more research (or someone with more experience with NATS) and adjusting the configuration and probably the code too. 

## Available Channels 

#### For local testing: 
 - email-local - sends emails to the local mail server from the Helo app (basic SMTP)
 - sms-log - simply writes the message and recipient to the console
 - email-log - simply writes the message and recipient to the console
 - slack-log - simply writes the message and recipient to the console

#### Actual implementations
 - email - sends emails using SendinBlue, it uses the same SMTP implementation as the local one. But can be easily modified to use the api. 
   - It requires an account with SendinBlue, and their SMTP credentials for transactional emails.
 - sms - sends sms using Twilio. This one imports the Go library they provide, creates an API client and sends the message.
   - It requires an account with Twilio, SID, token and twilio number to send from.
 - slack - sends messages to a slack channel using the Slack API. It uses the net/http package to makea POST request to the Slack API for sending messages.
   - It requires an APP to be created and added to a workspace with the `channels:read` and `chat:write` scope and an API token. 

## Running the project

The channel configurations are loaded from the environment. 
The .evn-example file contains all the necessary variables.

For testing purposes I've only run it locally on my machine after building the api and cmd packages as executables and providing the necessary environment variables.
However it should be able to run in containers as well and be managed by something like Docker Swarm or Kubernetes.

I've provided a docker-compose file which will start the NATS server and the API and Worker containers.
It uses a .env file for configuration. However the only usable channels in that configuration are the log channels because 
of a few networking errors which I didn't have time to fix. 

It can simply run with `docker-compose up` and the API will be available on port 8000.
We can add the `--scale cmd=n` flag to run *n* workers and demonstrate how they spread.
Because the api service binds directly to port 8000 and compose exposes it for simplicity, it can't be scaled in the same way without modification. 

## Testing

TODO

For production the system will need actual unit tests and integration tests.
Due to my limited experience with the language some of the code will probably need to be rewritten to be more easy to test.

## Future improvements

There are several areas where the system can be improved. 

First it needs proper api documentation, authentication and validation.

The channels and NATS queue connections can be made configurable from a file during startup instead of hardcoding the process in the executable. 
This could allow running workers which only handle specific channels. And in cases where a similar structure is used (like SMTP emails), they can be replaced or added without a full rebuild.

The encryption issues which prevent actually sending notifications from the docker-compose setup should be fixed for easier local testing. 

As mentioned above a full test and benchmark suite should be added as well. 

The messages could be stored in the databse as a second layer of persistence, more stats and history.

Currently theres no observability in the services so that is also an area to research and develop further. 

If the service is going to be used only internally the api part can be replaced with something like grpc which is well integrated in the Go ecosystem and might be more performant.