package kernels

import (
	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "initialize a directory for the kernel builder",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			err := kernels.InitDir(log, dirName, nil)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}
