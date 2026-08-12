package app

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/rahulkumarparida/roxkv/cli"
	"github.com/rahulkumarparida/roxkv/docs"
	"github.com/rahulkumarparida/roxkv/server"
)

type Command string

const (
	CommandInteractive Command = "interactive"
	CommandServer      Command = "server"
	CommandManual      Command = "manual"
	CommandAbout       Command = "about"
	CommandIntro       Command = "intro"
	CommandHelp        Command = "help"
)

type Invocation struct {
	ExecutablePath string
	ExecutableName string
	Args           []string
	Command        Command
}

var errUnknownCommand = errors.New("unknown command")

func Run(argv []string, stdout, stderr io.Writer) int {
	invocation, err := Parse(argv)
	if err != nil {
		if errors.Is(err, errUnknownCommand) {
			_, _ = fmt.Fprintf(stderr, "roxkv: %v\n\n%s", err, usageText())
		} else {
			_, _ = fmt.Fprintf(stderr, "roxkv: %v\n", err)
		}
		return 1
	}

	switch invocation.Command {
	case CommandInteractive:
		cli.CLIUI()
		return 0
	case CommandServer:
		if err := server.Run(server.StartOptions{
			Command:        invocation.DisplayCommand(),
			ExecutablePath: invocation.ExecutablePath,
		}); err != nil {
			_, _ = fmt.Fprintf(stderr, "roxkv: %v\n", err)
			return 1
		}
		return 0
	case CommandManual:
		_, _ = fmt.Fprintln(stdout, string(docs.ManualData))
		return 0
	case CommandAbout, CommandIntro:
		_, _ = fmt.Fprintln(stdout, string(docs.IntroData))
		return 0
	case CommandHelp:
		_, _ = fmt.Fprint(stdout, usageText())
		return 0
	default:
		_, _ = fmt.Fprintf(stderr, "roxkv: unsupported command %q\n", invocation.Command)
		return 1
	}
}

func Parse(argv []string) (Invocation, error) {
	if len(argv) == 0 {
		return Invocation{}, fmt.Errorf("missing argv")
	}

	executablePath := argv[0]
	executableName := normalizeToken(trimExecutableName(executablePath))
	args := normalizedArgs(argv[1:])
	invocation := Invocation{
		ExecutablePath: executablePath,
		ExecutableName: executableName,
		Args:           args,
	}

	switch executableName {
	case "man":
		if matchesRoxKVTopic(args) {
			invocation.Command = CommandManual
			return invocation, nil
		}
	case "about":
		if matchesRoxKVTopic(args) {
			invocation.Command = CommandAbout
			return invocation, nil
		}
	case "intro":
		if matchesRoxKVTopic(args) {
			invocation.Command = CommandIntro
			return invocation, nil
		}
	}

	if len(args) == 0 {
		invocation.Command = CommandInteractive
		return invocation, nil
	}

	switch args[0] {
	case "tcp", "roxkv-tcp", "serve", "server":
		invocation.Command = CommandServer
		return invocation, nil
	case "man":
		if len(args) == 1 || matchesRoxKVTopic(args[1:]) {
			invocation.Command = CommandManual
			return invocation, nil
		}
	case "about":
		if len(args) == 1 || matchesRoxKVTopic(args[1:]) {
			invocation.Command = CommandAbout
			return invocation, nil
		}
	case "intro":
		if len(args) == 1 || matchesRoxKVTopic(args[1:]) {
			invocation.Command = CommandIntro
			return invocation, nil
		}
	case "help", "-h", "--help":
		invocation.Command = CommandHelp
		return invocation, nil
	}

	return Invocation{}, fmt.Errorf("%w: %s", errUnknownCommand, strings.Join(argv, " "))
}

func (i Invocation) DisplayCommand() string {
	switch i.Command {
	case CommandServer:
		return "roxkv tcp"
	case CommandManual:
		return "roxkv man"
	case CommandAbout:
		return "roxkv about"
	case CommandIntro:
		return "roxkv intro"
	case CommandInteractive:
		return "roxkv"
	default:
		if len(i.Args) == 0 {
			return i.ExecutableName
		}
		return i.ExecutableName + " " + strings.Join(i.Args, " ")
	}
}

func normalizedArgs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		trimmed := strings.TrimSpace(arg)
		if trimmed == "" {
			continue
		}
		out = append(out, normalizeToken(trimmed))
	}
	return out
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func trimExecutableName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext == "" {
		return base
	}
	return strings.TrimSuffix(base, ext)
}

func matchesRoxKVTopic(args []string) bool {
	return len(args) == 1 && normalizeToken(args[0]) == "roxkv"
}

func usageText() string {
	return `Usage:
  roxkv
  roxkv tcp
  roxkv man
  roxkv about
  roxkv intro

Compatibility aliases:
  roxkv roxkv-tcp
  man roxkv
  about roxkv
  intro roxkv
`
}
