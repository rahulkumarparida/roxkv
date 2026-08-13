package cli

import (
	"context"

	"strings"

	"github.com/chzyer/readline"
	"github.com/rahulkumarparida/roxkv/internal/commands"
	"github.com/rahulkumarparida/roxkv/internal/logger"
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
			logger.ErrorLog("CLI: error initialising readline: " + err.Error())
			break
		}
		defer rl.Close()

		line, err := rl.Readline()
		if err != nil {
			logger.ErrorLog("CLI: error reading line: " + err.Error())
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
