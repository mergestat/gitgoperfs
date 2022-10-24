package gitgoperfs

import (
	"context"
	"errors"
	"io"
	"log"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	git2go "github.com/libgit2/git2go/v34"
	gitlog "github.com/mergestat/gitutils/gitlog"
)

var (
	GoGitCommitRevWalk  *object.Commit
	GitCLICommitRevWalk *gitlog.Commit
	Git2GoCommitRevWalk *git2go.Commit
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
	if err := iter.ForEach(func(commit *object.Commit) error {
		count++
		GoGitCommitRevWalk = commit
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func GitCLIRevWalk(repoPath string) error {
	iter, err := gitlog.Exec(context.Background(), repoPath, gitlog.WithStats(false))
	if err != nil {
		return err
	}

	count := 0

	for {
		if commit, err := iter.Next(); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Fatal(err)
		} else {

			count++

			GitCLICommitRevWalk = commit
		}

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

	if err := revWalk.Push(headRef.Target()); err != nil {
		return err
	}

	revWalk.Sorting(git2go.SortTime)

	count := 0
	if err := revWalk.Iterate(func(commit *git2go.Commit) bool {
		defer commit.Free()
		count++
		Git2GoCommitRevWalk = commit
		return true
	}); err != nil {
		return err
	}

	return nil
}
