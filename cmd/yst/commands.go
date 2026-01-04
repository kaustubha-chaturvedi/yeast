package main

import (
	"os"

	"github.com/kaustubha-chaturvedi/yeast/internal/cli"
	"github.com/kaustubha-chaturvedi/yeast/internal/plugins"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
	"github.com/spf13/cobra"
)

func createPluginsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugins",
		Short: "Plugin authoring commands",
	}

	createCmd := &cobra.Command{
		Use:   "create-new <name>",
		Short: "Create a new plugin skeleton",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias, _ := cmd.Flags().GetString("alias")
			domain, _ := cmd.Flags().GetString("domain")
			return plugins.CreateNew(args[0], alias, domain, "")
		},
	}
	createCmd.Flags().StringP("alias", "a", "", "Plugin alias (required)")
	createCmd.Flags().StringP("domain", "d", "", "Plugin domain (required)")
	createCmd.MarkFlagRequired("alias")
	createCmd.MarkFlagRequired("domain")

	cmd.AddCommand(
		createCmd,
		&cobra.Command{
			Use:   "build",
			Short: "Build plugin from current directory",
			RunE: func(cmd *cobra.Command, args []string) error {
				cwd, _ := os.Getwd()
				return plugins.Build(cwd)
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all installed plugins",
			RunE: func(cmd *cobra.Command, args []string) error {
				return plugins.List()
			},
		},
	)

	return cmd
}

func createAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "<alias>",
		Short: "Execute plugin command",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]
			rest := args[1:]

			if len(rest) == 0 {
				return utils.HandleErrorf("execute", "no command specified. Use 'yst %s help' for available commands", alias)
			}

			switch rest[0] {
			case "metadata":
				return cli.ShowMetadata(alias)
			case "help":
				if cli.HasHelpCommand(alias) {
					return cli.Route(alias, rest)
				}
				cmdName := ""
				if len(rest) > 1 {
					cmdName = rest[1]
				}
				return cli.ShowHelp(alias, cmdName)
			default:
				return cli.Route(alias, rest)
			}
		},
	}
}