package cmd

import (
	"gitops-release-manager/pkg/core"
	"gitops-release-manager/pkg/gits"
	"gitops-release-manager/pkg/markdown"
	"gitops-release-manager/pkg/provider"

	"github.com/spf13/cobra"
)

// ddCmd represents the dd command
var release = &cobra.Command{
	Use:     "release",
	Short:   releaseShortDisc,
	Long:    releaseLongDisc,
	PreRun:  core.ToggleDebug,
	Example: releaseReleaseDisc,
	Run: func(cmd *cobra.Command, args []string) {
		option := gits.FlagsOptions{GitAuth: o.GitAuth, GitBranch: o.GitBranch, GitUser: o.GitUser, GitEmail: o.GitEmail, GitKeyPath: o.GitKeyPath, Output: o.Output, CommitHash: o.CommitHash, Orgenization: o.Orgenization, Pat: o.Pat, Project: o.Project, RepoPath: o.RepoPath, DryRun: o.DryRun, Gitpush: o.Gitpush}
		r, commentsArray, newVersionTag, latestTag := option.MainGits()
		sortingForMD, workitemsID := markdown.SortCommitsForMD(commentsArray, option.Orgenization, option.Project, option.Pat)
		var setBool bool
		if !option.DryRun {
			provider.CreateNewAzureDevopsWorkItemTag(option.Orgenization, option.Pat, option.Project, newVersionTag, workitemsID)
			setBool, err := option.SetTag(r, newVersionTag)
			core.OnErrorFail(err, "failed to create Tag")
			if setBool {
				err = option.PushTags(r)
				core.OnErrorFail(err, "failed to push the tag")
			}
		}
		if setBool || option.DryRun {
			markdown.WriteToMD(sortingForMD, latestTag, newVersionTag, option.Output)
		}
	},
}

func init() {
	release.Flags().StringVar(&o.Output, "output", "./Report.md", "Set path to report output")
	release.Flags().StringVar(&o.CommitHash, "hash", "", "Set new TAG Hash")
	release.Flags().StringVar(&o.Orgenization, "org", "", "Set Azure DevOps orgenziation")
	release.Flags().StringVar(&o.Pat, "pat", "", "Set PAT for API calls")
	release.Flags().StringVar(&o.Project, "project", "", "Set Azure DevOps project")
	release.Flags().StringVar(&o.RepoPath, "repo-path", ".", "Set Path to Git repo root")
	release.Flags().StringVar(&o.GitBranch, "git-branch", "main", "Set Brnach to tag")
	release.Flags().StringVar(&o.GitUser, "git-user", ".", "Set userName to tag with")
	release.Flags().StringVar(&o.GitEmail, "git-email", ".", "Set email to tag with")
	release.Flags().StringVar(&o.GitKeyPath, "git-keyPath", "~/.ssh/id_rsa", "Set email to tag with")
	release.Flags().StringVar(&o.GitAuth, "auth", "https", "Set Auth type (ssh or https")
	release.Flags().BoolVar(&o.DryRun, "dry-run", false, "If true, only run a dry-run with cli output")
	release.Flags().BoolVar(&o.Gitpush, "git-push", false, "If true, only run a dry-run with cli output")
	release.MarkFlagsRequiredTogether("org", "project", "pat")
	release.MarkFlagRequired("hash")
	release.MarkFlagRequired("repo-path")
	rootCmd.AddCommand(release)
}
