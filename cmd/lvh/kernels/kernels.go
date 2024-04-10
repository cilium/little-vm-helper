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
		Use:     "kernels",
		Aliases: []string{"kernel", "k"},
		Short:   "build and pull kernels",
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
		catalogCommand(),
		pullCommand(),
	)

	return cmd
}
