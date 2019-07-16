package cloudflare

import (
	"../httpmethods"
	"github.com/iancoleman/strcase"
)

type Script struct {
	ID      string
	Content string
}

func (api *API) NewScript(id string, content string) *Script {
	return &Script{ID: strcase.ToSnake(id), Content: content}
}

func (api *API) UploadScript(script *Script) (*APIResponse, error) {
	url := "/zones/" + api.ZoneID + "/workers/scripts/" + script.ID
	return api.PerformAPIRequest(&APIRequest{Path: url, Method: httpmethods.Put, Payload: APIRequestPayload{Content: script.Content, ContentType: "application/javascript"}})
}
