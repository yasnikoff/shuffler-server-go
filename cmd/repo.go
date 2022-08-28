package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/repo"
)


var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Repo manages shuffler projects.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Help())
	},
}

func init() {

	repoCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Init creates and saves new repo in uninitialized path.",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {
			var path string
			path, err := GetPathFromArgs(args, 0)
			if err != nil {
				retErr(err, cmd)
				return
			}
			r := repo.New()
			if err := r.SaveTo(path); err != nil {
				retErr(err, cmd)
				return
			}
			App.Logger().Printf("Initialized %q with %s", r.Path(), r)
			return
		},
	})
	repoCmd.AddCommand(&cobra.Command{
		Use:        "index",
		Short:"Index updates index of the repo with projects from the given path.",
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			repo, ok:= App.Repo().(*repo.Repo)
			if !ok{
				panic("can only use *repo.Repo") // TODO: remove dependency on repo.Repo
			}
			App.Logger().Printf("Start indexing %q\n", repo.Path())
			if err := repo.UpdateIndex(); err != nil {
				retErr(err, cmd)
				return
			}
			App.Logger().Printf("Index updated.")

		},
	})
	Root.AddCommand(repoCmd)
}



