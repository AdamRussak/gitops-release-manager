package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"giops-reelase-manager/pkg/core"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	log "github.com/sirupsen/logrus"
)

func CreateNewAzureDevopsWorkItemTag(organization, personalAccessToken, project, version string, workItems []string) {
	intArray := converWorkItemToInt(workItems)
	wiBatch := getWorkItemBatch(organization, personalAccessToken, project, intArray)
	newTagsNeeded := checkExistingVersion(tagValidator(wiBatch), version)
	for _, u := range newTagsNeeded {
		UpdateWorkItemTag(organization, personalAccessToken, project, u, version)
	}

}
func UpdateWorkItemTag(organization, personalAccessToken, project, id, version string) {
	organizationUrl := "https://dev.azure.com/" + organization + "/" + project
	p := base64.StdEncoding.EncodeToString([]byte(" :" + personalAccessToken))
	as := "Basic " + p
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", organizationUrl+"/_apis/wit/workitems/"+id+"?api-version=7.0", bytes.NewBuffer(tagBody(fmt.Sprintf(`[{"op": "add","path": "/fields/System.Tags","value": "%s"}]`, version))))
	core.OnErrorFail(err, "faild to create http request")
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
		log.Info(sb)
	}
}

func getWorkItemBatch(organization, personalAccessToken, project string, ids []int) []byte {
	organizationUrl := "https://dev.azure.com/" + organization + "/" + project
	p := base64.StdEncoding.EncodeToString([]byte(" :" + personalAccessToken))
	as := "Basic " + p
	var intString string
	for _, i := range ids {
		if intString == "" {
			intString = intString + fmt.Sprint(i)
		} else {
			intString = intString + "," + fmt.Sprint(i)
		}

	}
	client := &http.Client{}
	payload := fmt.Sprintf(`{"ids": [%s],"fields": ["System.Id","System.Tags"]}`, intString)
	log.Debug(payload)
	req, err := http.NewRequest("POST", organizationUrl+"/_apis/wit/workitemsbatch?api-version=7.0", bytes.NewBuffer([]byte(payload)))
	core.OnErrorFail(err, "faild to create http request")
	req.Header.Add("Authorization", as)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	core.OnErrorFail(err, "faild to use http request")
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		return body
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		log.Warning(string("body: " + fmt.Sprint(resp.StatusCode)))
		log.Warning(string("body: " + string(body)))
		return nil
	}
}

func converWorkItemToInt(wi []string) []int {
	var intReturn []int
	for _, i := range wi {
		in, err := strconv.Atoi(i)
		core.OnErrorFail(err, "failed to convert string to int")
		intReturn = append(intReturn, in)
	}
	return intReturn
}

func tagBody(body string) []byte {
	log.Debug(body)
	payLoad, err := jsonpatch.DecodePatch([]byte(body))
	core.OnErrorFail(err, "faild to create Payload")
	p, err := json.Marshal(payLoad)
	core.OnErrorFail(err, "faild to marshel Payload")
	log.Debug(string(p))
	return []byte(p)
}

func tagValidator(jsonByte []byte) BatchWorkItems {
	res := BatchWorkItems{}
	json.Unmarshal(jsonByte, &res)
	return res
}
func checkExistingVersion(existingTags BatchWorkItems, newVersion string) []string {
	var workItemNeedTag []string
	for _, workItem := range existingTags.Value {
		var split []string
		if workItem.Fields.SystemTags == "" {
			log.Infof("Work-Item n`%s needs a version %s as a Tag", fmt.Sprint(workItem.ID), newVersion)
			workItemNeedTag = append(workItemNeedTag, fmt.Sprint(workItem.ID))
		} else {
			split = strings.Split(workItem.Fields.SystemTags, ";")
			var counter int
			for _, s := range split {
				if core.IsSemVer(s) {
					counter++
					log.Warningf("Work-Item n`%s already has a version %s as a Tag", fmt.Sprint(workItem.ID), s)
				}
			}
			if counter == 0 {
				log.Infof("Work-Item n`%s needs a version %s as a Tag", fmt.Sprint(workItem.ID), newVersion)
				workItemNeedTag = append(workItemNeedTag, fmt.Sprint(workItem.ID))
			}
		}
	}
	return workItemNeedTag
}
