package plugin

import (
	"os"
	"os/exec"

	"github.com/kaustubha-chaturvedi/yeast/internal/utils"
)

func Execute(pluginPath string, args []string) error {
	cmd := exec.Command(pluginPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return utils.HandleError("execute plugin", cmd.Run())
}

