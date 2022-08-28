package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)


var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project manages shuffler projects.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Help())
	},
}

func init() {

	Root.AddCommand(projectCmd)

	projectCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Init creates and saves new project in uninitialized path.",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {
			var path string

			path, err := GetPathFromArgs(args, 0)
			if err != nil {
				retErr(err, cmd)
			}
			if _, err := App.NewProject(path); err != nil {
				retErr(err, cmd)
			}
			return
		},
	})

}



