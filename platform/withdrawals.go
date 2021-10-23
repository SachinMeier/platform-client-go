package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type WithdrawalDetail struct {
	Network  string `json:"network"`
	Invoice  string `json:"destination"`
	FeeLimit int    `json:"fee_limit"`
}

type Withdrawal struct {
	Amount   int              `json:"amount"`
	Currency string           `json:"currency"`
	Details  WithdrawalDetail `json:"withdrawal_details"`
	State    string           `json:"state"`
	Id       string           `json:"id"`
}

type WithdrawalRequest struct {
	Amount   sats
	Invoice  string
	FeeLimit sats
	Currency string `default:"BTC"`
	Network  string `default:"LN"`
}

const (
	LN              string = "LN"
	BTC             string = "BTC"
	DefaultFeeLimit sats   = 300
)

func (pc *PlatformClient) handleWithdrawalRequest(req *http.Request, err error) (Withdrawal, error) {
	if err != nil {
		log.Error("Internal Error Creating Request")
		return Withdrawal{}, err
	}

	var withdrawal Withdrawal
	err = pc.sendRequest(req, &withdrawal)
	if err != nil {
		log.Errorf("Querying Withdrawal Failed: %s", err.Error())
		return Withdrawal{}, err
	}
	return withdrawal, nil
}

// NewWithdrawalRequest returns a WithdrawalRequest object to be passed to SubmitWithdrawal
func NewWithdrawalRequest(amount sats, invoice string) *WithdrawalRequest {
	return &WithdrawalRequest{
		Amount:   amount,
		Invoice:  invoice,
		FeeLimit: DefaultFeeLimit,
		Currency: BTC,
		Network:  LN,
	}
}

// NewWithdrawalRequest returns a WithdrawalRequest object with a defined fee_limit to be passed to SubmitWithdrawal
func NewWithdrawalRequestWithFeeLimit(amount sats, invoice string, fee_limit sats) *WithdrawalRequest {
	return &WithdrawalRequest{
		Amount:   amount,
		Invoice:  invoice,
		FeeLimit: fee_limit,
		Currency: BTC,
		Network:  LN,
	}
}

// SubmitWithdrawalRequest takes a WithdrawalRequest and passes it to InitiateWithdrawal
func (pc *PlatformClient) SubmitWithdrawalRequest(wreq *WithdrawalRequest) (Withdrawal, error) {
	return pc.InitiateWithdrawal(wreq.Amount, wreq.Invoice, wreq.Currency, wreq.Network, wreq.FeeLimit)
}

// InitiateWithdrawal initiates a withdrawal from River Platform API by paying a specific invoice
func (pc *PlatformClient) InitiateWithdrawal(amount sats, invoice, currency, network string, fee_limit sats) (Withdrawal, error) {
	log.Infof("Initiating Withdrawal: %d sats to %s", amount, invoice)
	data := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
		"withdrawal_details": map[string]interface{}{
			"network":     network,
			"destination": invoice,
			"fee_limit":   fee_limit,
		},
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Error("JSON encoding error")
		return Withdrawal{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts/%s/withdrawals", pc.BaseURL, pc.accountId), bytes.NewBuffer(body))

	return pc.handleWithdrawalRequest(req, err)
}

// GetWithdrawal returns a withdrawal based on the passed withdrawal_id
func (pc *PlatformClient) GetWithdrawal(withdrawal_id string) (Withdrawal, error) {
	log.Infof("Querying Withdrawal %s", withdrawal_id)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/accounts/%s/withdrawals/%s",
			pc.BaseURL,
			pc.accountId,
			withdrawal_id),
		nil,
	)
	return pc.handleWithdrawalRequest(req, err)
}
