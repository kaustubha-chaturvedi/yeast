package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kaustubha-chaturvedi/yeast/internal/plugin"
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
		if _, err := plugin.FindPlugin(os.Args[1]); err == nil {
			handlePluginCommand(os.Args[1:])
			return
		}
	}

	if err := rootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), "unknown command") && len(os.Args) > 1 {
			handlePluginCommand(os.Args[1:])
		} else {
			fmt.Printf("[ERROR] %v\n", err)
			os.Exit(1)
		}
	}
}

func handlePluginCommand(args []string) {
	if len(args) == 0 {
		fmt.Printf("[ERROR] No alias specified\n")
		os.Exit(1)
	}

	alias, rest := args[0], args[1:]
	pluginPath, err := plugin.FindPlugin(alias)
	if err != nil {
		fmt.Printf("[ERROR] Plugin not found: %s\n", alias)
		os.Exit(1)
	}

	if len(rest) == 0 {
		rest = []string{"-h"}
	}

	cmd := exec.Command(pluginPath, rest...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}