// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"fmt"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/spf13/cobra"
)

func listCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list available kernels (by reading config file in directory)",
		Run: func(cmd *cobra.Command, _ []string) {
			log := slogger.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			for _, k := range kd.Conf.Kernels {
				fmt.Printf("%-13s %s\n", k.Name, k.URL)
			}
		},
	}

	cmd.Flags().StringVar(&dirName, dirNameCommand, "", dirNameHelp)
	cmd.MarkFlagRequired(dirNameCommand)

	return cmd
}
