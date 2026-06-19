package main

import (
	"fmt"
	"os"

	"github.com/rahulkumarparida/roxkv/cli"
	"github.com/rahulkumarparida/roxkv/docs"
)




func main(){
 

	args := os.Args

	if args[1] == "roxkv"  {
		fmt.Println(args[0],args[1])
		cli.CLIUI()
	}
		
	
	if args[1] == "man" && args[2] == "roxkv" {
		fmt.Println(string(docs.ManualData))
	}

	if (args[1] == "about" || args[1] == "intro") && args[2] == "roxkv" {
		fmt.Println(string(docs.IntroData))
	}
	

}
