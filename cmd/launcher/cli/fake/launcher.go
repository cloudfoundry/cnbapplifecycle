package fake

import (
	"errors"
	"strings"

	"code.cloudfoundry.org/cnbapplifecycle/cmd/launcher/cli"
	"github.com/buildpacks/lifecycle/launch"
)

type FakeLifecycleLauncher struct {
	ExecutedCmd  string
	ExecutedSelf string
}

var _ cli.TheLauncher = &FakeLifecycleLauncher{}

func (l *FakeLifecycleLauncher) Launch(launcher *launch.Launcher, self string, cmd []string) error {
	if launcher.DefaultProcessType == "" && len(cmd) == 0 {
		return errors.New("when there is no default process a command is required")
	}
	l.ExecutedCmd = strings.Join(cmd, " ")
	l.ExecutedSelf = self

	return nil
}
