// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"context"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/spf13/cobra"
)

func removeCommand() *cobra.Command {
	var backupConf bool
	cmd := &cobra.Command{
		Use:   "remove <kernel>",
		Short: "remove kernel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := slogger.New()
			kname := args[0]
			err := kernels.RemoveKernel(context.Background(), log, dirName, kname, backupConf)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().BoolVar(&backupConf, "backup-conf", false, "backup configuration")
	cmd.Flags().StringVar(&dirName, dirNameCommand, "", dirNameHelp)
	cmd.MarkFlagRequired(dirNameCommand)
	return cmd
}
