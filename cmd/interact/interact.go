// Package interact provides functions to interact with user in command line.
package interact

import (
	"fmt"
	"strings"
)

func Confirm(question string)bool{
	var ans string
	usage:="(y/n):"
	fmt.Printf("%s (y/n):",question)

	for{
		fmt.Scan(&ans)
		ans = strings.ToLower(ans)
		if ans=="y" || ans=="yes"{
			return true
		}
		if ans=="n" || ans=="no"{
			return false
		}
		fmt.Printf("Invalid input. Please enter %s", usage)
	}

}
