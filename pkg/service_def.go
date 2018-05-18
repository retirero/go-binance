package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"net/http"
)



type apiService struct {
	URL    string
	APIKey string
	Signer Signer
	Logger log.Logger
	Ctx    context.Context
}

// NewAPIService creates instance of Service.
//
// If logger or ctx are not provided, NopLogger and Background context are used as default.
// You can use context for one-time request cancel (e.g. when shutting down the app).
func NewAPIService(url, apiKey string, signer Signer, logger log.Logger, ctx context.Context) Service {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return &apiService{
		URL:    url,
		APIKey: apiKey,
		Signer: signer,
		Logger: logger,
		Ctx:    ctx,
	}
}

func (as *apiService) request(method string, endpoint string, params map[string]string,
	apiKey bool, sign bool) (*http.Response, error) {
	transport := &http.Transport{}
	client := &http.Client{
		Transport: transport,
	}

	url := fmt.Sprintf("%s/%s", as.URL, endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}
	req.WithContext(as.Ctx)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	if apiKey {
		req.Header.Add("X-MBX-APIKEY", as.APIKey)
	}
	if sign {
		level.Debug(as.Logger).Log("queryString", q.Encode())
		q.Add("signature", as.Signer.Sign([]byte(q.Encode())))
		level.Debug(as.Logger).Log("signature", as.Signer.Sign([]byte(q.Encode())))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (as *apiService) handleError(textRes []byte) error {
	err := &Error{}
	level.Info(as.Logger).Log("errorResponse", textRes)
	if err := json.Unmarshal(textRes, err); err != nil {
		return errors.Wrap(err, "error unmarshal failed")
	}
	return err
}
