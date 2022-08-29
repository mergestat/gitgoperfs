package gitgoperfs

import (
	"context"

	gitlog "github.com/augmentable-dev/gitpert/pkg/gitlog"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	git2go "github.com/libgit2/git2go/v31"
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
	iter.ForEach(func(commit *object.Commit) error {
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
	})

	return nil
}

func GitCLIRevWalkStats(repoPath string) error {
	res, err := gitlog.Exec(context.Background(), repoPath, "", nil)
	if err != nil {
		return err
	}

	count := 0
	for _, commit := range res {
		count++
		for _, s := range commit.Stats {
			additions += s.Additions
			deletions += s.Deletions
		}
		GitCLICommitRevWalkStats = commit
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

	revWalk.Push(headRef.Target())
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
		defer diff.Free()

		stats, err := diff.Stats()
		if err != nil {
			panic(err)
		}
		defer stats.Free()

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
