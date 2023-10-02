package nomad

import (
	"context"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
)

type UninstallAction struct {
	Uninstalls []Uninstall `yaml:"uninstall"`
}
type Uninstall struct {
	Nomad Nomad `yaml:"nomad"`
}

func (m *Mixin) Uninstall(context context.Context) error {
	action, err := parseUninstallAction(m)
	if err != nil {
		return err
	}
	return m.DoAction(action)
}

func parseUninstallAction(m *Mixin) (*Nomad, error) {
	data, err := io.ReadAll(m.In)
	if err != nil {
		return nil, errors.Wrap(err, "could not read payload from STDIN")
	}
	var action UninstallAction
	err = yaml.Unmarshal(data, &action)
	if err != nil {
		return nil, err
	}
	if len(action.Uninstalls) != 1 {
		return nil, errors.Errorf("expected 1 uninstall steps, got %d", len(action.Uninstalls))
	}
	return &action.Uninstalls[0].Nomad, nil
}
