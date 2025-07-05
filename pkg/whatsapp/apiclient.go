package whatsapp

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type WhatsappClient interface {
	SendMessage(msg string, toNumber string) error
	SendInteractiveTemplate(templateSID string, toNumber string) error
	SendTemplateWithVariables(templateSID string, variables []string, toNumber string) error
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

func (client *whatsappClientImpl) SendTemplateWithVariables(templateSID string, variables []string, toNumber string) error {
	// Build content variables map for Twilio
	contentVariables := make(map[string]interface{})
	for i, variable := range variables {
		contentVariables[fmt.Sprintf("%d", i+1)] = variable
	}

	contentVariablesJSON, err := json.Marshal(contentVariables)
	if err != nil {
		return fmt.Errorf("failed to marshal content variables: %w", err)
	}

	params := &api.CreateMessageParams{}
	params.SetFrom(whatsappNumber(client.fromNumber))
	params.SetTo(whatsappNumber(toNumber))
	params.SetContentSid(templateSID)
	params.SetContentVariables(string(contentVariablesJSON))

	_, err = client.twilioClient.Api.CreateMessage(params)

	return err
}
