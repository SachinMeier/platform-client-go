package main

import (
	"context"
	"fmt"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
	platform "github.com/SachinMeier/platform-client-go/platform-client"
)

func main() {
	ctx := context.Background()
	pc := platform.NewPlatformClient(ctx, "Au4FlGJ4DnnA6L12p4ZJEBOdikEKSqhkivOhc3pPlRs=", "acc_2CTR45VR", "TEST")
	// test_decodeinvoice(pc)
	// test_decodeincorrectinvoice(pc)
	test_ping(pc)
}

func test_getdeposits(pc *platform.PlatformClient) {
	deposits, err := pc.GetDeposits(10, 0)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Success: %d\n", deposits.Count())
	fmt.Printf("Amount: %d\n", deposits.Deposits[0].Amount)
}

func test_createdepositinvoice(pc *platform.PlatformClient) {
	invoice, err := pc.CreateDepositInvoice(4000, "test go api", "LN")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("%s\n", invoice.Destination)
}

func test_getaccount(pc *platform.PlatformClient) {
	summary, err := pc.AccountBalance()
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("account: %s\n", summary.Id)
	fmt.Printf("balance: %d sats\n", summary.Balance)
	fmt.Printf("avl_bal: %d sats\n", summary.AvailableBalance)
}

func test_initwithdrawal(pc *platform.PlatformClient) {
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
	fmt.Printf("Invoice: %s\n", w.Details.Destination)
}

func test_getwithdrawal(pc *platform.PlatformClient) {
	w, err := pc.GetWithdrawal("wd_QLXKQKMP")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Withdrawal ID: %s\n", w.Id)
	fmt.Printf("State: %s\n", w.State)
	fmt.Printf("Invoice: %s\n", w.Details.Destination)
}

func test_getincorrectwithdrawal(pc *platform.PlatformClient) {
	w, err := pc.GetWithdrawal("wd_QLXKQKM")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Withdrawal ID: %s\n", w.Id)
	fmt.Printf("State: %s\n", w.State)
	fmt.Printf("Invoice: %s\n", w.Details.Destination)
}

func test_decodeinvoice(pc *platform.PlatformClient) {
	decinv, err := pc.DecodeInvoice("lnbcrt2500u1ps3a29qpp5y6nnh6eew828r8r7473234dfx3xhmn8tht78xkxne4m2vj50pknqdpquwpc4curk03c9wlrswe78q4eyqc7d8d0cqzpgxqyz5vqsp5cgt2kg6k8yrhzxvh4ek6dgmlqak5hdkyzynhqknm34wujxqruhcs9qyyssqj2wg3g8w2cu2aj0k6x5cy322m6nl0y0t5jcckwdsnl2gzq4czag9exst73kgd4etgg94xwzry2lrgrkt39qqdsxfcws7q4tnt6nkzqspdy72xk")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Invoice: %s\n", decinv.Invoice)
}

func test_decodeincorrectinvoice(pc *platform.PlatformClient) {
	decinv, err := pc.DecodeInvoice("asdf")
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Invoice: %s\n", decinv.Invoice)
}

func test_ping(pc *platform.PlatformClient) {
	if pc.Ping() {
		fmt.Print("PING SUCCEEDED\n")

	} else {
		fmt.Print("PING FAILED\n")
	}
}
