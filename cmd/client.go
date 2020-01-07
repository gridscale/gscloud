package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

const (
	requestBase                    = "/requests/"
	apiPaasServiceBase             = "/objects/paas/services"
	defaultAPIURL                  = "https://api.gridscale.io"
	bodyType                       = "application/json"
	requestDoneStatus              = "done"
	requestFailStatus              = "failed"
	defaultCheckRequestTimeoutSecs = 120
	defaultDelayIntervalMilliSecs  = 500
	requestUUIDHeaderParam         = "X-Request-Id"
)

type clientConfig struct {
	apiURL     string
	userUUID   string
	userToken  string
	userAgent  string
	httpClient *http.Client
}

type gsclient struct {
	cfg *clientConfig
}

func newClient(c *clientConfig) *gsclient {
	client := &gsclient{
		cfg: c,
	}
	return client
}

type request struct {
	uri    string
	method string
	body   interface{}
}

func (r requestError) Error() string {
	message := r.Description
	if message == "" {
		message = "no error message received from server"
	}
	errorMessageFormat := "Status code: %v. Error: %s. Request UUID: %s. "
	if r.StatusCode >= 500 {
		errorMessageFormat += "Please report this error along with the request UUID."
	}
	return fmt.Sprintf(errorMessageFormat, r.StatusCode, message, r.RequestUUID)
}

type requestStatus map[string]requestStatusProperties

type requestStatusProperties struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type requestError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusCode  int
	RequestUUID string
}

type targetFunc func() (bool, error)

func (r *request) execute(c gsclient, output interface{}) error {

	url := c.cfg.apiURL + r.uri
	httpClient := c.cfg.httpClient

	//Convert the body of the request to json
	jsonBody := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(jsonBody).Encode(r.body)
		if err != nil {
			return err
		}
	}

	//Add authentication headers and content type
	request, err := http.NewRequestWithContext(context.TODO(), r.method, url, jsonBody)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", c.cfg.userAgent)
	request.Header.Add("X-Api-Client", c.cfg.userAgent)
	request.Header.Add("X-Auth-UserID", c.cfg.userUUID)
	request.Header.Add("X-Auth-Token", c.cfg.userToken)
	request.Header.Add("Content-Type", bodyType)

	var requestUUID string
	var responseBodyBytes []byte

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	requestUUID = response.Header.Get(requestUUIDHeaderParam)
	responseBodyBytes, err = ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	//if output is set
	if output != nil {
		err = json.Unmarshal(responseBodyBytes, output)
		if err != nil {
			return err
		}
	}

	return c.waitForRequestCompleted(requestUUID)
}

func (c *gsclient) waitForRequestCompleted(id string) error {

	return retryWithTimeout(func() (bool, error) {
		r := request{
			uri:    path.Join(requestBase, id),
			method: http.MethodGet,
		}
		var response requestStatus
		err := r.execute(*c, &response)
		if err != nil {
			return false, err
		}
		if response[id].Status == requestDoneStatus {
			return false, nil
		} else if response[id].Status == requestFailStatus {
			errMessage := fmt.Sprintf("request %s failed with error %s", id, response[id].Message)
			return false, errors.New(errMessage)
		}
		return true, nil
	}, time.Duration(defaultCheckRequestTimeoutSecs)*time.Second, time.Duration(defaultDelayIntervalMilliSecs)*time.Millisecond)
}

func retryWithTimeout(target targetFunc, timeout, delay time.Duration) error {
	timer := time.After(timeout)
	var err error
	var continueRetrying bool
	for {
		select {
		case <-timer:
			if err != nil {
				return err
			}
			return errors.New("timeout reached")
		default:
			time.Sleep(delay) //delay between retries
			continueRetrying, err = target()
			if !continueRetrying {
				return err
			}
		}
	}
}
