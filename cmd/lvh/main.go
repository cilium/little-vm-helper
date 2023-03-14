// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"os"

	"github.com/cilium/little-vm-helper/cmd/lvh/images"
	"github.com/cilium/little-vm-helper/cmd/lvh/kernels"
	"github.com/cilium/little-vm-helper/cmd/lvh/runner"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "lvh",
		Short:        "little-vm-helper -- helper to build and run VMs",
		SilenceUsage: true,
	}
)

func init() {
	rootCmd.AddCommand(
		images.ImagesCommand(),
		kernels.KernelsCommand(),
		runner.RunCommand(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
