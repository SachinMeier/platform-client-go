package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type DepositInvoice struct {
	Id          string `json:"id"`
	Destination string `json:"destination"`
	Network     string `json:"network"`
	Timestamp   int    `json:"timestamp"`
}

type DepositInvoiceListOptions struct {
	Page    int `url:"page"`
	PerPage int `url:"per_page"`
}

// CreateDepositInvoice creates an invoice to enable deposits to River Platform
func (pc *PlatformClient) CreateDepositInvoice(amount int, label, network string) (DepositInvoice, error) {
	log.Info("Requesting Deposit Invoice")

	data := map[string]interface{}{
		"amount":  amount,
		"network": network,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Errorf("JSON encoding error with amount: %d or network: %s", amount, network)
		return DepositInvoice{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/%s/deposit_intents", pc.BaseURL, pc.accountId), bytes.NewBuffer(body))
	if err != nil {
		log.Error("Internal Error")
		return DepositInvoice{}, err
	}

	var invoice DepositInvoice
	err = pc.sendRequest(req, &invoice)
	if err != nil {
		log.Errorf("Create Invoice Failed: %s", err.Error())
		return DepositInvoice{}, err
	}
	return invoice, nil
}

// GetDepositInvoices queries a list of invoices generated by River Platform
func (pc *PlatformClient) GetDepositInvoices(limit, next_timestamp int) ([]DepositInvoice, error) {
	log.Info("Querying Deposit Invoices")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/deposit_intents", pc.BaseURL, pc.accountId), nil)
	if err != nil {
		log.Errorf("Internal Error: %s", err.Error())
		return nil, err
	}

	var invoices []DepositInvoice
	err = pc.sendRequest(req, &invoices)
	if err != nil {
		log.Errorf("Create Invoice Failed: %s", err.Error())
		return nil, err
	}
	return invoices, nil
}