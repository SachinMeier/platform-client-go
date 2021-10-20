package main

type WithdrawalDetail struct {
	Network     string `json:"network"`
	Destination string `json:"destination"`
	FeeLimit    int    `json:"fee_limit"`
}

type Withdrawal struct {
	Amount   int              `json:"amount"`
	Currency string           `json:"currency"`
	Details  WithdrawalDetail `json:"withdrawal_details"`
	State    string           `json:"state"`
	Id       string           `json:"id"`
}

type WithdrawalListOptions struct {
	Page    int `url:"page"`
	PerPage int `url:"per_page"`
}

// func InitiateWithdrawal(amount int, invoice string, currency string, network string, fee_limit int) {

// }
