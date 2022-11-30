package markdown

type WorkItem struct {
	Name        string
	ServiceName string
	Hash        string
}

type WorkItemOutput struct {
	Hash     string
	itemID   int
	workItem string
}
