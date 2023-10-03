package nomad

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type UpgradeAction struct {
	Upgrade []Upgrade `yaml:"upgrade"`
}
type Upgrade struct {
	Nomad Nomad `yaml:"nomad"`
}

func (m *Mixin) Upgrade(context context.Context) error {
	action, err := parseUpgradeAction(m)
	if err != nil {
		return err
	}
	return m.DoAction(action)
}

func parseUpgradeAction(m *Mixin) (*Nomad, error) {
	data, err := io.ReadAll(m.In)
	if err != nil {
		return nil, errors.Wrap(err, "could not read payload from STDIN")
	}
	var action UpgradeAction
	err = yaml.Unmarshal(data, &action)
	if err != nil {
		return nil, err
	}
	if len(action.Upgrade) != 1 {
		return nil, errors.Errorf("expected 1 upgrade steps, got %d", len(action.Upgrade))
	}
	return &action.Upgrade[0].Nomad, nil
}
