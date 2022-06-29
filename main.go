package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/kkourt/vamp/pkg/images"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	exampleConfigCmd = ExampleConfigCommand()
	buildImagesCmd   = BuildImagesCommand()
	rootCmd          = &cobra.Command{
		Use:   "vamp",
		Short: "vamp builds and runs VMs",
	}
)

func init() {
	rootCmd.AddCommand(buildImagesCmd)
	rootCmd.AddCommand(exampleConfigCmd)
}

func ExampleConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example-config",
		Short: "Print an example config file",
		Run: func(cmd *cobra.Command, _ []string) {
			images := []images.ImageConf{
				{
					Name: "base",
					Packages: []string{
						"less",
						"vim",
						"sudo",
						"openssh-server",
						"curl",
					},
				},
				{
					Name: "k8s",
					Packages: []string{
						"docker.io",
					},
				},
			}

			confb, err := json.MarshalIndent(images, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(confb)
		},
	}
	return cmd
}

func BuildImagesCommand() *cobra.Command {
	var configFname, dirName string

	cmd := &cobra.Command{
		Use:   "build-images",
		Short: "Build VM images",
		Run: func(cmd *cobra.Command, _ []string) {
			log := logrus.New()
			if configFname == "" {
				configFname = path.Join(dirName, images.DefaultConfFile)
			}

			configData, err := os.ReadFile(configFname)
			if err != nil {
				log.Fatal(err)
			}

			var cnf images.BuilderConf
			cnf.ImageDir = dirName
			err = json.Unmarshal(configData, &cnf.Images)
			if err != nil {
				log.Fatal(err)
			}

			builder, err := images.NewImageBuilder(&cnf)
			if err != nil {
				log.Fatal(err)
			}

			start := time.Now()
			res := builder.BuildAllImages(&images.BuildConf{
				Log:    log,
				DryRun: false,
			})
			elapsed := time.Since(start)

			if err := res.Err(); err != nil {
				log.WithError(err).Error("building images failed")
			} else {
				log.WithField("time-elapsed", elapsed).Info("images built succesfully")
			}

			for img, ir := range res.ImageResults {
				if ir.Error == nil {
					fmt.Printf("image:%-10s cachedImageUsed:%t cachedImageDeleted:%s\n", img, ir.CachedImageUsed, ir.CachedImageDeleted)
				}
			}
		},
	}
	cmd.Flags().StringVar(&dirName, "dir", "", "directory  to place images")
	cmd.MarkFlagRequired("dir")
	cmd.Flags().StringVar(&configFname,
		"config", "",
		fmt.Sprintf("config file (default is <dir>/%s)", images.DefaultConfFile),
	)
	return cmd
}

func main() {
	rootCmd.Execute()
}
