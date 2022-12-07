package provider

type BatchWorkItems struct {
	Count int        `json:"count,omitempty"`
	Value []WorkItem `json:"value,omitempty"`
}
type Payload struct {
	Ids    []int    `json:"ids"`
	Fields []string `json:"fields"`
}
type BaseInfo struct {
	BaseUrl   string
	BaseCreds string
}

type WorkItem struct {
	ID     int      `json:"id,omitempty"`
	Rev    int      `json:"rev,omitempty"`
	Fields WiFields `json:"fields,omitempty"`
	URL    string   `json:"url,omitempty"`
}

type WiFields struct {
	SystemID           int    `json:"System.Id,omitempty"`
	SystemTags         string `json:"System.Tags,omitempty"`
	SystemTitle        string `json:"System.Title,omitempty"`
	SystemWorkItemType string `json:"System.WorkItemType,omitempty"`
}
