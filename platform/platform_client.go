package platform

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

type sats int

type PlatformClient struct {
	BaseURL    string
	credential string
	accountId  string
	HTTPClient *http.Client
	Context    context.Context
}

// setHeaders sets the headers for all HTTP requests
func setHeaders(req *http.Request, credential string) {
	req.Header.Set("Content-Type", "application/json; charset-utf-8")
	req.Header.Set("Accept", "application/json; charset-utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", credential))
}

// handleResponse handles HTTP responses and unmarshals JSON to the appropriate object
func handleResponse(res *http.Response, response interface{}) error {
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		// var errresp ErrorResponse
		var errmsg string
		body, err := io.ReadAll(res.Body)
		log.Errorf(string(body))
		// err is not json
		if err != nil {
			// body is empty
			if len(body) == 0 {
				errmsg = "[Response body is empty]"
			}
		} else {
			errmsg = string(body)
		}
		errmsg = fmt.Sprintf("Error %d: %s", res.StatusCode, errmsg)
		log.Error(errmsg)
		return errors.New(errmsg)
	}

	if response != nil {
		err := json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			msg, err := io.ReadAll(res.Body)
			log.Errorf("%s", string(msg))
			return err
		}
	}
	return nil
}

// createCredential creates the basic auth credential used to authenticate requests to Platform API
func createCredential(apiKey string) string {
	key := fmt.Sprintf("%s:%s", apiKey, apiKey)
	return b64.StdEncoding.EncodeToString([]byte(key))
}

// sendRequest handles sending HTTP requests
func (pc *PlatformClient) sendRequest(req *http.Request, response interface{}) error {
	req = req.WithContext(pc.Context)
	setHeaders(req, pc.credential)

	res, err := pc.HTTPClient.Do(req)
	if err != nil {
		select {
		case <-pc.Context.Done():
			// log.Errorf("Context Expired: %s", pc.Context.Err().Error())
			return pc.Context.Err()
		default:
			// log.Error(err.Error())
			return err
		}
	}
	defer res.Body.Close()
	// log.Infof("%s %s %d", req.Method, req.URL, res.StatusCode)
	return handleResponse(res, response)
}

// NewPlatformClient creates a new PlatformClient
func NewPlatformClient(ctx context.Context, baseUrl, accountId, apiKey string) *PlatformClient {
	if ctx == nil {
		ctx = context.Background()
	}

	return &PlatformClient{
		BaseURL:    baseUrl,
		accountId:  accountId,
		credential: createCredential(apiKey),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		Context: ctx,
	}
}

// LoadEnv reads the necessary variables for creating a PlatformClient from environment and returns them
func LoadEnv() (string, string, string, error) {
	env := os.Getenv("PLATFORM_ENV")
	baseUrl := os.Getenv(fmt.Sprintf("%s_URL", env))
	accountId := os.Getenv(fmt.Sprintf("%s_RIVER_ACCOUNT_ID", env))
	apiKey := os.Getenv(fmt.Sprintf("%s_RIVER_API_SECRET", env))
	if baseUrl == "" {
		errmsg := fmt.Sprintf("%s_URL not set", env)
		log.Error(errmsg)
		return "", "", "", errors.New(errmsg)
	}
	if accountId == "" {
		errmsg := fmt.Sprintf("%s_RIVER_ACCOUNT_ID not set", env)
		log.Error(errmsg)
		return "", "", "", errors.New(errmsg)
	}
	if apiKey == "" {
		errmsg := fmt.Sprintf("%s_RIVER_API_SECRET not set", env)
		log.Error(errmsg)
		return "", "", "", errors.New(errmsg)
	}
	return baseUrl, accountId, apiKey, nil
}

// NewPlatformClientFromEnv combines LoadEnv and NewPlatformClient to create a client directly from env variables
func NewPlatformClientFromEnv() *PlatformClient {
	baseUrl, accountId, apiKey, err := LoadEnv()
	if err != nil {
		log.Errorf("Failed to Load Environment Variables: %s", err.Error())
		return nil
	}
	return NewPlatformClient(nil, apiKey, accountId, baseUrl)
}
