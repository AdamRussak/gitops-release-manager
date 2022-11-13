package cmd

import (
	"giops-reelase-manager/pkg/gits"

	"github.com/spf13/cobra"
)

// ddCmd represents the dd command
var apikey string
var appkey string
var release = &cobra.Command{
	Use:   "release",
	Short: "Send Metrics to Data Dog",
	Run: func(cmd *cobra.Command, args []string) {
		option := gits.FlagsOptions{Output: o.Output, CommitHash: o.CommitHash, Orgenization: o.Orgenization, Pat: o.Pat, Project: o.Project, RepoPath: o.RepoPath, DryRun: o.DryRun, Gitpush: o.Gitpush}
		option.MainGits()

	},
}

func init() {
	release.Flags().StringVar(&o.Output, "output", "./Report.md", "Set path to report output")
	release.Flags().StringVar(&o.CommitHash, "hash", "", "Set new TAG Hash")
	release.Flags().StringVar(&o.Orgenization, "org", "", "Set Azure DevOps orgenziation")
	release.Flags().StringVar(&o.Pat, "pat", "", "Set PAT for API calls")
	release.Flags().StringVar(&o.Project, "project", "", "Set Azure DevOps project")
	release.Flags().StringVar(&o.RepoPath, "repo-path", ".", "Set Path to Git repo root")
	release.Flags().BoolVar(&o.DryRun, "dry-run", false, "If true, only run a dry-run with cli output")
	release.Flags().BoolVar(&o.Gitpush, "git-push", false, "If true, only run a dry-run with cli output")
	rootCmd.AddCommand(release)
}
