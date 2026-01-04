package main

import (
	"os"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "yst",
		Short: "YEAST - Plugin host for media tools",
		Long:  "YEAST is a local-only CLI tool that works as a plugin host.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(createPluginsCmd())
	rootCmd.AddCommand(createAliasCmd())

	if err := rootCmd.Execute(); err != nil {
		utils.Logf("[ERROR] Command execution failed: %v", err)
		utils.Close()
		os.Exit(1)
	}
	utils.Close()
}