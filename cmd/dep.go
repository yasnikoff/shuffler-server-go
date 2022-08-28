package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var depCmd = &cobra.Command{
	Use:   "dep",
	Short: "dep manages project's dependencies.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Help())
	},
}

var dep_update_recursive bool

func init() {

	projectCmd.AddCommand(depCmd)

	depCmd.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "add adds and saves project's dependency.",
		Long:  `Command takes two arguments:
	- dependency's subpath inside projects folder
	- dependency's ID

		`,

		Run:   func(cmd *cobra.Command, args []string) {

			if len(args) < 2 {
				cmd.Help()
			}

			pr, err := App.CurrentProject()
			if err != nil {
				retErr(err, cmd)
			}
			if err := App.AddDep(pr.ID().String(), args[0], args[1]); err != nil {
				retErr(err, cmd)
			}

			return
		},
	})
	cmd_dep_update := &cobra.Command{
		Use:   "update",
		Short: "update rewrites all project's dependencies with new versions.",
		Long:  ``,

		Run:   func(cmd *cobra.Command, args []string) {

			pr, err := App.CurrentProject();
			if err != nil {
				retErr(err, cmd)
			}
			if err := pr.Deps.UpdateAll(); err != nil {
				retErr(err, cmd)
			}
			return
		},
	}
	//cmd_dep_update.Flags().BoolVarP(&dep_update_recursive, "recursive", "r", false, "recursive update of dependencies.")
	depCmd.AddCommand(cmd_dep_update)


	depCmd.AddCommand(&cobra.Command{
		Use:   "remove",
		Short: "remove deletes project's dependency.",
		Long:  `Command takes one argument: dependency's subpath inside project folder

		`,

		Run:   func(cmd *cobra.Command, args []string) {

			if len(args) < 1 {
				cmd.Help()
			}

			pr, err := App.CurrentProject()
			if err != nil {
				retErr(err, cmd)
			}
			if err := App.RemoveDep(pr.ID().String(), args[0]); err != nil {
				retErr(err, cmd)
			}

			return
		},
	})
}
