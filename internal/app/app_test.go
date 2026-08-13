package app

import "testing"

func TestParseInvocation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		argv    []string
		command Command
	}{
		{name: "interactive", argv: []string{"/usr/local/bin/roxkv"}, command: CommandInteractive},
		{name: "server subcommand", argv: []string{"./roxkv", "tcp"}, command: CommandServer},
		{name: "server compatibility", argv: []string{"./roxkv", "roxkv-tcp"}, command: CommandServer},
		{name: "manual subcommand", argv: []string{"./roxkv", "man"}, command: CommandManual},
		{name: "manual compatibility alias", argv: []string{"./man", "roxkv"}, command: CommandManual},
		{name: "about compatibility alias", argv: []string{"./about", "roxkv"}, command: CommandAbout},
		{name: "intro compatibility alias", argv: []string{"./intro.exe", "roxkv"}, command: CommandIntro},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			invocation, err := Parse(tc.argv)
			if err != nil {
				t.Fatalf("Parse(%v) returned error: %v", tc.argv, err)
			}

			if invocation.Command != tc.command {
				t.Fatalf("Parse(%v) command = %q, want %q", tc.argv, invocation.Command, tc.command)
			}
		})
	}
}

func TestParseUnknownCommand(t *testing.T) {
	t.Parallel()

	if _, err := Parse([]string{"./roxkv", "unknown"}); err == nil {
		t.Fatal("Parse should reject unknown commands")
	}
}
