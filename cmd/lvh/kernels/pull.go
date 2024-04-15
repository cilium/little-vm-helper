package kernels

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	cranev1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/spf13/cobra"
)

func pullCommand() *cobra.Command {
	var (
		platform string
		repo     string
	)

	cmd := &cobra.Command{
		Use:   "pull <tag>",
		Short: "pull built kernels by cilium/little-vm-helper-images",
		Long: `Pull built kernels artifacts by cilium/little-vm-helper-images

Examples:
  # Pull latest main tag for version 6.6
  lvh kernels pull 6.6-main

  # Pull the bpf-next-20240404.012646 tag for version bpf-next
  lvh kernels pull bpf-next-20240404.012646

  # Pull the main tag for version 5.4 for arm64
  lvh kernels pull 5.4-main --platform linux/arm64

  # Pull the main tag for version 5.10 in directory mykernels
  lvh kernels pull 5.10-main --dir mykernels

  # Pull the latest tags available for version 5.15 without using 5.15-main
  lvh kernels pull $(lvh kernels catalog 5.15 | tail -n 2 | head -n 1)

  # Pull the latest CI-generated images for version bpf-next
  lvh kernels pull bpf-next-main --repo quay.io/lvh-images/kernel-images-ci`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			re := regexp.MustCompile(kernelTagRegex)
			match := re.FindStringSubmatch(args[0])
			if len(match) < 2 {
				return fmt.Errorf("tag is malformed, it should follow %s: %s", re, args[0])
			}

			sPlatform := strings.SplitN(platform, "/", 2)
			if len(sPlatform) < 2 {
				return fmt.Errorf("platform is malformed, it must be <os>/<arch>: %s", platform)
			}

			srcImage := fmt.Sprintf("%s:%s", repo, args[0])
			dstTarFile := fmt.Sprintf("%s.tar", args[0])

			err := FetchTarImage(srcImage, dstTarFile, sPlatform[0], sPlatform[1])
			defer os.Remove(dstTarFile) // this is on purpose before if != nil check
			if err != nil {
				return fmt.Errorf("failed fetching and creating tar image: %w", err)
			}

			pathToExtract := fmt.Sprintf("data/kernels/%s", match[1])

			err = ExtractTarPath(dstTarFile, pathToExtract, filepath.Join(dirName, args[0]))
			if err != nil {
				return fmt.Errorf("failed extracting the file from the tar image: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&platform, "platform", "p", runtime.GOOS+"/"+runtime.GOARCH, "platform for the kernel image <os>/<arch>")
	cmd.Flags().StringVarP(&dirName, dirNameCommand, "d", ".", dirNameHelp)
	cmd.Flags().StringVar(&repo, repoCommand, kernelImageRepository, repoHelp)

	return cmd
}

func FetchTarImage(src, dst, platformOS, platformArch string) error {
	var img cranev1.Image

	desc, err := crane.Get(src, crane.WithPlatform(&cranev1.Platform{
		OS:           platformOS,
		Architecture: platformArch,
	}))
	if err != nil {
		return fmt.Errorf("pulling %s: %w", src, err)
	}
	if desc.MediaType.IsSchema1() {
		img, err = desc.Schema1()
		if err != nil {
			return fmt.Errorf("pulling schema 1 image %s: %w", src, err)
		}
	} else {
		img, err = desc.Image()
		if err != nil {
			return fmt.Errorf("pulling Image %s: %w", src, err)
		}
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	return crane.Export(img, file)
}

func ExtractTarPath(tarFile string, path string, targetDir string) error {
	file, err := os.Open(tarFile)
	if err != nil {
		return fmt.Errorf("failed to open the tar file: %w", err)
	}
	tarReader := tar.NewReader(file)
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break // end of archive
		}

		if err != nil {
			return err
		}

		if filepath.IsAbs(header.Name) {
			return fmt.Errorf("absolute path file in archive: %s", header.Name)
		}

		name, found := strings.CutPrefix(header.Name, path)
		if !found {
			// skip all non interesting files
			continue
		}

		dstPath := filepath.Join(targetDir, name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dstPath, err)
			}
		case tar.TypeReg:
			dstFile, err := os.Create(dstPath)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", dstPath, err)
			}
			defer dstFile.Close()
			n, err := io.CopyN(dstFile, tarReader, header.Size)
			if err != nil {
				return fmt.Errorf("failed to copy %s from tar %s: %w", dstPath, tarFile, err)
			}
			if n != header.Size {
				return fmt.Errorf("tar header reports file %s size %d, but only %d bytes were pulled", header.Name, header.Size, n)
			}

		default:
			// skip all other tar header types (symlinks, etc.)
			continue
		}
	}
	return nil
}
