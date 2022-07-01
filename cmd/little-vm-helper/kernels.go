package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cilium/little-vm-helper/pkg/kernels"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func KernelsCommand() *cobra.Command {
	var dirName string
	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "build kernels",
	}
	cmd.PersistentFlags().StringVar(&dirName, "dir", "", "directory to to place kernels")
	cmd.MarkFlagRequired("dir")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "initialize a directory for the kernel builder",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			err := kernels.InitDir(dirName, nil)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list available kernels (by reading config file in directory)",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			kd, err := kernels.LoadDir(dirName)
			if err != nil {
				log.Fatal(err)
			}

			for _, k := range kd.Conf.Kernels {
				fmt.Printf("%-13s %s", k.Name, k.URL)
			}
		},
	}

	var addConfigGroups []string
	var addPrintConfig bool
	addCmd := &cobra.Command{
		Use:     "add <name> <url> (e.g., add bpf-next )",
		Short:   "add kernel (by updating config file in directory)",
		Example: fmt.Sprintf("add bpf-next git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git\nadd 5.18.8 git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#v5.18.8"),
		Run: func(cmd *cobra.Command, args []string) {
			log := logrus.New()

			kconf := kernels.KernelConf{
				Name: args[0],
				URL:  args[1],
			}

			if err := kconf.AddGroups(addConfigGroups...); err != nil {
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

			err := kernels.AddKernel(dirName, &kconf)
			if err != nil {
				log.Fatal(err)
			}

		},
		Args: cobra.ExactArgs(2),
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
	cmd.AddCommand(initCmd, listCmd, addCmd)

	/*


		add := &cobra.Command{
			Use:   "build",
			Short: "build kernel",
		}

		rm := &cobra.Command{}
	*/

	return cmd
}
