package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// manualAnnotation is the Annotations map key under which each top-level command
// stores its brief, LLM-friendly manual text.
const manualAnnotation = "manual"

// importantCommands are the core commands documented by `g man` without --full,
// listed in the order they should be printed.
var importantCommands = []string{
	"create", "delete", "switch", "merge", "currentBranch", "list", "worktree",
}

// manCmd represents the man command
var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Prints an LLM-friendly guide on how to use this tool",
	Long: `Prints a concise, LLM-friendly guide on how to use this tool.

By default only the most important commands are documented. Use --full to print
the manual for every command. Run "g <command> --help" for the full flags and
details of any individual command.`,
	Aliases: []string{"manual"},
	Args:    cobra.NoArgs,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		full, _ := cmd.Flags().GetBool("full")
		fmt.Print(buildManual(full))
	},
}

// buildManual assembles the usage guide from the manual annotation of each
// top-level command. When full is false, only importantCommands are included.
func buildManual(full bool) string {
	var b strings.Builder

	b.WriteString("# g (gitBranchTool) — Command Guide\n\n")
	b.WriteString("`g` manages git branches by registering short, memorable aliases (and notes)\n")
	b.WriteString("for long or cryptic branch names, and can manage git worktrees. Most commands\n")
	b.WriteString("accept a branch NAME or its ALIAS interchangeably.\n")
	b.WriteString("Run `g <command> --help` for the full flags and details of any command.\n\n")

	// Index top-level commands by name for ordered lookup.
	byName := map[string]*cobra.Command{}
	for _, c := range rootCmd.Commands() {
		byName[c.Name()] = c
	}

	// Always emit the important commands first, in their defined order.
	emitted := map[string]bool{}
	for _, name := range importantCommands {
		if c, ok := byName[name]; ok {
			writeCommandManual(&b, c)
			emitted[name] = true
		}
	}

	if !full {
		b.WriteString("---\n")
		b.WriteString("Only the core commands are shown. Run `g manual --full` for the full list of commands.\n")
		return b.String()
	}

	// With --full, append every remaining documented command, sorted by name.
	// Hidden helper commands are skipped from the printed guide.
	var rest []string
	for name, c := range byName {
		if emitted[name] || c.Hidden {
			continue
		}
		if c.Annotations[manualAnnotation] == "" {
			continue
		}
		rest = append(rest, name)
	}
	sort.Strings(rest)
	for _, name := range rest {
		writeCommandManual(&b, byName[name])
	}

	return b.String()
}

// writeCommandManual writes a single command's manual block to b.
func writeCommandManual(b *strings.Builder, c *cobra.Command) {
	manual := c.Annotations[manualAnnotation]
	if manual == "" {
		return
	}
	b.WriteString("## " + c.Name())
	if len(c.Aliases) > 0 {
		b.WriteString("  (aliases: " + strings.Join(c.Aliases, ", ") + ")")
	}
	b.WriteString("\n")
	b.WriteString("Usage: g " + c.Use + "\n")
	b.WriteString(manual + "\n\n")
}

func init() {
	rootCmd.AddCommand(manCmd)
	manCmd.Flags().BoolP("full", "f", false, "Print the manual for all commands, not just the important ones")
	manCmd.Annotations = map[string]string{
		manualAnnotation: `Print this LLM-friendly usage guide. By default only the core commands are shown; use -f/--full to document every command.`,
	}
}
