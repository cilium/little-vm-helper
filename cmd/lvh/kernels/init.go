// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package kernels

import (
	"fmt"
	"strings"

	"github.com/cilium/little-vm-helper/pkg/kernels"
	"github.com/cilium/little-vm-helper/pkg/slogger"
	"github.com/spf13/cobra"
)

func initCommand() *cobra.Command {
	var force, backupConf bool
	var configGroups []string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize a directory for the kernel builder",
		Run: func(cmd *cobra.Command, _ []string) {
			log := slogger.New()
			conf := &kernels.Conf{
				Kernels:    make([]kernels.KernelConf, 0),
				CommonOpts: make([]kernels.ConfigOption, 0),
			}
			conf.AddGroupsCommonOpts(configGroups...)
			err := kernels.InitDir(log, dirName, conf, kernels.InitDirFlags{
				Force:      force,
				BackupConf: backupConf,
			})
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "force init")
	cmd.Flags().BoolVar(&backupConf, "backup-conf", false, "backup configuration")
	cmd.Flags().StringSliceVar(
		&configGroups,
		"config-groups", kernels.DefaultConfigGroups,
		fmt.Sprintf(
			"add configuration options based on the following predefined groups: %s",
			strings.Join(kernels.GetConfigGroupNames(), ","),
		),
	)
	cmd.Flags().StringVar(&dirName, dirNameCommand, "", dirNameHelp)
	cmd.MarkFlagRequired(dirNameCommand)
	return cmd
}
