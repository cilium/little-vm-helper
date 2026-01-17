// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"context"
	"runtime"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/spf13/cobra"
)

func buildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build <kernel>",
		Short: "build kernel",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slogger.New()
			kname := args[0]
			arch := cmd.Flag(archFlag).Value.String()
			return kernels.BuildKernel(context.Background(), log, dirName, kname, false /* TODO: add fetch flag */, arch)
		},
	}

	cmd.Flags().String(archFlag, runtime.GOARCH, archHelp)
	cmd.Flags().StringVar(&dirName, dirNameCommand, "", dirNameHelp)
	cmd.MarkFlagRequired(dirNameCommand)

	return cmd
}
