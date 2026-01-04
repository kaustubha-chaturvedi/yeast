package main

import (
	"os"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/cli"
	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "yst",
		Short:         "YEAST - Plugin host for media tools",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(createPluginsCmd())

	if len(os.Args) > 1 && os.Args[1] != "plugins" {
		plugin.BuildIndex()
		if _, err := plugin.FindPlugin(os.Args[1]); err == nil {
			cli.HandlePluginCommand(os.Args[1:])
			utils.Close()
			return
		}
	}

	if err := rootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), "unknown command") && len(os.Args) > 1 {
			cli.HandlePluginCommand(os.Args[1:])
		} else {
			utils.Printf("[ERROR] %v\n", err)
			os.Exit(1)
		}
	}
	utils.Close()
}