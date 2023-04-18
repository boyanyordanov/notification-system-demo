package notifications

var Channels = make(map[string]NotificationChannel)

type NotificationChannel interface {
	Send(n Notification) (string, error)
	GetType() string
	GetName() string
}

func RegisterChannels(channelConfigurations []map[string]string) {
	for _, channelConfiguration := range channelConfigurations {
		channel := RegisterChannel(channelConfiguration["type"], channelConfiguration["name"], channelConfiguration)
		Channels[channel.GetType()] = channel
	}
}

func RegisterChannel(channelType, name string, configs map[string]string) NotificationChannel {
	switch channelType {
	case "sms-log":
		return &SMSLogChannel{Type: channelType, Name: name, Configuration: configs}
	case "email-log":
		return &EmailLogChannel{Type: channelType, Name: name, Configuration: configs}
	case "slack-log":
		return &SlackLogChannel{Type: channelType, Name: name, Configuration: configs}
	case "email-local":
		return &SMTPEmailChannel{Type: channelType, Name: name, Configuration: configs}
	case "email":
		return &SMTPEmailChannel{Type: channelType, Name: name, Configuration: configs}
	case "sms":
		return &SMSTwilioChannel{Type: channelType, Name: name, Configuration: configs}
	case "slack":
		return &SlackChannel{Type: channelType, Name: name, Configuration: configs}
	default:
		return nil
	}
}
