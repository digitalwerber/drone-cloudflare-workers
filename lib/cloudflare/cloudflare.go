package cloudflare

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"../httpmethods"

	"../errormessages"
	"github.com/pkg/errors"
)

// API
type API struct {
	APIKey   string
	APIEmail string
	ZoneID   string
}

type APIRequestPayload struct {
	Content     interface{}
	ContentType string
}

type APIRequest struct {
	Path    string
	Payload APIRequestPayload
	Method  httpmethods.HttpRequestMethod
}

type APIError struct {
	Message string
	Code    int
}

type APIResponse struct {
	Success bool
	Result  interface{}
	Errors  []APIError
}

const (
	ErrorInvalidZoneID = "Invalid zone id: Must look like \"9a7806061c88ada191ed06f989cc3dac\""
)

// New creates a new cloudflare api client
func New(apiEmail string, apiKey string, zoneID string) (*API, error) {

	if apiEmail == "" || apiKey == "" {
		return nil, errors.New(errormessages.EmptyCredentials)
	}

	if zoneID == "" {
		return nil, errors.New(errormessages.EmptyZoneID)
	}

	if validateRouteID(zoneID) != true {
		return nil, errors.New(ErrorInvalidZoneID)
	}

	api := &API{
		APIEmail: apiEmail,
		APIKey:   apiKey,
		ZoneID:   zoneID,
	}

	response, err := api.PerformAPIRequest(&APIRequest{Path: "/user", Method: httpmethods.Get})

	if err != nil {
		return nil, err
	}

	if response.Success != true {
		return nil, errors.New(errormessages.InvalidCredentials)
	}

	return api, nil

}

func (api *API) PerformAPIRequests(requests ...*APIRequest) ([]*APIResponse, error) {

	c := make(chan *APIResponse)

	for _, request := range requests {
		go api.performChannelAPIRequest(request, c)
	}

	result := make([]*APIResponse, len(requests))

	for i, _ := range result {
		result[i] = <-c
	}

	return result, nil

}

func (api *API) performChannelAPIRequest(request *APIRequest, c chan *APIResponse) {
	res, _ := api.PerformAPIRequest(request)
	c <- res
}

func (api *API) PerformAPIRequest(request *APIRequest) (*APIResponse, error) {

	url := "https://api.cloudflare.com/client/v4" + request.Path

	//payload := strings.NewReader(string(scriptContent))

	if request == nil {
		request = &APIRequest{}
	}

	if request.Method == "" {
		request.Method = httpmethods.Get
	}

	var requestPayload io.Reader

	if request.Payload.Content != nil {

		var payloadString string

		switch request.Payload.Content.(type) {

		case string:
			payloadString = request.Payload.Content.(string)
			log.Println("Payload is String")

		default:
			payloadJSON, _ := json.Marshal(request.Payload.Content)
			payloadString = string(payloadJSON)
			log.Println("Payload is JSON")

		}

		requestPayload = strings.NewReader(payloadString)
	} else {
		requestPayload = nil
	}

	req, _ := http.NewRequest(string(request.Method), url, requestPayload)

	req.Header.Add("X-Auth-Email", api.APIEmail)
	req.Header.Add("X-Auth-Key", api.APIKey)

	if request.Payload.Content != "" {

		if request.Payload.ContentType == "" {
			request.Payload.ContentType = "application/json"
		}

		req.Header.Add("Content-Type", request.Payload.ContentType)
	}
	//

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	var response APIResponse

	json.Unmarshal([]byte(body), &response)

	return &response, nil
}

func (api *API) FormatAPIErrors(errs []APIError) error {
	var errStrs []string

	for _, err := range errs {
		errStrs = append(errStrs, string(err.Message+" (Code: "+string(err.Code)+")"))
	}

	return errors.New(strings.Join(errStrs[:], ", "))

}

func validateZoneID(zoneID string) bool {
	re := regexp.MustCompile("^[a-f0-9]{32}$")
	return re.MatchString(zoneID)
}
