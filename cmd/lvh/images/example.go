// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package images

import (
	"encoding/json"
	"log"
	"os"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/spf13/cobra"
)

func ExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-config",
		Short: "Print an example config",
		Run: func(cmd *cobra.Command, _ []string) {
			conf := &images.ExampleImagesConf
			confb, err := json.MarshalIndent(conf, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(confb)
		},
	}
	return cmd
}
