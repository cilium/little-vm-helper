package kernels

import (
	"context"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func buildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "build <kernel>",
		Short: "build kernel",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			kname := args[0]
			kconf := kd.KernelConfig(kname)
			if kconf == nil {
				log.Fatalf("kernel `%s` not found", kname)
			}

			kURL, err := kernels.ParseURL(kconf.URL)
			if err != nil {
				log.Fatal(err)
			}

			err = kURL.Fetch(context.Background(), log, dirName, kconf.Name)
			if err != nil {
				log.Fatal(err)
			}

			err = kconf.Build(context.Background(), log, dirName)
			if err != nil {
				log.Fatal(err)
			}

		},
	}
}
