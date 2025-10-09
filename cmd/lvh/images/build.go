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
	var dirName, pkgRepository string
	var forceRebuild, dryRun, mergeSteps bool
	var imageNames []string

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
				PkgRepo:      pkgRepository,
			}
			start := time.Now()
			if imageNames == nil {
				res = forest.BuildAllImages(bldConf)
			} else {
				res = forest.BuildImages(bldConf, imageNames)
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
	cmd.MarkFlagRequired("dir")
	cmd.Flags().StringArrayVarP(&imageNames, "image", "i", nil, "images to build. If empty, all images will be built.")
	cmd.Flags().BoolVar(&forceRebuild, "force-rebuild", false, "rebuild all images, even if they exist")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "do the whole thing, but instead of building actual images create empty files")
	cmd.Flags().BoolVar(&mergeSteps, "merge-steps", true, "Merge steps when possible to improve performance. Disabling this might be useufl to investigate action issues.")
	cmd.Flags().StringVar(&pkgRepository, "pkg-repo", "sid", "repository used to get packages from when building a base image (ex: sid, unstable, bookworm, stable, oldstable, etc..)")
	return cmd
}
