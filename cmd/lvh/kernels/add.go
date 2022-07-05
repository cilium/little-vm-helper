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
	addExamples := func() string {
		var sb strings.Builder

		for _, ex := range kernels.UrlExamples {
			sb.WriteString(fmt.Sprintf("  add %s %s\n", ex.Name, ex.URL))
		}

		return sb.String()
	}

	var addConfigGroups []string
	var addPrintConfig, addFetch bool
	addCmd := &cobra.Command{
		Use:     "add <name> <url> (e.g., add bpf-next )",
		Short:   "add kernel (by updating config file in directory)",
		Example: addExamples(),
		Args:    cobra.ExactArgs(2),
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

			kURL, err := kernels.ParseURL(kconf.URL)
			if err != nil {
				log.Fatal(err)
			}

			if addFetch {
				err := kURL.Fetch(context.Background(), log, dirName, kconf.Name)
				if err != nil {
					log.Fatal(err)
				}
			}

		},
	}
	addCmd.Flags().StringSliceVar(
		&addConfigGroups,
		"config-groups", kernels.DefaultConfigGroups,
		fmt.Sprintf(
			"add configuration options based on the following predefined groups: %s",
			strings.Join(kernels.GetConfigGroupNames(), ","),
		),
	)
	addCmd.Flags().BoolVar(&addPrintConfig, "just-print-config", false, "do not actually add the kernel. Just print its config.")
	addCmd.Flags().BoolVar(&addFetch, "fetch", false, "fetch URL")
	return addCmd
}
