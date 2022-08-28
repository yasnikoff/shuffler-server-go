package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

var projectTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag manages project's tags.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cmd.Help())
	},
}

func init() {
	projectCmd.AddCommand(projectTagCmd)

	projectTagCmd.AddCommand(
		&cobra.Command{
			Use:"add",
			Short:"add tag to the project",
			Long:"",
			Run: func(cmd *cobra.Command, args []string)  {

				pr, err:=App.CurrentProject()
				if err!=nil{
					App.Logger().Printf("can't get gurrent project: %s", err)
					return
				}

				if err := App.AddTags(pr.ID(), args...); err != nil {
					App.Logger().Printf("Can't add tags %q. error:\n%s", args, err)
				}
				App.Logger().Printf("Added tags %s", args)
			},
		})
}
