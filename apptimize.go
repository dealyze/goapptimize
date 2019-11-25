package apptimize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

// debug should be true for the test logging
var debug = false

// apiURL is the apptimize api url
const apiURL = "https://api.apptimize.com"

// ErrBadAttributes is the error returned when too many attribute maps or structs are passed into the track method
var ErrBadAttributes = fmt.Errorf("track only supports a single map or struct of key -> value pairs")

// Apptimize is a partially implemented aptimize sdk in golang.
type Apptimize interface {
	// Variant returns the variant of the code block that should be executed for a user in an experiment
	Variant(userID, experiment string) (string, error)

	// Track tracks a significant user event within an experiment
	Track(userID, event string, attributes ...interface{}) error
}

// Config are the config parameters of the apptimze interface
type Config struct {
	// APIToken is the API Token from the Apptimize dashboard (required)
	APIToken string

	// AppName is the name of the customer application (optional)
	AppName string

	// AppVersion is the version number of the customer application (optional)
	AppVersion string
}

// New returns a new Apptimize instance
func New(cfg *Config) Apptimize {
	return &apptimize{
		apiToken:   cfg.APIToken,
		appName:    cfg.AppName,
		appVersion: cfg.AppVersion,
		os:         "go",
		osVersion:  runtime.Version(),
	}
}

// apptimize is the implementation of the Apptimize interface
type apptimize struct {
	apiToken   string
	appName    string
	appVersion string
	os         string
	osVersion  string
}

// Variant returns the variant of the code block that should be executed for a user in an experiment
func (a *apptimize) Variant(userID, experiment string) (string, error) {
	var resp struct {
		CodeBlockVariant string `json:"codeBlockVariant,omitempy"`
	}
	if err := a.request("GET", fmt.Sprintf("%s/v1/users/%s/code-blocks/%s", apiURL, userID, experiment), nil, &resp); err != nil {
		return "", err
	}
	return resp.CodeBlockVariant, nil
}

// Track tracks a significant user event within an experiment
func (a *apptimize) Track(userID, event string, attributes ...interface{}) error {
	if l := len(attributes); l > 1 {
		return ErrBadAttributes
	} else if l == 1 {
		return a.request("POST", fmt.Sprintf("%s/v1/users/%s/events/%s", apiURL, userID, event), attributes[0], nil)
	}
	return a.request("POST", fmt.Sprintf("%s/v1/users/%s/events/%s", apiURL, userID, event), nil, nil)
}

// request issues a request to the api and parses the response
func (a *apptimize) request(method, url string, reqBody, respBody interface{}) error {
	var buffer bytes.Buffer
	if reqBody == nil {
		// no request body
	} else if err := json.NewEncoder(&buffer).Encode(&reqBody); err != nil {
		logd(err)
		return err
	}
	req, err := http.NewRequest(method, url, &buffer)
	if err != nil {
		logd(err)
		return err
	}
	req.Header.Add("ApptimizeApiToken", a.apiToken)
	req.Header.Add("ApptimizeApplicationName", a.appName)
	req.Header.Add("ApptimizeApplicationVersion", a.appVersion)
	req.Header.Add("ApptimizeOperatingSystem", a.os)
	req.Header.Add("ApptimizeOperatingSystemVersion", a.osVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logd(err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err := fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
		logd(err)
		return err
	} else if respBody == nil {
		// no response body
	} else if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		logd(err)
		return err
	}
	return nil
}

// logf is a degub log
func logd(is ...interface{}) {
	if debug {
		for _, i := range is {
			log.Output(2, fmt.Sprintf("%s", i))
		}
	}
}
