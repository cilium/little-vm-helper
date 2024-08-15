// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/cilium/little-vm-helper/pkg/runner"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rcnf      runner.RunConf
	dirName   string
	pullImage bool

	ports []string
)

func RunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "run",
		Short:        "run/start VMs based on generated base images and kernels",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			rcnf.Logger = logrus.New()

			rcnf.ForwardedPorts, err = runner.ParsePortForward(ports)
			if err != nil {
				return fmt.Errorf("Port flags: %w", err)
			}

			t0 := time.Now()

			if _, err := os.Stat(rcnf.Image); errors.Is(err, os.ErrNotExist) && pullImage {
				pcnf := images.PullConf{
					Image:     rcnf.Image,
					TargetDir: dirName,
				}
				// Not a local file reference, could this be an OCI image?
				if err := images.PullImage(context.Background(), pcnf); err != nil {
					fmt.Fprintf(os.Stderr, "unable to find local file %q\n", rcnf.Image)
					return fmt.Errorf("unable to pull image: %w", err)
				}

				result, err := images.ExtractImage(context.Background(), pcnf)
				switch {
				case err == nil:
				case errors.Is(err, os.ErrExist):
				default:
					fmt.Fprintf(os.Stderr, "unable to find local file %q\n", rcnf.Image)
					return fmt.Errorf("unable to extract image: %w", err)
				}
				fmt.Printf("Autodetected image %s inside %s\n", result.Images[0], pcnf.Image)
				rcnf.Image = result.Images[0]
			}

			err = runner.StartQemu(rcnf)
			dur := time.Since(t0).Round(time.Millisecond)
			fmt.Printf("Execution took %v\n", dur)
			if err != nil {
				return fmt.Errorf("Qemu exited with an error: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&rcnf.Image, "image", "", "VM image file path")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringVar(&dirName, "dir", "_data", "directory to keep the images (configuration will be saved in <dir>/images.json and images in <dir>/images)")
	cmd.Flags().BoolVar(&pullImage, "pull-image", true, "Pull image from an OCI repository if it is not found locally")
	cmd.Flags().StringVar(&rcnf.KernelFname, "kernel", "", "kernel filename to boot with. (if empty no -kernel option will be passed to qemu)")
	cmd.Flags().BoolVar(&rcnf.QemuPrint, "qemu-cmd-print", false, "Do not run the qemu command, just print it")
	cmd.Flags().BoolVar(&rcnf.DisableHardwareAccel, "qemu-disable-kvm", false, "Do not use KVM acceleration, even if /dev/kvm exists")
	cmd.Flags().MarkDeprecated("qemu-disable-kvm", "use --no-hw-accel")
	cmd.Flags().BoolVar(&rcnf.DisableHardwareAccel, "no-hw-accel", false, "Do not use hardware acceleration, KVM for Linux or HVF for macOS")
	cmd.Flags().BoolVar(&rcnf.Daemonize, "daemonize", false, "daemonize QEMU after initializing")
	cmd.Flags().StringVar(&rcnf.ConsoleLogFile, "console-log-file", "", "Save VM console output to given file")
	cmd.Flags().StringVar(&rcnf.HostMount, "host-mount", "", "Mount the specified host directory in the VM using a 'host_mount' tag")
	cmd.Flags().StringArrayVarP(&ports, "port", "p", nil, "Forward a port (hostport[:vmport[:tcp|udp]])")
	cmd.Flags().IntVar(&rcnf.SerialPort, "serial-port", 0, "Port for serial console")
	cmd.Flags().IntVar(&rcnf.CPU, "cpu", 2, "CPU count (-smp)")
	cmd.Flags().StringVar(&rcnf.Mem, "mem", "4G", "RAM size (-m)")
	cmd.Flags().StringVar(&rcnf.CPUKind, "cpu-kind", "", "CPU kind to use (-cpu), has no effect when KVM is disabled (default 'kvm64' on amd64 and 'max' on arm64)")
	cmd.Flags().IntVar(&rcnf.QemuMonitorPort, "qemu-monitor-port", 0, "Port for QEMU monitor")
	cmd.Flags().StringVar(&rcnf.RootDev, "root-dev", "vda", "type of root device (hda or vda)")
	cmd.Flags().BoolVarP(&rcnf.Verbose, "verbose", "v", false, "Print qemu command before running it")
	cmd.Flags().StringVar(&rcnf.QemuArch, "qemu-arch", runtime.GOARCH, "specify qemu arch to use")

	return cmd
}
