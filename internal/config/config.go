package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	TwilioAccountSID   string
	TwilioAuthToken    string
	TwilioWhatsAppFrom string
	WhatsAppRecipients []string
}

// NewConfig creates a new Config from environment variables
func NewConfig() *Config {
	recipients := strings.Split(os.Getenv("WHATSAPP_RECIPIENTS"), ",")
	// Clean up recipients (remove spaces)
	cleanedRecipients := []string{}
	for _, r := range recipients {
		cleaned := strings.TrimSpace(r)
		if cleaned != "" {
			cleanedRecipients = append(cleanedRecipients, cleaned)
		}
	}

	return &Config{
		TwilioAccountSID:   os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:    os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioWhatsAppFrom: os.Getenv("TWILIO_WHATSAPP_FROM"),
		WhatsAppRecipients: cleanedRecipients,
	}
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	if c.TwilioAccountSID == "" {
		return fmt.Errorf("TWILIO_ACCOUNT_SID is required")
	}
	if c.TwilioAuthToken == "" {
		return fmt.Errorf("TWILIO_AUTH_TOKEN is required")
	}
	if c.TwilioWhatsAppFrom == "" {
		return fmt.Errorf("TWILIO_WHATSAPP_FROM is required")
	}
	if len(c.WhatsAppRecipients) == 0 {
		return fmt.Errorf("WHATSAPP_RECIPIENTS is required")
	}
	return nil
}
