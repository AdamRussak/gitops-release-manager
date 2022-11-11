package gits

import (
	"fmt"
	"giops-reelase-manager/pkg/core"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
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
