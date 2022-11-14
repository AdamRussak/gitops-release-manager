package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"giops-reelase-manager/pkg/core"
	"io/ioutil"
	"net/http"

	jsonpatch "github.com/evanphx/json-patch"
	log "github.com/sirupsen/logrus"
)

// TODO: check if tags already exist in work-item or another version tag exist
func UpdateTag(organization, personalAccessToken, project, id, version string) {
	organizationUrl := "https://dev.azure.com/" + organization + "/" + project
	p := base64.StdEncoding.EncodeToString([]byte(" :" + personalAccessToken))
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", organizationUrl+"/_apis/wit/workitems/"+id+"?api-version=7.0", bytes.NewBuffer(tagBody(version)))
	core.OnErrorFail(err, "faild to create http request")
	as := "Basic " + p
	req.Header.Add("Authorization", as)
	req.Header.Add("Content-Type", "application/json-patch+json")
	resp, err := client.Do(req)
	core.OnErrorFail(err, "faild to use http request")
	//We Read the response body on the line below.
	if resp.StatusCode == 412 {
		log.Warningf("Work-Item N`%s already has a tag", id)
	} else if resp.StatusCode == 200 {
		log.Infof("Work-Item N`%s was taged with version %s", id, version)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		//Convert the body to type string
		sb := string(body)
		log.Println(sb)
	}

}

func tagBody(newVersion string) []byte {
	apiCallPayload := fmt.Sprintf(`[{"op": "test","path": "/fields/System.Tags","value": ""},{"op": "add","path": "/fields/System.Tags","value": "%s"}]`, newVersion)
	payLoad, err := jsonpatch.DecodePatch([]byte(apiCallPayload))
	core.OnErrorFail(err, "faild to create Payload")
	p, err := json.Marshal(payLoad)
	core.OnErrorFail(err, "faild to marshel Payload")
	log.Debug(string(p))
	return []byte(p)
}
