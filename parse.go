// Package parse provides wrappers around the Parse API.
package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	IOS       = "ios"
	Increment = "Increment"
)

const (
	apiPrefix = `https://api.parse.com/1/`
)

type Envelope struct {
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

type InstallationMessage struct {
	Envelope       `json:",inline"`
	Badge          int      `json:"badge,omitempty"`
	Channels       []string `json:"channels,omitempty"`
	TimeZone       string   `json:"timeZone,omitempty"`
	DeviceType     string   `json:"deviceType,omitempty"`
	PushType       string   `json:"pushType,omitempty"`
	InstallationID string   `json:"timeZone,omitempty"`
	DeviceToken    string   `json:"deviceToken,omitempty"`
	AppName        string   `json:"appName,omitempty"`
	AppVersion     string   `json:"appVersion,omitempty"`
	ParseVersion   string   `json:"parseVersion,omitempty"`
	AppIdentifier  string   `json:"appIndentifier,omitempty"`
	ObjectID       string   `json:"objectId,omitempty"`
}

// Parse struct
type Parse struct {
	applicationID string
	clientKey     string
}

// PushResponse is the message replied by Parse after sending a message to the
// push endpoint.
type PushResponse struct {
	Result bool `json:"result"`
}

type Notification struct {
	Alert            string `json:"alert,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Sound            string `json:"sound,omitempty"`
	ContentAvailable string `json:"content-available"`
	Category         string `json:"category"`
}

// PushMessage is the general structure of a message that is sent to parse.
type PushMessage struct {
	Where    map[string]interface{} `json:"where,omitempty"`
	Channels []string               `json:"channels,omitempty"`
	Data     Notification           `json:"data"`
}

// New creates and returns a Parse client.
func New(applicationID, clientKey string) *Parse {
	p := &Parse{
		applicationID: applicationID,
		clientKey:     clientKey,
	}
	return p
}

func (p *Parse) send(endpoint string, content []byte) (*http.Response, error) {
	var err error
	var req *http.Request
	var res *http.Response

	endpointURL := apiPrefix + endpoint
	client := http.Client{}

	if req, err = http.NewRequest("POST", endpointURL, bytes.NewBuffer(content)); err != nil {
		return nil, err
	}

	req.Header.Add("X-Parse-Application-Id", p.applicationID)
	req.Header.Add("X-Parse-REST-API-Key", p.clientKey)
	req.Header.Add("Content-Type", "application/json")

	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	return res, nil
}

func (p *Parse) Installation(message InstallationMessage) (*InstallationMessage, error) {
	var ret InstallationMessage
	var res *http.Response
	var buf []byte
	var err error

	if buf, err = json.Marshal(message); err != nil {
		return nil, err
	}

	if res, err = p.send("/installations", buf); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if buf, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, &ret); err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return &ret, errors.New(ret.Error)
	}

	return &ret, nil
}

// Push sends a push request.
func (p *Parse) Push(message PushMessage) (*PushResponse, error) {
	var pr PushResponse
	var res *http.Response
	var buf []byte
	var err error

	if buf, err = json.Marshal(message); err != nil {
		return nil, err
	}

	if res, err = p.send("/push", buf); err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if buf, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(buf, &pr); err != nil {
		return nil, err
	}

	return &pr, nil
}
