package main

import (
	"os"

	"github.com/rahulkumarparida/roxkv/cli"
)




func main(){
	// Accespt the arguments but from UI
	// Starts from command Roxkv
	// Initializes CLI
	// roxkv > SET name rahul
	// rahul
	// roxkv > GET name
	// rahul  

	args := os.Args

	if args[1] == "roxkv" {
		
		cli.CLIUI()
	
	}

	

}
