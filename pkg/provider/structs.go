package provider

type BatchWorkItems struct {
	Count int `json:"count,omitempty"`
	Value []struct {
		ID     int `json:"id,omitempty"`
		Rev    int `json:"rev,omitempty"`
		Fields struct {
			SystemID   int    `json:"System.Id,omitempty"`
			SystemTags string `json:"System.Tags,omitempty"`
		} `json:"fields,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"value,omitempty"`
}
