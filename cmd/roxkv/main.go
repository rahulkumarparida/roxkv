package main

import (
	"os"

	"github.com/rahulkumarparida/roxkv/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args, os.Stdout, os.Stderr))
}
