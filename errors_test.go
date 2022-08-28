package shuffler

import (
	"bitbucket.org/yasnikoff/shuffler/internal/md"
	"testing"
)

func TestIsInitialized(t *testing.T) {
	t.Parallel()
	path := "path does not exist"
	_, err:=LoadProject(path)

	if err!= md.ErrPathNotInitialized {
		t.Fatal("ErrPathInitialized is not detected.")
	}
}
