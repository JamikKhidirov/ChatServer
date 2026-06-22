package sms

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Sender struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewSender(accountSID, authToken, fromNumber string) *Sender {
	return &Sender{AccountSID: accountSID, AuthToken: authToken, FromNumber: fromNumber}
}

func (s *Sender) Send(to, message string) error {
	if s.AccountSID == "" || s.AuthToken == "" {
		return fmt.Errorf("SMS not configured: set TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN env vars")
	}
	addr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.AccountSID)
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", s.FromNumber)
	data.Set("Body", message)
	req, err := http.NewRequest("POST", addr, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	auth := base64.StdEncoding.EncodeToString([]byte(s.AccountSID + ":" + s.AuthToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Twilio API error: %s", resp.Status)
	}
	return nil
}
