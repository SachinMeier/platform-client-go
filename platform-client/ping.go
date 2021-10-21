package platform

import (
	"fmt"
	"net/http"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

// Ping does ping pong with the API server at /
func (pc *PlatformClient) Ping() bool {
	log.Info("Ping Server")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/", pc.BaseURL), nil)
	if err != nil {
		log.Error("Internal Error")
		return false
	}
	// empty body response
	err = pc.sendRequest(req, nil)
	if err != nil {
		log.Errorf("Ping Failed: %s", err.Error())
		return false
	}
	log.Info("Ping Succeeded")
	return true
}
