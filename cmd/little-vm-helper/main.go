package main

import (
	"github.com/spf13/cobra"
)

var (
	exampleConfigCmd = ExampleConfigCommand()
	buildImagesCmd   = BuildImagesCommand()
	rootCmd          = &cobra.Command{
		Use:   "little-vm-helper",
		Short: "hellper to build and run VMs",
	}
)

func init() {
	rootCmd.AddCommand(buildImagesCmd)
	rootCmd.AddCommand(exampleConfigCmd)
}

func main() {
	rootCmd.Execute()
}
