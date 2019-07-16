package cloudflare

import (
	"regexp"

	"../httpmethods"
	"github.com/pkg/errors"
)

type Route struct {
	ID      string
	Pattern string
	Script  string
}

const (
	ErrorInvalidRouteID = "Invalid route id: Must look like \"9a7806061c88ada191ed06f989cc3dac\""
)

func (api *API) GetAllRoutes() (routes []*Route, err error) {
	url := "/zones/" + api.ZoneID + "/workers/routes"
	resp, err := api.PerformAPIRequest(&APIRequest{Path: url})

	if err != nil {
		return nil, err
	}

	if resp.Success != true {
		return nil, api.FormatAPIErrors(resp.Errors)
	}

	arr := resp.Result.([]interface{})

	for _, route := range arr {
		dynRoute := route.(map[string]interface{})
		routes = append(routes, &Route{ID: dynRoute["id"].(string), Pattern: dynRoute["pattern"].(string), Script: dynRoute["script"].(string)})
	}

	return
}

func (api *API) DeleteRoute(routeID string) (*APIResponse, error) {

	if validateRouteID(routeID) != true {
		return nil, errors.New(ErrorInvalidRouteID)
	}

	url := "/zones/" + api.ZoneID + "/workers/routes/" + routeID
	return api.PerformAPIRequest(&APIRequest{Path: url, Method: httpmethods.Delete})
}

func (api *API) DeleteRoutes(routeIDs []string) ([]*APIResponse, error) {

	for _, routeID := range routeIDs {
		if validateRouteID(routeID) != true {
			return nil, errors.New(ErrorInvalidRouteID)
		}
	}

	requests := make([]*APIRequest, len(routeIDs))

	for i, _ := range requests {
		url := "/zones/" + api.ZoneID + "/workers/routes/" + routeIDs[i]
		requests[i] = &APIRequest{Path: url, Method: httpmethods.Delete}
	}

	return api.PerformAPIRequests(requests...)
}

func validateRouteID(routeID string) bool {
	re := regexp.MustCompile("^[a-f0-9]{32}$")
	return re.MatchString(routeID)
}
