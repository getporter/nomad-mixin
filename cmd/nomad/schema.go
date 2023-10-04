package main

import (
	"github.com/getporter/nomad-mixin/pkg/nomad"
	"github.com/spf13/cobra"
)

func buildSchemaCommand(m *nomad.Mixin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the json schema for the nomad mixin",
		Run: func(cmd *cobra.Command, args []string) {
			m.PrintSchema()
		},
	}
	return cmd
}
