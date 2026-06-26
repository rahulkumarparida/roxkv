package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"strings"

	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/utils"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)
func CLIUI(){
	ctx , cancel := context.WithCancel(context.Background())
	scanner := bufio.NewScanner(os.Stdin)
	defer cancel()
	var input string

	store := store.StoreInMemory()

	for{
		worker.ExpiryWorker(ctx,store)
		fmt.Print("roxkv > ")
		
		
		if scanner.Scan() {
			input = scanner.Text()
		}else {
			fmt.Println("Some error occured while scanning")
			cancel()
		}
		comandArgs := strings.Fields(input)


		if input == "q" {
			cancel()
			break
		}
		nullclient := utils.NewClient{}
		commands.ParseCommands(store,comandArgs,&nullclient)


	}
			

				

				
			
}