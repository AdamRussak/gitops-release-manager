package gits

import (
	"errors"
	"fmt"
	"gitops-release-manager/pkg/core"
	"gitops-release-manager/pkg/markdown"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/sirupsen/logrus"
)

// https://github.com/src-d/go-git/issues/1101
func (c FlagsOptions) MainGits() (*git.Repository, []markdown.WorkItem, string, string) {
	r, err := git.PlainOpen(c.RepoPath)
	core.OnErrorFail(err, "faild to get git repo")
	gitsOptions := GitsOptions{Output: c.Output, GitBranch: c.GitBranch, GitUser: c.GitUser, GitEmail: c.GitEmail, GitKeyPath: c.GitKeyPath, gitInstance: r, CommitHash: c.CommitHash, Orgenization: c.Orgenization, Pat: c.Pat, Project: c.Project, RepoPath: c.RepoPath, DryRun: c.DryRun, Gitpush: c.Gitpush}
	gitsOptions.CheckOutBranch()
	latestTag := core.EvaluateVersion(gitsOptions.getTagsArray())
	newVersionTag := core.BumpVersion(latestTag)
	log.Infof("New Version is: %s", newVersionTag)
	latestTagObject, err := r.Tag(latestTag)
	core.OnErrorFail(err, "failed to get Tag Object")
	tagObjectCommit, err := r.TagObject(latestTagObject.Hash())
	core.OnErrorFail(err, "failed to get tag object commit")
	commits := gitsOptions.GetCommits(tagObjectCommit.Target, plumbing.NewHash(c.CommitHash))
	var commentsArray []markdown.WorkItem
	for _, commit := range commits {
		if IsCommitConvention(commit.Comment) {
			split := core.SplitCommitMessage(commit.Comment)
			commentsArray = append(commentsArray, markdown.WorkItem{ServiceName: split[0], Name: split[2], Hash: split[1]})
		} else {
			commentsArray = append(commentsArray, markdown.WorkItem{ServiceName: "untracked", Name: commit.Comment, Hash: ""})
		}
	}
	return r, commentsArray, newVersionTag, latestTag
	//TODO: creat validation that commit dosent have a tag already

}
func (c GitsOptions) CheckOutBranch() {
	w, err := c.gitInstance.Worktree()
	core.OnErrorFail(err, "failed to get worktree")

	// ... checking out branch
	log.Infof("git checkout %s", c.GitBranch)

	branchRefName := plumbing.NewBranchReferenceName(c.GitBranch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}
	if err := w.Checkout(&branchCoOpts); err != nil {
		log.Warningf("local checkout of branch '%s' failed, will attempt to fetch remote branch of same name.", c.GitBranch)
		log.Warning("like `git checkout <branch>` defaulting to `git checkout -b <branch> --track <remote>/<branch>`")

		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", c.GitBranch, c.GitBranch)
		err = c.fetchOrigin(mirrorRemoteBranchRefSpec)
		core.OnErrorFail(err, "failed to featch branch origin")

		err = w.Checkout(&branchCoOpts)
		core.OnErrorFail(err, "failed to checkout branch")
	}
	log.Infof("checked out branch: %s", c.GitBranch)
}

func (c GitsOptions) fetchOrigin(refSpecStr string) error {
	remote, err := c.gitInstance.Remote("origin")
	core.OnErrorFail(err, "failed in reachging Origin")

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Info("refs already up to date")
		} else {
			return errors.New(fmt.Sprintf("fetch origin failed: %v", err))
		}
	}

	return nil
}

func (c GitsOptions) GetCommits(tagHash, newVersionHash plumbing.Hash) []commit {
	var comments []commit
	until := c.getHashObject(newVersionHash)
	fromCommit := c.getHashObject(tagHash)
	from := fromCommit.Author.When.Add(time.Second * 1)
	cIter, err := c.gitInstance.Log(&git.LogOptions{Since: &from, Until: &until.Author.When})
	core.OnErrorFail(err, "fail to get commits from tag to now")
	// ... just iterates over the commits, printing it
	_ = cIter.ForEach(func(c *object.Commit) error {
		log.Tracef("found commit Hash: %s", c.Hash.String())
		comments = append(comments, commit{Hash: c.Hash.String(), Comment: c.Message})
		return nil
	})
	return comments
}

// git tag process
func tagExists(r *git.Repository, tag string) bool {
	tagFoundErr := "tag was found"
	tags, err := r.TagObjects()
	if err != nil {
		log.Errorf("get tags error: %s", err)
		return false
	}
	res := false
	err = tags.ForEach(func(t *object.Tag) error {
		if t.Name == tag {
			res = true
			return fmt.Errorf(tagFoundErr)
		}
		return nil
	})
	if err != nil && err.Error() != tagFoundErr {
		log.Errorf("iterate tags error: %s", err)
		return false
	}
	return res
}

// setting the tag with the new version in the commit
func (c FlagsOptions) SetTag(r *git.Repository, tag string) (bool, error) {
	if tagExists(r, tag) {
		log.Infof("tag %s already exists", tag)
		return false, nil
	}
	log.Infof("Set tag %s", tag)
	_, err := r.CreateTag(tag, plumbing.NewHash(c.CommitHash), &git.CreateTagOptions{
		Message: tag,
	})

	if err != nil {
		log.Errorf("create tag error: %s", err)
		return false, err
	}

	return true, nil
}

// pushing the new tag the the remote repo
func (c FlagsOptions) PushTags(r *git.Repository) error {
	auth, err := c.publicKey()
	core.OnErrorFail(err, "Failed to get the SSH")
	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
		Auth:       auth,
	}
	err = r.Push(po)
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Info("origin remote was up to date, no push done")
			return nil
		}
		log.Errorf("push to remote origin error: %s", err)
		return err
	}
	return nil
}

// gitops commit logic
// validates the commit comment is with the agreed convention
func IsCommitConvention(commit string) bool {
	isCommit := regexp.MustCompile(`\[([A-Za-z0-9]+(-[A-Za-z0-9]+)+)]\[[A-Za-z0-9]+]\[[^\]]*]`)
	return isCommit.MatchString(commit)
}

// getting the commit from a Hash
func (c GitsOptions) getHashObject(tagHash plumbing.Hash) *object.Commit {
	fromCommit, err := c.gitInstance.CommitObject(tagHash)
	core.OnErrorFail(err, "fail to get commit object for tag")
	return fromCommit
}

// loadiong the ssh for the push tag
func (c FlagsOptions) publicKey() (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	log.Debugf("path for SSH Key: %s", c.GitKeyPath)
	sshKey, err := ioutil.ReadFile(c.GitKeyPath)
	core.OnErrorFail(err, "failed to read SSH file")
	publicKey, err = ssh.NewPublicKeys("", []byte(sshKey), "")
	core.OnErrorFail(err, "fail to get publick key")
	return publicKey, err
}

// getting all the tags in the repo
func (c GitsOptions) getTagsArray() []string {
	tags, _ := c.gitInstance.TagObjects()
	var tagsArray []string
	err := tags.ForEach(func(t *object.Tag) error {
		log.Debugf("found tag %s", t.Name)
		if c.CommitHash == t.Target.String() && isVersionTag(t.Name) {
			core.OnErrorFail(errors.New(t.Name), "Already Taged with version")
		}
		tagsArray = append(tagsArray, t.Name)
		return nil
	})
	core.OnErrorFail(err, "err of ForEach tags process")
	return tagsArray
}

// validating the tag is in SemVer convention
func isVersionTag(tag string) bool {
	if core.IsSemVer(tag) {
		return true
	} else {
		return false
	}
}
