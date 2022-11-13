package gits

import (
	"fmt"
	"giops-reelase-manager/pkg/core"
	"giops-reelase-manager/pkg/markdown"
	"os"
	"regexp"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

func (c FlagsOptions) MainGits() {

	r, err := git.PlainOpen(c.RepoPath)
	core.OnErrorFail(err, "faild to get git repo")
	CheckOutBranch(r, "main")
	tags, _ := r.TagObjects()
	var tagsArray []string
	err = tags.ForEach(func(t *object.Tag) error {
		log.Infof("found tag %s", t.Name)
		tagsArray = append(tagsArray, t.Name)
		return nil
	})
	core.OnErrorFail(err, "err of ForEach tags process")
	latestTag := core.EvaluateVersion(tagsArray)
	newVersionTag := core.BumpVersion(latestTag)
	log.Infof("New Version is: %s", newVersionTag)
	latestTagObject, err := r.Tag(latestTag)
	core.OnErrorFail(err, "failed to get Tag Object")
	tagObjectCommit, err := r.TagObject(latestTagObject.Hash())
	core.OnErrorFail(err, "failed to get tag object commit")
	commits := GetCommits(r, tagObjectCommit.Target, plumbing.NewHash(c.CommitHash))
	var commentsArray []markdown.WorkItem
	for _, commit := range commits {
		if IsCommitConvention(commit.Comment) {
			split := core.SplitCommitMessage(commit.Comment)
			commentsArray = append(commentsArray, markdown.WorkItem{ServiceName: split[0], Name: split[2], Hash: split[1]})
		} else {
			commentsArray = append(commentsArray, markdown.WorkItem{ServiceName: "untracked", Name: commit.Comment, Hash: ""})
		}
	}
	sortingForMD := markdown.SortCommitsForMD(commentsArray, c.Orgenization, c.Project, c.Pat, newVersionTag)
	markdown.WriteToMD(sortingForMD, latestTag, newVersionTag)
}
func CheckOutBranch(r *git.Repository, branch string) {
	// ... retrieving the commit being pointed by HEAD
	log.Info("git show-ref --head HEAD")
	ref, err := r.Head()
	core.OnErrorFail(err, "failed to get head")

	w, err := r.Worktree()
	core.OnErrorFail(err, "failed to get worktree")

	// ... checking out branch
	log.Info("git checkout %s", branch)

	branchRefName := plumbing.NewBranchReferenceName(branch)
	branchCoOpts := git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchRefName),
		Force:  true,
	}
	if err := w.Checkout(&branchCoOpts); err != nil {
		log.Warning("local checkout of branch '%s' failed, will attempt to fetch remote branch of same name.", branch)
		log.Warning("like `git checkout <branch>` defaulting to `git checkout -b <branch> --track <remote>/<branch>`")

		mirrorRemoteBranchRefSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branch, branch)
		err = fetchOrigin(r, mirrorRemoteBranchRefSpec)
		core.OnErrorFail(err, "failed to featch branch origin")

		err = w.Checkout(&branchCoOpts)
		core.OnErrorFail(err, "failed to checkout branch")
	}
	core.OnErrorFail(err, "failed in process")

	log.Info("checked out branch: %s", branch)

	// ... retrieving the commit being pointed by HEAD (branch now)
	log.Info("git show-ref --head HEAD")
	ref, err = r.Head()
	core.OnErrorFail(err, "failed in getting head")
	fmt.Println(ref.Hash())
}
func fetchOrigin(repo *git.Repository, refSpecStr string) error {
	remote, err := repo.Remote("origin")
	core.OnErrorFail(err, "failed in reachging Origin")

	var refSpecs []config.RefSpec
	if refSpecStr != "" {
		refSpecs = []config.RefSpec{config.RefSpec(refSpecStr)}
	}

	if err = remote.Fetch(&git.FetchOptions{
		RefSpecs: refSpecs,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			fmt.Print("refs already up to date")
		} else {
			return fmt.Errorf("fetch origin failed: %v", err)
		}
	}

	return nil
}

func GetCommits(r *git.Repository, tagHash, newVersionHash plumbing.Hash) []commit {
	var comments []commit
	until := getHashObject(r, newVersionHash)
	fromCommit := getHashObject(r, tagHash)
	from := fromCommit.Author.When.Add(time.Second * 1)
	cIter, err := r.Log(&git.LogOptions{Since: &from, Until: &until.Author.When})
	core.OnErrorFail(err, "fail to get commits from tag to now")
	// ... just iterates over the commits, printing it
	_ = cIter.ForEach(func(c *object.Commit) error {
		comments = append(comments, commit{Hash: c.Hash.String(), Comment: c.Message})
		return nil
	})
	return comments
}

// git tag process
func tagExists(tag string, r *git.Repository) bool {
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

func SetTag(r *git.Repository, tag string, tagger *object.Signature) (bool, error) {
	if tagExists(tag, r) {
		log.Infof("tag %s already exists", tag)
		return false, nil
	}
	log.Infof("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		log.Errorf("get HEAD error: %s", err)
		return false, err
	}

	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger:  tagger,
		Message: tag,
	})

	if err != nil {
		log.Errorf("create tag error: %s", err)
		return false, err
	}

	return true, nil
}
func pushTags(r *git.Repository) error {

	po := &git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")},
	}

	err := r.Push(po)

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

func IsCommitConvention(commit string) bool {
	isCommit := regexp.MustCompile(`\[([a-zA-Z]+(-[a-zA-Z]+)+)]\[[A-Za-z0-9]+]\[[^\]]*]`)
	return isCommit.MatchString(commit)
}

func getHashObject(r *git.Repository, tagHash plumbing.Hash) *object.Commit {
	fromCommit, err := r.CommitObject(tagHash)
	core.OnErrorFail(err, "fail to get commit object for tag")
	return fromCommit
}
