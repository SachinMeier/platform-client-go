package platform_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	platform "github.com/SachinMeier/platform-client-go/platform"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	tps := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte(""))
			}),
	)
	defer tps.Close()

	tpc := platform.NewPlatformClient(
		context.Background(),
		"acc_test",
		"apisecret",
		tps.URL,
	)

	assert.True(t, tpc.Ping())

}
