package gits

import "github.com/go-git/go-git/v5"

type commit struct {
	Hash    string
	Comment string
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
	GitBranch    string
	DryRun       bool
	Gitpush      bool
}
type GitsOptions struct {
	Output       string
	gitInstance  *git.Repository
	CommitHash   string
	Orgenization string
	Pat          string
	Project      string
	RepoPath     string
	GitUser      string
	GitEmail     string
	GitKeyPath   string
	GitBranch    string
	DryRun       bool
	Gitpush      bool
}
