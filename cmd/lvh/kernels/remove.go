package kernels

import (
	"context"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func removeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <kernel>",
		Short: "remove kernel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kname := args[0]
			err := kernels.RemoveKernel(context.Background(), log, dirName, kname)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
