// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"github.com/spf13/cobra"
)

var dirName string

const (
	dirNameHelp    = "directory to place kernels"
	dirNameCommand = "dir"
)

func KernelsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "build kernels",
	}

	cmd.AddCommand(
		initCommand(),
		listCommand(),
		addCommand(),
		removeCommand(),
		configureCommand(),
		rawConfigureCommand(),
		buildCommand(),
		fetchCommand(),
	)

	return cmd
}
