package main

import (
	"os"

	"github.com/rahulkumarparida/roxkv/cli"
)




func main(){
 

	args := os.Args

	if args[1] == "roxkv" {
		
		cli.CLIUI()
	
	}

	

}
