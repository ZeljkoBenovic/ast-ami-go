package types

type Config struct {
	Username        string
	Password        string
	Host            string
	Port            string
	InboundContext  string
	OutboundContext string
	LogFileLocation string
	WebhookURL      string
	WebhookMethod   string
	LogLevel        string
	AMIDebug        bool
}
