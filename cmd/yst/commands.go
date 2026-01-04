package main

import (
	"github.com/kaustubha-chaturvedi/yeast/internal/plugins"
	"github.com/spf13/cobra"
)

func createPluginsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "plugins", Short: "Plugin authoring commands"}

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
			Use:   "install <author:plugin-name>",
			Short: "Install a plugin from the registry",
			Args:  cobra.ExactArgs(1),
			RunE:  func(_ *cobra.Command, args []string) error { return plugins.Install(args[0]) },
		},
		&cobra.Command{
			Use:   "list",
			Short: "List all installed plugins",
			RunE:  func(*cobra.Command, []string) error { return plugins.List() },
		},
	)
	return cmd
}