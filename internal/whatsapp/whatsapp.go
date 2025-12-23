package whatsapp

import (
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"stock-picker/internal/config"
)

// Sender handles WhatsApp message sending via Twilio
type Sender struct {
	client *twilio.RestClient
	from   string
}

// NewSender creates a new WhatsApp sender instance
func NewSender(cfg *config.Config) *Sender {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &Sender{
		client: client,
		from:   cfg.TwilioWhatsAppFrom,
	}
}

// SendMessage sends a WhatsApp message using Twilio
func (w *Sender) SendMessage(to, message string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(w.from)
	params.SetBody(message)

	resp, err := w.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent successfully. SID: %s", *resp.Sid)
	return nil
}

