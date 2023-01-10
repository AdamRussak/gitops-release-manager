package cmd

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
	GitBranch    string
	GitAuth      string
	DryRun       bool
	Gitpush      bool
}
