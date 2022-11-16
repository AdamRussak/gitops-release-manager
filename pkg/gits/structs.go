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

type FlagsOptions struct {
	Output       string
	CommitHash   string
	Orgenization string
	Pat          string
	Project      string
	RepoPath     string
	GitUser      string
	GitEmail     string
	GitKeyPath   string
	DryRun       bool
	Gitpush      bool
}
