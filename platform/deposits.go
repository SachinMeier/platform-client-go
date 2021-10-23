package platform

import (
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type Deposit struct {
	Id        string         `json:"id"`
	Invoice   DepositInvoice `json:"deposit_intent"`
	Amount    sats           `json:"amount"`
	Detail    DepositDetail  `json:"deposit_details"`
	State     string         `json:"state"`
	Timestamp int            `json:"timestamp"`
}

type DepositDetail struct {
	Network string `json:"network"`
	Proof   string `json:"proof"`
}

type DepositList struct {
	Deposits      []Deposit `json:"deposits"`
	NextTimestamp int       `json:"next_timestamp"`
}

func (dl *DepositList) Count() int {
	return len(dl.Deposits)
}

// GetDeposits returns a list of deposits (settled invoices) to River Platform
func (pc *PlatformClient) GetDeposits(limit, next_timestamp int) (DepositList, error) {
	log.Info("Querying Deposits")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/deposits", pc.BaseURL, pc.accountId), nil)
	if err != nil {
		log.Error("Internal Error")
		return DepositList{}, err
	}

	// Add query params
	query := req.URL.Query()
	query.Add("limit", fmt.Sprint(limit))
	if next_timestamp != 0 {
		query.Add("next_timestamp", fmt.Sprint(next_timestamp))
	}
	req.URL.RawQuery = query.Encode()

	var deposits DepositList
	err = pc.sendRequest(req, &deposits)
	if err != nil {
		log.Errorf("Deposits Query Failed: %s", err.Error())
		return DepositList{}, err
	}
	return deposits, nil
}
