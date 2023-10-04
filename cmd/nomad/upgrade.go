package main

import (
	"github.com/getporter/nomad-mixin/pkg/nomad"
	"github.com/spf13/cobra"
)

func buildUpgradeCommand(m *nomad.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Execute the upgrade functionality of this nomad mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Install(cmd.Context())
		},
	}
	return cmd
}
