package manager

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func compareOps(t *testing.T, expNames []string, actual []*Operation) {
	actualNames := opNames(actual)
	chunks := diff.DiffChunks(expNames, actualNames)
	if !isEqual(chunks) {
		printDiff(chunks)
		t.Fatal("histories not equal")
	}

}

func opNames(ops []*Operation) []string {
	var out []string
	for _, op := range ops {
		out = append(out, op.Name)
	}
	return out
}

func printDiff(chunks []diff.Chunk) {
	for _, c := range chunks {
		for _, line := range c.Added {
			fmt.Printf("+ %s\n", line)
		}
		for _, line := range c.Deleted {
			fmt.Printf("- %s\n", line)
		}
		for _, line := range c.Equal {
			fmt.Printf("  %s\n", line)
		}
	}
}

func isEqual(chunks []diff.Chunk) bool {
	if len(chunks) == 1 {
		chunk := chunks[0]
		return len(chunk.Equal) == 1 && len(chunk.Added) == 0 && len(chunk.Deleted) == 0
	}
	return false
}
