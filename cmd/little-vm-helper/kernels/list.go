package kernels

import (
	"fmt"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func listCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list available kernels (by reading config file in directory)",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			for _, k := range kd.Conf.Kernels {
				fmt.Printf("%-13s %s %s\n", k.Name, k.URL)
			}
		},
	}
}
