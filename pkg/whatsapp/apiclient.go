package whatsapp

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappClient interface {
	SendMessage(msg string, toNumber string) error
	SendInteractiveTemplate(templateSID string, toNumber string) error
}

type whatsappClientImpl struct {
	fromNumber   string
	twilioClient *twilio.RestClient
	log          echo.Logger
}

func NewWhatsappClient(fromNumber string, log echo.Logger) WhatsappClient {
	twilioClient := twilio.NewRestClient()
	return &whatsappClientImpl{fromNumber, twilioClient, log}
}

func whatsappNumber(phoneNumber string) string {
	return fmt.Sprintf("whatsapp:%s", phoneNumber)
}

func (client *whatsappClientImpl) SendMessage(msg string, toNumber string) error {
	params := &api.CreateMessageParams{}
	params.SetFrom(whatsappNumber(client.fromNumber))
	params.SetTo(whatsappNumber(toNumber))
	params.SetBody(msg)

	_, err := client.twilioClient.Api.CreateMessage(params)

	return err
}

func (client *whatsappClientImpl) SendInteractiveTemplate(templateSID string, toNumber string) error {
	params := &api.CreateMessageParams{}
	params.SetFrom(whatsappNumber(client.fromNumber))
	params.SetTo(whatsappNumber(toNumber))
	params.SetContentSid(templateSID)

	_, err := client.twilioClient.Api.CreateMessage(params)

	return err
}
