// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package images

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cilium/little-vm-helper/pkg/images"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func BuildCmd() *cobra.Command {
	var dirName, imageName, arch string
	var forceRebuild, dryRun, mergeSteps bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build VM images",
		RunE: func(cmd *cobra.Command, _ []string) error {
			log := logrus.New()
			imagesDir := filepath.Join(dirName, "images")
			configFname := filepath.Join(dirName, images.DefaultConfFile)

			configData, err := os.ReadFile(configFname)
			if err != nil {
				return err
			}

			var cnf images.ImagesConf
			cnf.Dir, err = filepath.Abs(imagesDir)
			if err != nil {
				return err
			}
			err = json.Unmarshal(configData, &cnf.Images)
			if err != nil {
				return err
			}

			forest, err := images.NewImageForest(&cnf, false)
			if err != nil {
				return err
			}

			var res *images.BuilderResult
			bldConf := &images.BuildConf{
				Log:          log,
				DryRun:       dryRun,
				ForceRebuild: forceRebuild,
				MergeSteps:   mergeSteps,
				Arch:         arch,
			}
			start := time.Now()
			if imageName == "" {
				res = forest.BuildAllImages(bldConf)
			} else {
				res, err = forest.BuildImage(bldConf, imageName)
				if err != nil {
					log.WithField("image", imageName).WithError(err).Error("error bulding image")
					return err
				}
			}
			elapsed := time.Since(start)

			err = res.Err()
			if err != nil {
				log.WithError(err).Error("building images failed")
			} else {
				log.WithField("time-elapsed", elapsed).Info("images built succesfully")
			}

			for img, ir := range res.ImageResults {
				if ir.Error == nil {
					fmt.Printf("image:%-10s cachedImageUsed:%t cachedImageDeleted:%s\n", img, ir.CachedImageUsed, ir.CachedImageDeleted)
				}
			}

			return err
		},
	}

	cmd.Flags().StringVar(&dirName, "dir", "", "directory to keep the images (configuration will be saved in <dir>/images.json and images in <dir>/images)")
	cmd.Flags().StringVar(&imageName, "image", "", "image to build. If empty, all images will be build.")
	cmd.MarkFlagRequired("dir")
	cmd.Flags().BoolVar(&forceRebuild, "force-rebuild", false, "rebuild all images, even if they exist")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "do the whole thing, but instead of building actual images create empty files")
	cmd.Flags().BoolVar(&mergeSteps, "merge-steps", true, "Merge steps when possible to improve performance. Disabling this might be useufl to investigate action issues.")
	cmd.Flags().StringVar(&arch, "arch", "", "target architecture to build the image, e.g. 'amd64' or 'arm64' (default to native architecture)")
	return cmd
}
