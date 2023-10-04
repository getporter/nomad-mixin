package main

import (
	"github.com/getporter/nomad-mixin/pkg/nomad"
	"github.com/spf13/cobra"
)

func buildUninstallCommand(m *nomad.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Execute the uninstall functionality of this nomad mixin",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Uninstall(cmd.Context())
		},
	}
	return cmd
}
