package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"./lib/cloudflare"
	"./lib/errormessages"
	"./lib/params"
)

type cloudflareWorkerPluginInput struct {
	ZoneID     string
	ApiEmail   string
	ApiKey     string
	ScriptName string
	FileName   string
	Routes     []string
}

const (
	regexID    = "^[a-f0-9]{32,64}$"
	regexEmail = "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$"
)

func main() {

	input := &cloudflareWorkerPluginInput{
		ZoneID: params.GetPluginSetting("zoneid").
			Required(errormessages.EmptyZoneID).
			Lower().
			MustMatch(regexID, errormessages.InvalidZoneID).
			Value.(string),
		ApiEmail: params.GetPluginSetting("api_email").
			Required(errormessages.EmptyCredentials).
			Lower().
			MustMatch(regexEmail, errormessages.InvalidEmail).
			Value.(string),
		ApiKey: params.GetPluginSetting("api_key").
			Required(errormessages.EmptyCredentials).
			Lower().
			MustMatch(regexID, errormessages.InvalidAPIKey).
			Value.(string),
		ScriptName: params.GetPluginSetting("script_name").
			Required(errormessages.EmptyScript).
			Lower().
			Snake().
			Value.(string),
		FileName: params.GetPluginSetting("file_name").
			Required(errormessages.EmptyScript).
			Value.(string),
		Routes: params.GetPluginSetting("routes").
			Array().
			Value.([]string),
	}

	//zoneID := "ff15f3fc733dedd649f10ea23bc1ad60"

	scriptContent, err := ioutil.ReadFile(input.FileName)

	if err != nil {
		log.Fatal(err)
	}

	api, err := cloudflare.New(input.ApiEmail, input.ApiKey, input.ZoneID)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := api.UploadScript(api.NewScript(input.ScriptName, string(scriptContent)))

	if err != nil {
		log.Fatal(err)
	}

	if resp.Success != true {
		log.Fatal("Uploading script failed;", api.FormatAPIErrors(resp.Errors))
	}

	respJSON, _ := json.Marshal(resp)

	log.Println(string(respJSON))

	fmt.Println(api.APIEmail)

	//url := "https://api.cloudflare.com/client/v4/zones/" + input.ZoneID + "/workers/scripts/" + input.ScriptName

	//payload := strings.NewReader(string(scriptContent))

	//req, _ := http.NewRequest("PUT", url, payload)

	//req.Header.Add("X-Auth-Email", input.ApiEmail)
	//req.Header.Add("X-Auth-Key", input.ApiKey)
	//req.Header.Add("Content-Type", "application/javascript")

	//res, _ := http.DefaultClient.Do(req)

	//defer res.Body.Close()
	//body, _ := ioutil.ReadAll(res.Body)

	//if res.StatusCode != 200 {
	//	log.Fatal("Unable to upload worker script: \n", string(body))
	//}

	//fmt.Println(string(body))

	if input.Routes != nil && len(input.Routes) > 0 && input.Routes[0] != "" {

		allRoutes, err := api.GetAllRoutes()

		if err != nil {
			log.Fatal(err)
		}

		payloadJSON, _ := json.Marshal(allRoutes)

		log.Println(string(payloadJSON))

		var routesToRemove []string

		for _, route := range allRoutes {
			if route.Script == input.ScriptName {
				routesToRemove = append(routesToRemove, route.ID)
			}
		}

		if len(routesToRemove) > 0 {
			log.Println("Removing old routes from cloudflare...")
			api.DeleteRoutes(routesToRemove)
		}

		for _, element := range input.Routes {
			log.Println(element)
		}
	}

}
