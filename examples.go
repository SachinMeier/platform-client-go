package main

import (
	"fmt"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
	platform "github.com/SachinMeier/platform-client-go/platform"
)

func Getdeposits(pc *platform.PlatformClient) {
	deposits, err := pc.GetDeposits(10, 0)
	if err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("Success: %d\n", deposits.Count())
		fmt.Printf("Amount: %d\n", deposits.Deposits[0].Amount)
	}
}

func Getdepositintents(pc *platform.PlatformClient) {
	deposits, err := pc.GetDepositInvoices(2, 0)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Success: %d\n", deposits.Count())
	fmt.Printf("Time: %d\n", deposits.DepositInvoices[0].Timestamp)
}

func Createdepositinvoice(pc *platform.PlatformClient) {
	invoice, err := pc.CreateDepositInvoice(4000, "test go api", "LN")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("%s\n", invoice.Invoice)
}

func Createdepositinvoicenegamt(pc *platform.PlatformClient) {
	invoice, err := pc.CreateDepositInvoice(-4000, "test go api", "LN")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("%s\n", invoice.Invoice)
}
func Createdepositinvoiceinvalidnetwork(pc *platform.PlatformClient) {
	invoice, err := pc.CreateDepositInvoice(21000000, "test go api", "BSV")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("%s\n", invoice.Invoice)
}

func Getaccount(pc *platform.PlatformClient) {
	summary, err := pc.AccountBalance()
	if err != nil {
		msg := err.Error()
		log.Error(msg)
	} else {
		fmt.Printf("account: %s\n", summary.Id)
		fmt.Printf("balance: %d sats\n", summary.Balance)
		fmt.Printf("avl_bal: %d sats\n", summary.AvailableBalance)
		fmt.Printf("res_bal: %d sats\n", summary.ReservedBalance())
	}
}

func Initwithdrawal(pc *platform.PlatformClient) {
	wreq := platform.NewWithdrawalRequest(
		4000,
		"lnbcrt40u1pshzfdppp5m8p6s7xzjzp70hnrpzl4hqc2qxmz5k6fh5srcutu3evl9l49fhtqdqqcqzpgxqrrsssp53xzsme4p2fs8wl9u480x66f8cuhg0rcs64ftx0hek0cm6l0kq04s9qyyssqmd9aepx79rz2z9lyt8mpe7sj0cru4qj32y2k9t402hglam8wsupjck9p2qwymd256tfm66p8v9qk74e2s4ktc79fa5thj5vm2y49a2spvejz8a",
	)
	w, err := pc.SubmitWithdrawalRequest(wreq)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Withdrawal ID: %s\n", w.Id)
	fmt.Printf("State: %s\n", w.State)
	fmt.Printf("Invoice: %s\n", w.Details.Invoice)
}

func Getwithdrawal(pc *platform.PlatformClient) {
	w, err := pc.GetWithdrawal("wd_QLXKPKMP")
	if err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("Withdrawal ID: %s\n", w.Id)
		fmt.Printf("State: %s\n", w.State)
		fmt.Printf("Invoice: %s\n", w.Details.Invoice)
	}
}

func Getincorrectwithdrawal(pc *platform.PlatformClient) {
	w, err := pc.GetWithdrawal("wd_QLXKQKM")
	if err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("Withdrawal ID: %s\n", w.Id)
		fmt.Printf("State: %s\n", w.State)
		fmt.Printf("Invoice: %s\n", w.Details.Invoice)
	}
}

func Decodeinvoice(pc *platform.PlatformClient) {
	decinv, err := pc.DecodeInvoice("lnbcrt2500u1ps3a29qpp5y6nnh6eew828r8r7473234dfx3xhmn8tht78xkxne4m2vj50pknqdpquwpc4curk03c9wlrswe78q4eyqc7d8d0cqzpgxqyz5vqsp5cgt2kg6k8yrhzxvh4ek6dgmlqak5hdkyzynhqknm34wujxqruhcs9qyyssqj2wg3g8w2cu2aj0k6x5cy322m6nl0y0t5jcckwdsnl2gzq4czag9exst73kgd4etgg94xwzry2lrgrkt39qqdsxfcws7q4tnt6nkzqspdy72xk")
	if err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("Invoice: %s\n", decinv.Invoice)
	}
}

func Decodeincorrectinvoice(pc *platform.PlatformClient) {
	decinv, err := pc.DecodeInvoice("asdf")
	if err != nil {
		log.Error(err.Error())
	} else {
		fmt.Printf("Invoice: %s\n", decinv.Invoice)
	}
}

func Ping(pc *platform.PlatformClient) {
	if pc.Ping() {
		fmt.Print("PING SUCCEEDED\n")
	} else {
		fmt.Print("PING FAILED\n")
	}
}
