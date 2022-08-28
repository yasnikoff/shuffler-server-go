package cmd
import (
	"bitbucket.org/yasnikoff/shuffler/fs"
	"github.com/spf13/cobra"
)

func GetPathFromArgs(args []string, pos int)(string, error){
	var path string
	if len(args)<pos+1{
		path = ""
	}else{
		path = args[pos]
	}
	return fs.NormalizePath(path)
}


func retErr(err error, cmd *cobra.Command)error{
	App.Logger().Printf("Command: %q. Error: %s\n", cmd.Name(), err)
	return err
}
