package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

//SlackRequestBody struct used to define slack request body
type SlackRequestBody struct {
	Text     string `json:"text"`
	Username string `json:"username,omitempty"`
	Channel  string `json:"channel,omitempty"`
}

//SendSlackNotification send a notification to slack. Webhook URL is fetched from config file (slack.url) using Viper
func SendSlackNotification(msg SlackRequestBody) error {
	if !viper.GetBool("slack.enabled") {
		return nil
	}
	webhookURL := viper.GetString("slack.url")
	slackBody, _ := json.Marshal(msg)
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
