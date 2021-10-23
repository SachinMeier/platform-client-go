package platform

import (
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type AccountSummary struct {
	Id               string `json:"id"`
	Balance          sats   `json:"balance"`
	AvailableBalance sats   `json:"available_balance"`
}

// ReservedBalance returns the reserved balance of an account, calculated as Balance - AvailableBalance
func (as *AccountSummary) ReservedBalance() sats {
	return as.Balance - as.AvailableBalance
}

// AccountBalance returns a summary of the account's balance and available balance
func (pc *PlatformClient) AccountBalance() (AccountSummary, error) {
	log.Info("Querying Account Balance")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/", pc.BaseURL, pc.accountId), nil)
	if err != nil {
		log.Error("Internal Error")
		return AccountSummary{}, err
	}

	var acct AccountSummary
	err = pc.sendRequest(req, &acct)
	if err != nil {
		log.Errorf("Account Query Failed: %s", err.Error())
		return AccountSummary{}, err
	}
	return acct, nil
}
