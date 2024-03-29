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
	GoGitCommitRevWalkStats  *object.Commit
	GitCLICommitRevWalkStats *gitlog.Commit
	Git2GoCommitRevWalkStats *git2go.Commit
	additions                int
	deletions                int
)

func GoGitRevWalkStats(repoPath string) error {
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
		stats, err := commit.Stats()
		if err != nil {
			return err
		}

		for _, s := range stats {
			additions += s.Addition
			deletions += s.Deletion
		}

		GoGitCommitRevWalkStats = commit
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func GitCLIRevWalkStats(repoPath string) error {
	iter, err := gitlog.Exec(context.Background(), repoPath, gitlog.WithStats(true))

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

			for _, s := range commit.Stats {
				additions += s.Additions
				deletions += s.Deletions
			}

			GitCLICommitRevWalkStats = commit
		}
	}

	return nil
}

func Git2GoRevWalkStats(repoPath string) error {
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
	err = revWalk.Iterate(func(commit *git2go.Commit) bool {
		defer commit.Free()
		count++

		p := commit.Parent(0)
		if p == nil {
			return true
		}
		defer p.Free()

		tree, err := commit.Tree()
		if err != nil {
			panic(err)
		}

		pTree, err := p.Tree()
		if err != nil {
			panic(err)
		}

		diff, err := repo.DiffTreeToTree(pTree, tree, &git2go.DiffOptions{})
		if err != nil {
			panic(err)
		}

		defer func() {
			if err = diff.Free(); err != nil {
				log.Fatal(err)
			}
		}()

		stats, err := diff.Stats()
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := stats.Free(); err != nil {
				log.Fatal(err)
			}
		}()

		additions += stats.Insertions()
		deletions += stats.Deletions()

		Git2GoCommitRevWalkStats = commit
		return true
	})
	if err != nil {
		return err
	}

	return nil
}
