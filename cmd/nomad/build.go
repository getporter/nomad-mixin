package main

import (
	"github.com/ludfjig/nomad-mixin/pkg/nomad"
	"github.com/spf13/cobra"
)

func buildBuildCommand(m *nomad.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate Dockerfile lines for the bundle invocation image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Build(cmd.Context())
		},
	}
	return cmd
}
