package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type DecodedInvoice struct {
	Amount  sats   `json:"amount"`
	Memo    string `json:"memo"`
	NodeId  string `json:"node_id"`
	Invoice string `json:"destination"`
}

type FeeEstimate struct {
	Amount  sats   `json:"amount"`
	Invoice string `json:"destination"`
	Fee     sats   `json:"fee"`
}

// DecodeInvoice decodes a Lightning Invoice using River Platform using `lncli decodepayreq`
func (pc *PlatformClient) DecodeInvoice(invoice string) (DecodedInvoice, error) {
	log.Infof("Query Decode Invoice %s", invoice)
	data := map[string]string{
		"destination": invoice,
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Error("JSON encoding error")
		return DecodedInvoice{}, err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/lightning/parse_invoice", pc.BaseURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Error("Internal Error Creating Request")
		return DecodedInvoice{}, err
	}

	var decoded_invoice DecodedInvoice
	err = pc.sendRequest(req, &decoded_invoice)
	if err != nil {
		log.Errorf("Invoice Decode Failed: %s", err.Error())
		return DecodedInvoice{}, err
	}
	return decoded_invoice, nil
}

// EstimateLightningFee estimates Lightning Fee of an invoice using `lncli`
func (pc *PlatformClient) EstimateLightningFee(invoice string, amount sats) (FeeEstimate, error) {
	log.Infof("Estimate fee for invoice %s", invoice)
	data := map[string]string{
		"destination": invoice,
	}

	body, err := json.Marshal(data)
	if err != nil {
		log.Error("JSON encoding error")
		return FeeEstimate{}, err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/lightning/estimate_fee/", pc.BaseURL),
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Error("Internal Error Creating Request")
		return FeeEstimate{}, err
	}

	var fee_estimate FeeEstimate
	err = pc.sendRequest(req, &fee_estimate)
	if err != nil {
		log.Errorf("Invoice Decode Failed: %s", err.Error())
		return FeeEstimate{}, err
	}
	return fee_estimate, nil
}
