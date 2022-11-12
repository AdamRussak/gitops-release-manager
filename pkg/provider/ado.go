package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"giops-reelase-manager/pkg/core"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func UpdateTag(organization, personalAccessToken, project, id, version string) {
	organizationUrl := "https://dev.azure.com/" + organization + "/" + project
	p := base64.StdEncoding.EncodeToString([]byte(" :" + personalAccessToken))
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", organizationUrl+"/_apis/wit/workitems/"+id+"?api-version=7.0", bytes.NewBuffer(tagBody(version)))
	core.OnErrorFail(err, "faild to create http request")
	as := "Basic " + p
	log.Trace(as)
	req.Header.Add("Authorization", as)
	req.Header.Add("Content-Type", "application/json-patch+json")
	resp, err := client.Do(req)
	core.OnErrorFail(err, "faild to use http request")
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	core.OnErrorFail(err, "faild to read http body")
	//Convert the body to type string
	sb := string(body)
	log.Println(sb)
}

func tagBody(newVersion string) []byte {
	payload, err := json.Marshal(map[string]interface{}{
		"op":    "add",
		"path":  "/fields/System.Tags",
		"value": newVersion,
	})
	core.OnErrorFail(err, "faild to marshel body")
	manipulate := "[" + string(payload) + "]"
	return []byte(manipulate)
}
