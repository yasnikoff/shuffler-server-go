package cmd

import (
	"bitbucket.org/yasnikoff/goutil/fsutil/patherrors"
	"bitbucket.org/yasnikoff/shuffler/internal/md"
	"github.com/spf13/cobra"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler"
)

func Status(cmd *cobra.Command, args []string) {
	var status string = "Can't get status. Error: %s"

	defer func() {
		App.Logger().Println(status)
	}()

	path, err := GetPathFromArgs(args, 0)

	if err != nil {
		status = fmt.Sprintf(status, err)
		return
	}

	h, err := shuffler.LoadHeader(path)
	switch {
	case h == nil:
		fallthrough
	case patherrors.Is(err, md.ErrPathNotInitialized):
		status = "Path is not initialized."
	case err != nil:
		status = fmt.Sprintf(status, err)
	default:
		status = fmt.Sprintf("%s", h)
	}

}
