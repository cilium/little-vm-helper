package main

import (
	"github.com/cilium/little-vm-helper/cmd/little-vm-helper/images"
	"github.com/cilium/little-vm-helper/cmd/little-vm-helper/kernels"

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
		images.ImagesCommand(),
		kernels.KernelsCommand(),
	)
}

func main() {
	rootCmd.Execute()
}
