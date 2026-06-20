package main

import (
	"fmt"
	"os"

	"github.com/rahulkumarparida/roxkv/cli"
	"github.com/rahulkumarparida/roxkv/docs"
	"github.com/rahulkumarparida/roxkv/server"
)




func main(){

	



	args := os.Args

	if len(args) < 3 {

		if args[1] == "roxkv-cli"  {
			cli.CLIUI()
		}

		if args[1] == "roxkv-tcp" {
			server.Server()
		}
		
	}else{

		
		
		if args[1] == "roxkv" && args[2] == "man" {
			fmt.Println(string(docs.ManualData))
		}

		if args[1] == "roxkv" && (args[2] == "about" || args[2] == "intro")  {
			fmt.Println(string(docs.IntroData))
		}

	}

}




	// Use this When sending this to the bin directory
	// if len(args) < 2 {

	// 	if args[0] == "roxkv"  {
	// 		cli.CLIUI()
	// 	}
		
	// }else{

	// 	if args[0] == "roxkv" && args[1] == "tcp" {
	// 		server.Server()
	// 	}
		
	// 	if args[0] == "man" && args[1] == "roxkv" {
	// 		fmt.Println(string(docs.ManualData))
	// 	}

	// 	if (args[0] == "about" || args[0] == "intro") && args[1] == "roxkv" {
	// 		fmt.Println(string(docs.IntroData))
	// 	}

	// }