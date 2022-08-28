package cmd

import (
	"github.com/spf13/cobra"
	"bitbucket.org/yasnikoff/shuffler/app"
)

var App app.App

var Root = &cobra.Command{
	Use:   "shuffler",
	Short: "Shuffler is a general purpose asset manager.",
	Long:  ``,
}

func init() {

	Root.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Print current directory status.",
		Long:  ``,
		Run:   Status,
	})
	Root.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "Clear removes metadata from the path",
		Long:  ``,
		Run:   Clear,
	})
}

func Run() {

	App = app.New()
	/*if err := App.Start(); err != nil {
		log.Fatal(err)
	}*/

	if err := Root.Execute(); err != nil {
		App.Logger().Println(err)
	}

}


