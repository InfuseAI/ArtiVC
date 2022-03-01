package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "art",
	Short: "AritV is a version control system for large files",
	Example: `  # Push data to the repository
  cd /path/to/my/data
  art init s3://mybucket/path/to/repo
  art push -m "my first commit"

  # Pull data from the repository
  cd /path/to/download
  art init s3://mybucket/path/to/repo
  art pull

  # Download by quick command
  art get -o /path/to/download s3://mybucket/path/to/repo

  # Show command help
  art <command> -h

  For more information, please check https://github.com/infuseai/artiv`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.EnableCommandSorting = false
	rootCmd.SetUsageTemplate(usageTemplate)

	addCommandWithGroup(GROUP_QUICK,
		getCmd,
		putCmd,
	)

	addCommandWithGroup(GROUP_BASIC,
		initCommand,
		configCommand,
		pullCmd,
		pushCmd,
		tagCommand,
		listCommand,
		logCommand,
		diffCommand,
	)

	addCommandWithGroup("",
		versionCommand,
		docsCommand,
	)
}

func addCommandWithGroup(group string, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Annotations = map[string]string{
			"group": group,
		}
	}

	rootCmd.AddCommand(cmds...)
}

var usageTemplate = `{{- /* usage template */ -}}
{{define "command" -}}
{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end -}}
{{- end -}}
{{- /*
	Body
*/
-}}
Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}
{{if .HasAvailableSubCommands}}
{{- if not .HasParent}}
Basic Commands:{{range .Commands}}{{if (eq .Annotations.group "basic")}}{{template "command" .}}{{end}}{{end}}

Quick Commands (Download or upload without a workspace):{{range .Commands}}{{if (eq .Annotations.group "quick")}}{{template "command" .}}{{end}}{{end}}

Other Commands:{{range .Commands}}{{if not .Annotations.group}}{{template "command" .}}{{end}}{{end}}
{{- else}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
{{- end -}}
{{end}}
{{if .HasAvailableLocalFlags}}  
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
