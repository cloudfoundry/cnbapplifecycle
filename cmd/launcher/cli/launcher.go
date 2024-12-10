package cli

import (
	"strings"

	"github.com/buildpacks/lifecycle/launch"
)

type TheLauncher interface {
	Launch(launcher *launch.Launcher, self string, cmd []string) error
	LaunchProcess(launcher *launch.Launcher, self string, process launch.Process) error
}

type LifecycleLauncher struct {
}

var _ TheLauncher = &LifecycleLauncher{}

func (l *LifecycleLauncher) Launch(launcher *launch.Launcher, self string, cmd []string) error {
	return launcher.Launch(self, cmd)
}

func (l *LifecycleLauncher) LaunchProcess(launcher *launch.Launcher, self string, process launch.Process) error {
	return launcher.LaunchProcess(self, process)
}

type FakeLifecycleLauncher struct {
	ExecutedCmd    string
	ExecutedDirect bool
	ExecutedType   string
}

var _ TheLauncher = &FakeLifecycleLauncher{}

func (l *FakeLifecycleLauncher) Launch(launcher *launch.Launcher, self string, cmd []string) error {
	l.ExecutedCmd = strings.Join(cmd, " ")
	l.ExecutedDirect = false
	return nil
}

func (l *FakeLifecycleLauncher) LaunchProcess(launcher *launch.Launcher, self string, process launch.Process) error {
	l.ExecutedCmd = strings.Join(process.Command.Entries, " ")
	l.ExecutedDirect = process.Direct
	l.ExecutedType = process.Type
	return nil
}
