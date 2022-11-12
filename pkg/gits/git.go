package gits

import (
	"fmt"
	"giops-reelase-manager/pkg/core"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	log "github.com/sirupsen/logrus"
)

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
		Force:  false,
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

func GetCommits(r *git.Repository, tagHash plumbing.Hash) []commit {
	var comments []commit
	until := time.Now()
	fromCommit, _ := r.CommitObject(tagHash)
	cIter, _ := r.Log(&git.LogOptions{Since: &fromCommit.Author.When, Until: &until})
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
func GetWorkItem(s string) string {
	workItemRegex := regexp.MustCompile(`[0-9]+`)
	return workItemRegex.FindString(s)
}

func IsCommitConvention(commit string) bool {
	isCommit := regexp.MustCompile(`\[([a-zA-Z]+(-[a-zA-Z]+)+)]\[[A-Za-z0-9]+]\[[^\]]*]`)
	return isCommit.MatchString(commit)
}

func StringContains(s []string, e string) (bool, int) {
	for a := range s {
		if strings.Contains(s[a], e) {
			return true, a
		}
	}
	return false, 0
}
