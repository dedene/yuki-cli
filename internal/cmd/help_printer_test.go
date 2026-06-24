package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestCompactCommandsSectionUsesInlineRows(t *testing.T) {
	input := strings.Join([]string{
		"Usage: yuki <command> [flags]",
		"",
		"Commands:",
		"  a [flags]",
		"    Short command.",
		"  longer [flags]",
		"    Longer command.",
		"",
		"Flags:",
		"  -h, --help      Show help.",
		"",
	}, "\n")

	got := compactCommandsSection(input)
	want := strings.Join([]string{
		"Usage: yuki <command> [flags]",
		"",
		"Commands:",
		"  a [flags] ······· Short command.",
		"  longer [flags] ·· Longer command.",
		"",
		"Flags:",
		"  -h, --help      Show help.",
		"",
	}, "\n")

	if got != want {
		t.Fatalf("compactCommandsSection() =\n%s\nwant:\n%s", got, want)
	}
}

func TestColorizeHelpCanBeDisabled(t *testing.T) {
	input := sampleHelpText()

	if got := colorizeHelp(input, termenv.Ascii); got != input {
		t.Fatalf("colorizeHelp(ascii) = %q, want unchanged %q", got, input)
	}
}

func TestColorizeHelpStylesImportantHelpParts(t *testing.T) {
	got := colorizeHelp(sampleHelpText(), termenv.TrueColor)

	for _, want := range []string{"\x1b[", "Usage:", "Commands:", "auth", "Manage authentication.", "Show help."} {
		if !strings.Contains(got, want) {
			t.Fatalf("colorized help missing %q in:\n%s", want, got)
		}
	}
}

func TestExecuteHelpIsCompactAndColorCanBeForced(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("YUKI_COLOR", "always")

	var out bytes.Buffer
	err := Execute(context.Background(), []string{"--help"}, Runtime{Out: &out})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	got := out.String()
	for _, want := range []string{"\x1b[", "auth", "·", "Manage authentication."} {
		if !strings.Contains(got, want) {
			t.Fatalf("help output missing %q in:\n%s", want, got)
		}
	}
}

func TestExecuteHelpHonorsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("YUKI_COLOR", "always")

	var out bytes.Buffer
	err := Execute(context.Background(), []string{"--help"}, Runtime{Out: &out})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	got := out.String()
	if strings.Contains(got, "\x1b[") {
		t.Fatalf("help output contains ANSI escapes despite NO_COLOR:\n%s", got)
	}
	if !strings.Contains(got, "auth <command> [flags]") || !strings.Contains(got, "·") {
		t.Fatalf("plain help output is not compact:\n%s", got)
	}
}

func sampleHelpText() string {
	return strings.Join([]string{
		"Usage: yuki <command> [flags]",
		"",
		"Commands:",
		"  auth <command> [flags] ·· Manage authentication.",
		"",
		"Flags:",
		"  -h, --help      Show help.",
		"",
	}, "\n")
}
