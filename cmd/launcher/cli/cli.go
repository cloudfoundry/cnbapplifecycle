package cli

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/buildpacks/lifecycle/api"
	"github.com/buildpacks/lifecycle/cmd"
	cli "github.com/buildpacks/lifecycle/cmd/launcher/cli"
	"github.com/buildpacks/lifecycle/env"
	"github.com/buildpacks/lifecycle/launch"
	platform "github.com/buildpacks/lifecycle/platform/launch"
	"github.com/spf13/cobra"

	builderCli "code.cloudfoundry.org/cnbapplifecycle/cmd/builder/cli"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/errors"
	"code.cloudfoundry.org/cnbapplifecycle/pkg/log"
)

const defaultProcessType = "web"

func Execute() error {
	return launcherCmd.Execute()
}

func findLaunchProcessType(processes []launch.Process) (string, bool) {
	expectedCmd := strings.Join(os.Args[2:], "")
	for _, proc := range processes {
		command := append(proc.Command.Entries, proc.Args...)
		if expectedCmd == strings.Join(command, " ") {
			return proc.Type, false
		}
	}
	return "", true
}

var launcherCmd = &cobra.Command{
	Use:          "launcher",
	SilenceUsage: true,
	RunE: func(cobraCmd *cobra.Command, cmdArgs []string) error {
		var md launch.Metadata
		var args []string
		logger := log.NewLogger()
		defaultProc := defaultProcessType

		if _, err := toml.DecodeFile(launch.GetMetadataFilePath(cmd.EnvOrDefault(platform.EnvLayersDir, builderCli.DefaultLayersPath)), &md); err != nil {
			logger.Errorf("failed decoding, error: %s\n", err.Error())
			return errors.ErrLaunching
		}

		if err := verifyBuildpackAPIs(md.Buildpacks); err != nil {
			logger.Errorf("failed verifying buildpack API, error: %s\n", err.Error())
			return errors.ErrLaunching
		}

		var self string
		var isSidecar bool
		// Tasks are launched with a "--" prefix, all other processes are launched with "app"
		if len(os.Args) > 1 && os.Args[1] == "--" {
			self = "launcher"
			args = os.Args[2:]
			defaultProc = ""
		} else if len(os.Args) > 1 {
			self, isSidecar = findLaunchProcessType(md.Processes)
			logger.Infof("Detected process type: %q, isSidecar: %v", self, isSidecar)
			defaultProc = self
		}

		launcher := &launch.Launcher{
			DefaultProcessType: defaultProc,
			LayersDir:          cmd.EnvOrDefault(platform.EnvLayersDir, builderCli.DefaultLayersPath),
			AppDir:             cmd.EnvOrDefault(platform.EnvAppDir, builderCli.DefaultWorkspacePath),
			PlatformAPI:        api.MustParse(builderCli.PlatformAPI),
			Processes:          md.Processes,
			Buildpacks:         md.Buildpacks,
			Env:                env.NewLaunchEnv(os.Environ(), launch.ProcessDir, "/tmp/lifecycle"),
			Exec:               launch.OSExecFunc,
			ExecD:              launch.NewExecDRunner(),
			Shell:              launch.DefaultShell,
			Setenv:             os.Setenv,
		}

		if isSidecar {
			process := launch.Process{
				Type:    "sidecar",
				Command: launch.NewRawCommand([]string{os.Args[2]}),
				Args:    os.Args[2:],
				Direct:  true,
			}
			if err := launcher.LaunchProcess(self, process); err != nil {
				logger.Errorf("failed launching process %q, args: %#v, with self %q, error: %s\n", process.Command, process.Args, self, err.Error())
				return errors.ErrLaunching
			}
		} else {
			if err := launcher.Launch(self, args); err != nil {
				logger.Errorf("failed launching with self: %q, defaultProc: %q, args: %#v, error: %s\n", self, defaultProc, args, err.Error())
				return errors.ErrLaunching
			}
		}

		return nil
	},
}

func verifyBuildpackAPIs(bps []launch.Buildpack) error {
	for _, bp := range bps {
		if err := cmd.VerifyBuildpackAPI(cli.KindBuildpack, bp.ID, bp.API, cmd.DefaultLogger); err != nil {
			return err
		}
	}
	return nil
}
