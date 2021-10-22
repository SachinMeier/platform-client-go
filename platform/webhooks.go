package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type Webhook struct {
	Url     string `json:"url"`
	Secret  string `json:"secret" default:""`
	Enabled bool   `json:"enabled"`
}

func (pc *PlatformClient) handleWebhookRequest(req *http.Request, err error) (Webhook, error) {
	if err != nil {
		log.Error("Internal Error")
		return Webhook{}, err
	}

	var webhook Webhook
	err = pc.sendRequest(req, &webhook)
	if err != nil {
		log.Error("Webhook Request Failed")
		return Webhook{}, err
	}
	return webhook, nil
}

// SubscribeToWebhook subscribes to a webhook
func (pc *PlatformClient) SubscribeToWebhook(callback_url string) (Webhook, error) {
	log.Infof("Subscribing to Webhook %s", callback_url)

	data := map[string]string{
		"url": callback_url,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Errorf("JSON encoding error with url: %s", callback_url)
		return Webhook{}, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/accounts/%s/webhooks/", pc.BaseURL, pc.accountId),
		bytes.NewBuffer(body),
	)
	return pc.handleWebhookRequest(req, err)
}

// GetSubscribedWebhook queries subscribed webhook
func (pc *PlatformClient) GetSubscribedWebhook() (Webhook, error) {
	log.Infof("Querying Webhook")
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/accounts/%s/webhooks/", pc.BaseURL, pc.accountId),
		nil,
	)
	return pc.handleWebhookRequest(req, err)
}

// DeleteWebhook deletes the existing webhook
func (pc *PlatformClient) DeleteWebhook() bool {
	log.Infof("Querying Webhook")
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/accounts/%s/webhooks/", pc.BaseURL, pc.accountId),
		nil,
	)
	if err != nil {
		log.Error("Internal Error")
		return false
	}

	err = pc.sendRequest(req, nil)
	if err != nil {
		log.Error("Delete Webhook Failed")
		return false
	}
	return true
}
