package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/spf13/cobra"
)

func ExampleConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-config",
		Short: "Print an example config file",
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
