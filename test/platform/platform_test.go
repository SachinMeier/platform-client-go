package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/SachinMeier/platform-client-go/pkg/log"
	platform "github.com/SachinMeier/platform-client-go/platform"
)

func newServer(status int, response []byte) *httptest.Server {
	tps := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
				_, err := w.Write(response)
				if err != nil {
					log.Error(err.Error())
				}
			}),
	)

	return tps
}

// TestPing tests pinging a mock server
func TestPing(t *testing.T) {
	tps := newServer(http.StatusOK, []byte(""))
	defer tps.Close()

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)
	if !tpc.Ping() {
		t.Error("Ping Failed")
	}
}

// TestPingFail_NoServer runs ping against a nonexistent server
func TestPingFail_NoServer(t *testing.T) {
	tpc := platform.NewPlatformClient(
		context.Background(),
		"http://nohost:6969",
		"acc_test",
		"apisecret",
	)
	if tpc.Ping() {
		t.Error("Ping Failed")
	}
}

// TestPingFail_403 runs ping against a server that returns 403
func TestPingFail_403(t *testing.T) {
	tps := newServer(http.StatusForbidden, []byte(""))
	defer tps.Close()

	tpc := platform.NewPlatformClient(
		context.Background(),
		"http://nohost:6969",
		"acc_test",
		"apisecret",
	)
	if tpc.Ping() {
		t.Error("Ping Failed")
	}
}

// TestAccountBalance runs AccountBalance
func TestAccountBalance(t *testing.T) {
	data := platform.AccountSummary{
		Id:               "acc_satoshi",
		Balance:          21000000,
		AvailableBalance: 3092009,
	}
	resp, _ := json.Marshal(data)

	tps := newServer(http.StatusOK, []byte(resp))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"apisecret",
		"acc_test",
	)
	if _, err := tpc.AccountBalance(); err != nil {
		t.Error("GET Account Balance Failed")
	}
}

// TestAccountBalanceFail_InvalidAccountId runs against a server that simulates an invalid account_id
func TestAccountBalanceFail_InvalidAccountId(t *testing.T) {
	tps := newServer(http.StatusForbidden, []byte("Forbidden"))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)
	_, err := tpc.AccountBalance()
	if err == nil {
		t.Errorf("failed to fail")
	} else if msg := err.Error(); msg != "Error 403: Forbidden" {
		t.Errorf("Incorrect Error: %s", msg)
	}
}

// TestCreateDepositInvoice simulates a call to PlatformClient.CreateDepositInvoice
func TestCreateDepositInvoice(t *testing.T) {
	data := platform.DepositInvoice{
		Id:        "acc_satoshi",
		Invoice:   "lnbc2500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
		Network:   "LN",
		Timestamp: 1634975795000,
	}
	resp, _ := json.Marshal(data)

	tps := newServer(http.StatusOK, []byte(resp))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.CreateDepositInvoice(250000, "memo", "LN")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestCreateDepositInvoiceFail_InvalidAmount(t *testing.T) {
	tps := newServer(http.StatusInternalServerError, []byte("unable to process request"))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.CreateDepositInvoice(-250, "neg amt", "LN")
	if err == nil {
		t.Error("failed to fail")
	} else if msg := err.Error(); msg != "Error 500: unable to process request" {
		t.Errorf("Incorrect Error: %s", msg)
	}
}

func TestCreateDepositInvoiceFail_InvalidNetwork(t *testing.T) {
	tps := newServer(http.StatusInternalServerError, []byte("unable to process request"))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.CreateDepositInvoice(-250, "neg amt", "LN")
	if err == nil {
		t.Error("failed to fail")
	} else if msg := err.Error(); msg != "Error 500: unable to process request" {
		t.Errorf("Incorrect Error: %s", msg)
	}
}

func TestGetDepositInvoices(t *testing.T) {
	data := platform.DepositInvoiceList{
		DepositInvoices: []platform.DepositInvoice{
			{
				Id:        "acc_satoshi",
				Invoice:   "lnbc2500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
				Network:   "LN",
				Timestamp: 1634975794000,
			},
			{
				Id:        "acc_satoshi",
				Invoice:   "lnbc3500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
				Network:   "LN",
				Timestamp: 1634975795000,
			},
		},
		NextTimestamp: 1634975123333,
	}
	resp, _ := json.Marshal(data)

	tps := newServer(http.StatusOK, []byte(resp))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.GetDepositInvoices(2, 1634975794000)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestGetDeposits(t *testing.T) {
	data := platform.DepositList{
		Deposits: []platform.Deposit{
			{
				Id: "acc_satoshi",
				Invoice: platform.DepositInvoice{
					Id:        "acc_satoshi",
					Invoice:   "lnbc2500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
					Network:   "LN",
					Timestamp: 1634975794000,
				},
				Amount: 250000,
				Detail: platform.DepositDetail{
					Network: "LN",
					Proof:   "1c7272b3cb1d980b7701040e8afd537af886a437cd718ffdaf7c49c41171e11c",
				},
				Timestamp: 1634975794000,
			},
			{
				Id:     "acc_satoshi",
				Amount: 250000,
				Invoice: platform.DepositInvoice{
					Id:        "acc_satoshi",
					Invoice:   "lnbc3500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
					Network:   "LN",
					Timestamp: 1634975795000,
				},
				Detail: platform.DepositDetail{
					Network: "LN",
					Proof:   "1c7272b3cb1d980b7701040e8afd537af886a437cd718ffdaf7c49c41171e11c",
				},
				Timestamp: 1634975794000,
			},
		},
		NextTimestamp: 1634975123333,
	}
	resp, _ := json.Marshal(data)

	tps := newServer(http.StatusOK, []byte(resp))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.GetDeposits(2, 0)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestInitiateWithdrawal(t *testing.T) {

	data := platform.Withdrawal{
		Amount:   2100,
		Currency: "BTC",
		Details: platform.WithdrawalDetail{
			Network:  "LN",
			Invoice:  "lnbc3500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh",
			FeeLimit: 200,
		},
		State: "PENDING",
		Id:    "wd_tosilkroad",
	}

	resp, _ := json.Marshal(data)

	tps := newServer(http.StatusOK, []byte(resp))

	tpc := platform.NewPlatformClient(
		context.Background(),
		tps.URL,
		"acc_test",
		"apisecret",
	)

	_, err := tpc.InitiateWithdrawal(2100, "lnbc3500u1pvjluezsp5zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zygspp5qqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqqqsyqcyq5rqwzqfqypqdq5xysxxatsyp3k7enxv4jsxqzpu9qrsgquk0rl77nj30yxdy8j9vdx85fkpmdla2087ne0xh8nhedh8w27kyke0lp53ut353s06fv3qfegext0eh0ymjpf39tuven09sam30g4vgpfna3rh", "BTC", "LN", 200)
	if err != nil {
		t.Error(err.Error())
	}
}
