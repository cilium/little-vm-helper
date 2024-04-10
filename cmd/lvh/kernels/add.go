// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addCommand() *cobra.Command {
	var addConfigGroups []string
	var addPrintConfig, addFetch, backupConf bool
	addCmd := &cobra.Command{
		Use:     "add <name> <url>",
		Short:   "add kernel (by updating config file in directory)",
		Example: kernels.GetExamplesText(),
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()

			kconf := kernels.KernelConf{
				Name: args[0],
				URL:  args[1],
				Opts: make([]kernels.ConfigOption, 0),
			}

			if err := kconf.AddGroupsOpts(addConfigGroups...); err != nil {
				log.Fatal(err)
			}

			if err := kconf.Validate(); err != nil {
				log.Fatal(err)
			}

			if addPrintConfig {
				confb, err := json.MarshalIndent(kconf, "", "    ")
				if err != nil {
					log.Fatal(fmt.Errorf("failed to marshal config: %w", err))
				}
				os.Stdout.Write(confb)
				return
			}

			err := kernels.AddKernel(context.Background(), log, dirName, &kconf, kernels.AddKernelFlags{
				BackupConf: backupConf,
				Fetch:      addFetch,
			})
			if err != nil {
				log.Fatal(err)
			}

		},
	}
	addCmd.Flags().StringSliceVar(
		&addConfigGroups,
		"config-groups", []string{},
		fmt.Sprintf(
			"add configuration options based on the following predefined groups: %s",
			strings.Join(kernels.GetConfigGroupNames(), ","),
		),
	)
	addCmd.Flags().BoolVar(&addPrintConfig, "just-print-config", false, "do not actually add the kernel. Just print its config.")
	addCmd.Flags().BoolVar(&addFetch, "fetch", false, "fetch URL")
	addCmd.Flags().BoolVar(&backupConf, "backup-conf", false, "backup configuration")
	addCmd.Flags().StringVar(&dirName, dirNameCommand, "", dirNameHelp)
	addCmd.MarkFlagRequired(dirNameCommand)
	return addCmd
}
