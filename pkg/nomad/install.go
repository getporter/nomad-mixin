package nomad

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type InstallAction struct {
	Install []Install `yaml:"install"`
}
type Install struct {
	Nomad Nomad `yaml:"nomad"`
}

func (m *Mixin) Install(context context.Context) error {
	action, err := parseInstallAction(m)
	if err != nil {
		return err
	}
	return m.DoAction(action)
}

func parseInstallAction(m *Mixin) (*Nomad, error) {
	data, err := io.ReadAll(m.In)
	if err != nil {
		return nil, errors.Wrap(err, "could not read payload from STDIN")
	}
	var action InstallAction
	err = yaml.Unmarshal(data, &action)
	if err != nil {
		return nil, err
	}
	if len(action.Install) != 1 {
		return nil, errors.Errorf("expected 1 installation steps, got %d", len(action.Install))
	}
	return &action.Install[0].Nomad, nil
}
