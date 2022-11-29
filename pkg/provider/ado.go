package provider

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gitops-release-manager/pkg/core"
	"io"
	"net/http"
	"strconv"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	log "github.com/sirupsen/logrus"
)

// TODO: need to add a Get for releated work-items (per work-item)
// TODO: scann for child and other non-Parent work-items in close status
// TODO: list all work items that are closed and are linked to commited work-items
// TODO: tag all work-items with version tag
// TODO: validate work-item exist beofre tag, send warning if dosnt exist
func CreateNewAzureDevopsWorkItemTag(organization, personalAccessToken, project, version string, workItems []string) {
	commnads := BaseInfo{BaseUrl: "https://dev.azure.com/" + organization + "/" + project, BaseCreds: "Basic " + base64.StdEncoding.EncodeToString([]byte(" :"+personalAccessToken))}
	intArray := commnads.converWorkItemToInt(workItems)
	wiBatch := commnads.getWorkItemBatch(intArray)
	newTagsNeeded := checkExistingVersion(tagValidator(wiBatch), version)
	for _, u := range newTagsNeeded {
		commnads.UpdateWorkItemTag(u, version)
	}

}
func (b BaseInfo) UpdateWorkItemTag(id, version string) {
	resp := b.baseApiCall("PATCH", "/_apis/wit/workitems/"+id, fmt.Sprintf(`[{"op": "add","path": "/fields/System.Tags","value": "%s"}]`, version))
	//We Read the response body on the line below.
	if resp.StatusCode == 412 {
		log.Warningf("Work-Item N`%s already has a tag", id)
	} else if resp.StatusCode == 200 {
		log.Infof("Work-Item N`%s was taged with version %s", id, version)
	} else {
		body, err := io.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		//Convert the body to type string
		sb := string(body)
		log.Info(sb)
	}
}

func (b BaseInfo) getWorkItemBatch(ids []int) []byte {
	var intString string
	for _, i := range ids {
		if intString == "" {
			intString = intString + fmt.Sprint(i)
		} else {
			intString = intString + "," + fmt.Sprint(i)
		}
	}
	resp := b.baseApiCall("POST", "/_apis/wit/workitemsbatch", fmt.Sprintf(`{"ids": [%s],"fields": ["System.Id","System.Tags"]}`, intString))
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		return body
	} else {
		body, err := io.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		log.Warning(string("body: " + fmt.Sprint(resp.StatusCode)))
		log.Warning(string("body: " + string(body)))
		return nil
	}
}
func (b BaseInfo) isWorkItem(id string) bool {
	resp := b.baseApiCall("GET", "/_apis/wit/workitems/"+id, "")
	if resp.StatusCode == 200 {
		return true
	} else {
		body, err := io.ReadAll(resp.Body)
		core.OnErrorFail(err, "faild to read http body")
		log.Warningf("body: %s with Status code: %s"+string(body), fmt.Sprint(resp.StatusCode))
		return false
	}
}

func (b BaseInfo) converWorkItemToInt(wi []string) []int {
	var intReturn []int
	for _, i := range wi {
		if b.isWorkItem(i) {
			in, err := strconv.Atoi(i)
			core.OnErrorFail(err, "failed to convert string to int")
			intReturn = append(intReturn, in)
		}
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

func (b BaseInfo) baseApiCall(callType, apiPath, body string) *http.Response {
	payload := getPayload(body)
	client := &http.Client{}
	req, err := http.NewRequest(callType, b.BaseUrl+apiPath+"?api-version=7.0", payload)
	core.OnErrorFail(err, "faild to create http request")
	req.Header.Add("Authorization", b.BaseCreds)
	req.Header.Add("Content-Type", ContentType(callType))
	resp, err := client.Do(req)
	core.OnErrorFail(err, "faild to use http request")
	return resp
}
func getPayload(body string) *bytes.Buffer {
	if body != "" {
		return bytes.NewBuffer(tagBody(body))
	} else {
		return nil
	}
}
func ContentType(callType string) string {
	if callType == "PATCH" {
		return "application/json-patch+json"
	} else {
		return "application/json"
	}
}
