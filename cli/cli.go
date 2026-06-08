package cli

import (
	"bufio"
	"fmt"
	"os"

	"strings"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
)
func CLIUI(){
	
	scanner := bufio.NewScanner(os.Stdin)
	
	var input string

	store := store.StoreInMemory()

	for{
		fmt.Print("roxkv > ")
		if scanner.Scan() {
			input = scanner.Text()
		}else {
			fmt.Println("Some error occured while scanning")
		}
		comandArgs := strings.Fields(input)


		if input == "q" {
			break
		}

		commands.ParseCommands(store,comandArgs)


	}
			

				

				
			
}