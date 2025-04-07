package sms

import (
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"log"
)

type SmsClient struct {
	Client     *twilio.RestClient
	FromNumber string
}

func NewSmsClient(client *twilio.RestClient, AccountSID, AuthToken, FromNumber string) *SmsClient {
	return &SmsClient{client, FromNumber}
}

func (s *SmsClient) SendSMS(text, toNumber string) error {
	params := &api.CreateMessageParams{}
	params.SetBody(text)
	params.SetFrom(s.FromNumber)

	params.SetTo(toNumber)

	resp, err := s.Client.Api.CreateMessage(params)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if resp.Body != nil {
		log.Println(*resp.Body)
		return err
	}

	return nil
}
