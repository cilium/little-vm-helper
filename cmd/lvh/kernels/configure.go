// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"context"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func configureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure <kernel>",
		Short: "configure kernel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			kname := args[0]
			if err := kd.ConfigureKernel(context.Background(), log, kname); err != nil {
				log.Fatal(err)
			}

		},
	}
}

func rawConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "raw_configure <kernel_dir> [<kernel_name>]",
		Short: "configure a kernel prepared by means other than lvh",
		Args:  cobra.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			kdir := args[0]
			kname := ""
			if len(args) > 1 {
				kname = args[1]
			}
			if err := kd.RawConfigure(context.Background(), log, kdir, kname); err != nil {
				log.Fatal(err)
			}

		},
	}
}
