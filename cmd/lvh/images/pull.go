// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package images

import (
	"context"
	"fmt"
	"os"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/spf13/cobra"
)

var (
	dirName string
	cache   bool
)

func PullCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:        "pull <URL>",
		Short:      "Pull an image from an OCI repository",
		Args:       cobra.MinimumNArgs(1),
		ArgAliases: []string{"imageURL"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := os.MkdirAll(dirName, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dirName, err)
			}

			_, err := images.PullImage(context.Background(), images.PullConf{
				Image:     args[0],
				TargetDir: dirName,
				Cache:     cache,
			})
			return err
		},
	}

	cmd.Flags().StringVar(&dirName, "dir", "_data", "directory to keep the images (images will be saved in images in <dir>/images)")
	cmd.Flags().BoolVar(&cache, "cache", false, "cache a compressed version of the image")

	return cmd
}
