package platform

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/SachinMeier/platform-client-go/pkg/log"
)

const (
	Prod = "PROD"
	Test = "TEST"

	BaseURL_prod = "https://api.platform.river.com"
	BaseURL_test = "http://localhost:8080"
)

type PlatformClient struct {
	BaseURL    string
	credential string
	accountId  string
	HTTPClient *http.Client
	Context    context.Context
}

type ErrorResponse struct {
	Reason string `json:"error"`
}

func setHeaders(req *http.Request, credential string) {
	req.Header.Set("Content-Type", "application/json; charset-utf-8")
	req.Header.Set("Accept", "application/json; charset-utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("basic %s", credential))
}

func handleResponse(res *http.Response, response interface{}) error {
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errresp ErrorResponse
		var errmsg string
		err := json.NewDecoder(res.Body).Decode(&errresp)
		// err is not json
		if err != nil {
			body, err := io.ReadAll(res.Body)
			// body is empty
			if len(body) == 0 {
				errmsg = "Response Body Empty"
			} else if err != nil {
				// reading body failed
				errmsg = "Response Body Malformed"
				// err is raw string
			} else {
				errmsg = string(body)
			}
			// error is json: {"error": reason}
		} else {
			errmsg = errresp.Reason
		}
		errmsg = fmt.Sprintf("Error %d: %s", res.StatusCode, errmsg)
		log.Error(errmsg)
		return errors.New(errmsg)
	}

	if response != nil {
		err := json.NewDecoder(res.Body).Decode(response)
		if err != nil {
			log.Errorf("%s", err.Error())
			return err
		}
	}
	return nil
}

func createCredential(apiKey string) string {
	key := fmt.Sprintf("%s:%s", apiKey, apiKey)
	return b64.StdEncoding.EncodeToString([]byte(key))
}

func (pc *PlatformClient) sendRequest(req *http.Request, response interface{}) error {
	req = req.WithContext(pc.Context)
	setHeaders(req, pc.credential)

	res, err := pc.HTTPClient.Do(req)
	if err != nil {
		select {
		case <-pc.Context.Done():
			log.Errorf("Context Expired: %s", pc.Context.Err().Error())
			return pc.Context.Err()
		default:
			log.Error(err.Error())
			return err
		}
	}
	defer res.Body.Close()
	log.Infof("%s %s %d", req.Method, req.URL, res.StatusCode)
	return handleResponse(res, response)
}

func NewPlatformClient(ctx context.Context, apiKey string, accountId, env string) *PlatformClient {
	if ctx == nil {
		ctx = context.Background()
	}
	BaseURL := BaseURL_prod
	if env != Prod {
		BaseURL = BaseURL_test
	}
	return &PlatformClient{
		BaseURL:    BaseURL,
		accountId:  accountId,
		credential: createCredential(apiKey),
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
		Context: ctx,
	}
}
