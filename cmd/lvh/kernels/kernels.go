// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"github.com/spf13/cobra"
)

var dirName string

func KernelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "build kernels",
	}

	cmd.PersistentFlags().StringVar(&dirName, "dir", "", "directory to to place kernels")
	cmd.MarkPersistentFlagRequired("dir")

	cmd.AddCommand(
		initCommand(),
		listCommand(),
		addCommand(),
		removeCommand(),
		configureCommand(),
		buildCommand(),
		fetchCommand(),
	)

	return cmd
}
