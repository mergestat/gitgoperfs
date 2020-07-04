package gitgoperfs

import "testing"

var (
	repoPath string
)

func init() {
	// repoPath = "/Users/patrickdevivo/Desktop/linux"
	repoPath = "/Users/patrickdevivo/Desktop/go-git"
}

func BenchmarkGoGitRevWalk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := GoGitRevWalk(repoPath)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGitRevWalk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := GitCLIRevWalk(repoPath)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGit2GoRevwalk(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := Git2GoRevWalk(repoPath)
		if err != nil {
			b.Fatal(err)
		}
	}
}
