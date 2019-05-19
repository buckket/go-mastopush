package mastopush

import (
	"crypto/ecdsa"
)

type Config struct {
	PrivateKey   ecdsa.PrivateKey
	ServerKey    ecdsa.PublicKey
	SharedSecret []byte
}

type MastoPush struct {
	*Config
}

type Payload struct {
	AccessToken      string `json:"access_token"`
	PreferredLocale  string `json:"preferred_locale"`
	NotificationID   ID     `json:"notification_id"`
	NotificationType string `json:"notification_type"`
	Icon             string `json:"icon"`
	Title            string `json:"title"`
	Body             string `json:"body"`
}

func NewMastoPush(config *Config) *MastoPush {
	return &MastoPush{Config: config}
}
