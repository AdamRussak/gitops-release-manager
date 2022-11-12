package provider

import (
	"encoding/base64"
	"giops-reelase-manager/pkg/core"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func UpdateTag(organization, personalAccessToken, project, id string) {
	organizationUrl := "https://dev.azure.com/" + organization + "/" + project
	p := base64.StdEncoding.EncodeToString([]byte(" :" + personalAccessToken))
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", organizationUrl+"/_apis/wit/workitems/"+id+"?api-version=7.0", nil)
	core.OnErrorFail(err, "faild to create http request")
	as := "Basic " + p
	log.Trace(as)
	req.Header.Add("Authorization", as)
	resp, err := client.Do(req)
	core.OnErrorFail(err, "faild to use http request")
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	core.OnErrorFail(err, "faild to read http body")
	//Convert the body to type string
	sb := string(body)
	log.Println(sb)
}
