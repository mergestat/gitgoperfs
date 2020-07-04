package gitgoperfs

import (
	"context"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"

	gitlog "github.com/augmentable-dev/gitpert/pkg/gitlog"

	git2go "github.com/libgit2/git2go/v30"
)

var (
	GoGitCommit  *object.Commit
	GitCLICommit *gitlog.Commit
	Git2GoCommit *git2go.Commit
)

func GoGitRevWalk(repoPath string) error {
	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	iter, err := repo.Log(&gogit.LogOptions{
		From: headRef.Hash(),
	})
	if err != nil {
		return err
	}

	count := 0
	iter.ForEach(func(commit *object.Commit) error {
		count++
		GoGitCommit = commit
		return nil
	})

	return nil
}

func GitCLIRevWalk(repoPath string) error {
	res, err := gitlog.Exec(context.Background(), repoPath, "", nil)
	if err != nil {
		return err
	}

	count := 0
	for _, commit := range res {
		count++
		GitCLICommit = commit
	}

	return nil
}

func Git2GoRevWalk(repoPath string) error {
	repo, err := git2go.OpenRepository(repoPath)
	if err != nil {
		return err
	}
	defer repo.Free()

	headRef, err := repo.Head()
	if err != nil {
		return err
	}
	defer headRef.Free()

	revWalk, err := repo.Walk()
	if err != nil {
		return err
	}
	defer revWalk.Free()

	revWalk.Push(headRef.Target())
	revWalk.Sorting(git2go.SortTime)

	count := 0
	revWalk.Iterate(func(commit *git2go.Commit) bool {
		count++
		Git2GoCommit = commit
		return true
	})

	return nil
}
