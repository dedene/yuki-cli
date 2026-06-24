package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
)

func helpOptions() kong.HelpOptions {
	return kong.HelpOptions{
		NoExpandSubcommands: true,
	}
}

func helpPrinter() kong.HelpPrinter {
	return func(options kong.HelpOptions, ctx *kong.Context) error {
		var buf bytes.Buffer
		originalWriter := ctx.Stdout
		ctx.Stdout = &buf

		if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
			ctx.Stdout = originalWriter
			return err
		}
		ctx.Stdout = originalWriter

		help := compactCommandsSection(buf.String())
		help = colorizeHelp(help, helpColorProfile(ctx.Args))
		_, err := io.WriteString(originalWriter, help)
		return err
	}
}

type commandHelpRow struct {
	command     string
	description string
}

func compactCommandsSection(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line != "Commands:" {
			out = append(out, line)
			continue
		}

		out = append(out, line)
		i++

		rows := make([]commandHelpRow, 0)
		for i < len(lines) {
			if lines[i] == "" {
				i++
				continue
			}
			if !strings.HasPrefix(lines[i], "  ") || strings.HasPrefix(lines[i], "    ") {
				break
			}

			row := commandHelpRow{command: strings.TrimSpace(lines[i])}
			for i+1 < len(lines) && strings.HasPrefix(lines[i+1], "    ") {
				description := strings.TrimSpace(lines[i+1])
				if description != "" {
					if row.description != "" {
						row.description += " "
					}
					row.description += description
				}
				i++
			}

			rows = append(rows, row)
			i++
		}

		commandWidth := 0
		for _, row := range rows {
			if len(row.command) > commandWidth {
				commandWidth = len(row.command)
			}
		}

		for _, row := range rows {
			if row.description == "" {
				out = append(out, "  "+row.command)
				continue
			}

			leaderWidth := commandWidth - len(row.command) + 2
			if leaderWidth < 2 {
				leaderWidth = 2
			}
			out = append(out, "  "+row.command+" "+strings.Repeat("·", leaderWidth)+" "+row.description)
		}

		if i < len(lines) && lines[i] != "" {
			out = append(out, "")
		}
		i--
	}

	return strings.Join(out, "\n")
}

func helpColorProfile(args []string) termenv.Profile {
	if os.Getenv("NO_COLOR") != "" {
		return termenv.Ascii
	}

	switch strings.ToLower(strings.TrimSpace(os.Getenv("YUKI_COLOR"))) {
	case "always":
		return termenv.TrueColor
	case "never":
		return termenv.Ascii
	}

	if hasHelpArg(args, "--json") {
		return termenv.Ascii
	}

	return termenv.EnvColorProfile()
}

func hasHelpArg(args []string, name string) bool {
	for _, arg := range args {
		if arg == name || strings.HasPrefix(arg, name+"=") {
			return true
		}
	}
	return false
}

var helpSectionHeaders = map[string]bool{
	"Flags:":        true,
	"Global Flags:": true,
	"Commands:":     true,
	"Arguments:":    true,
}

const (
	helpColorUsage   = "#60a5fa"
	helpColorSection = "#a78bfa"
	helpColorCommand = "#38bdf8"
	helpColorDim     = "#9ca3af"
	helpColorDesc    = "#f8fafc"
)

func colorizeHelp(text string, profile termenv.Profile) string {
	if profile == termenv.Ascii {
		return text
	}

	out := termenv.NewOutput(io.Discard, termenv.WithProfile(profile))
	usage := func(s string) string {
		return out.String(s).Foreground(out.Color(helpColorUsage)).Bold().String()
	}
	section := func(s string) string {
		return out.String(s).Foreground(out.Color(helpColorSection)).Bold().String()
	}
	command := func(s string) string {
		return out.String(s).Foreground(out.Color(helpColorCommand)).Bold().String()
	}
	description := func(s string) string {
		return out.String(s).Foreground(out.Color(helpColorDesc)).String()
	}
	dim := func(s string) string {
		return out.String(s).Foreground(out.Color(helpColorDim)).String()
	}

	lines := strings.Split(text, "\n")
	inCommands := false

	for i, line := range lines {
		if line == "" {
			inCommands = false
			continue
		}

		if strings.HasPrefix(line, "Usage:") {
			lines[i] = usage("Usage:") + strings.TrimPrefix(line, "Usage:")
			continue
		}

		if helpSectionHeaders[line] {
			lines[i] = section(line)
			inCommands = line == "Commands:"
			continue
		}

		if inCommands && strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") {
			lines[i] = colorizeCommandHelpLine(line, command, description, dim)
			continue
		}

		if strings.HasPrefix(line, "  -") || strings.HasPrefix(line, "      --") {
			lines[i] = colorizeFlagHelpLine(line, dim)
			continue
		}
	}

	return strings.Join(lines, "\n")
}

func colorizeCommandHelpLine(line string, command, description, dim func(string) string) string {
	trimmed := strings.TrimPrefix(line, "  ")
	if trimmed == "" {
		return line
	}

	if idx := strings.Index(trimmed, " ·"); idx > 0 {
		commandPart := styleCommandPart(trimmed[:idx], command, dim)
		rest := strings.TrimSpace(trimmed[idx+1:])
		leader, desc, ok := strings.Cut(rest, " ")
		if !ok {
			return "  " + commandPart + " " + dim(rest)
		}
		return "  " + commandPart + " " + dim(leader) + " " + description(strings.TrimSpace(desc))
	}

	return "  " + styleCommandPart(trimmed, command, dim)
}

func styleCommandPart(commandPart string, command, dim func(string) string) string {
	words := strings.Fields(commandPart)
	styled := make([]string, 0, len(words))
	for _, word := range words {
		if strings.HasPrefix(word, "[") || strings.HasPrefix(word, "<") {
			styled = append(styled, dim(word))
			continue
		}
		styled = append(styled, command(word))
	}
	return strings.Join(styled, " ")
}

func colorizeFlagHelpLine(line string, dim func(string) string) string {
	trimmed := strings.TrimLeft(line, " ")
	indent := len(line) - len(trimmed)

	flagPart, descPart, ok := strings.Cut(trimmed, "  ")
	if !ok {
		return dim(line)
	}

	descPart = strings.TrimLeft(descPart, " ")
	if descPart == "" {
		return dim(line)
	}

	prefix := strings.Repeat(" ", indent) + flagPart
	spacing := len(line) - len(prefix) - len(descPart)
	if spacing < 2 {
		spacing = 2
	}
	return dim(prefix) + strings.Repeat(" ", spacing) + descPart
}
