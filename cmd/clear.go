package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"bitbucket.org/yasnikoff/shuffler/cmd/interact"
	"bitbucket.org/yasnikoff/shuffler/internal/md"
)

func Clear(cmd *cobra.Command, args []string){
	var err error
	defer func(){
		if err!=nil{
			App.Logger().Print(err)
		}
	}()
	path, err:=GetPathFromArgs(args, 0)
	if err!=nil{
		return
	}
	if !interact.Confirm(fmt.Sprintf("Do you really want to erase all the metadata from %q?", path)){
		return
	}
	if err:=md.Clear(path);err!=nil{
		return
	}
	App.Logger().Printf("Metadata removed from %q", path)

}
