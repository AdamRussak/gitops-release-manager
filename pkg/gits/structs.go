package gits

type commit struct {
	Hash    string
	Comment string
}

type workItem struct {
	Name        string
	ServiceName string
	Hash        string
}
