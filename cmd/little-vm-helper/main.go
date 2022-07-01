package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "little-vm-helper",
		Short: "hellper to build and run VMs",
	}
)

func init() {
	rootCmd.AddCommand(
		BuildImagesCommand(),
		ExampleConfigCommand(),
		KernelsCommand(),
	)
}

func main() {
	rootCmd.Execute()
}
