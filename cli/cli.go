package cli

import (
	"context"
	"fmt"

	"strings"

	"github.com/chzyer/readline"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/store"
	"github.com/rahulkumarparida/roxkv/internal/worker"
)

func CLIUI() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	stre, _ := store.StoreInMemory()

	for {
		worker.ExpiryWorker(ctx, stre)

		rl, err := readline.New("roxkv> ")
		if err != nil {
			fmt.Println("Error while initialising readline")
			break
		}
		defer rl.Close()

		line, err := rl.Readline()
		if err != nil {
			fmt.Println("Error while reading data")
			continue
		}

		// This is where we parse the data
		comandArgs := strings.Fields(line)

		if line == "q" {
			break
		}

		commands.ParseCommands(stre, comandArgs, nil)

	}

}
