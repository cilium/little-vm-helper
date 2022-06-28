package images

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// ImageConf describes the configuration of an image
type ImageConf struct {
	Name     string
	Parent   string
	Packages []string
	Actions  []Action
}

// ImageBuilder can be used to build images
type ImageBuilder struct {
	imageDir string
	confs    map[string]*ImageConf
	children map[string][]string
}

type ImageBuilderConf struct {
	ImageDir string
	Images   []ImageConf
}

// NewImageBuilder creates a new image builder
func NewImageBuilder(conf *ImageBuilderConf) (*ImageBuilder, error) {
	// image name -> ImageConf
	imgConfs := make(map[string]*ImageConf, len(conf.Images))
	// name -> parent name (if parent exists)
	imageParent := make(map[string]string)

	for i := range conf.Images {
		icnf := &conf.Images[i]
		if _, ok := imgConfs[icnf.Name]; ok {
			return nil, fmt.Errorf("duplicate image name: %s", icnf.Name)
		}
		imgConfs[icnf.Name] = icnf
		if icnf.Parent != "" {
			imageParent[icnf.Name] = icnf.Parent
		}
	}

	children := make(map[string][]string)
	for child, parent := range imageParent {
		if _, ok := imgConfs[parent]; !ok {
			return nil, fmt.Errorf("image '%s' specified as parent, but it is not defined", parent)
		}
		children[parent] = append(children[parent], child)
	}

	err := os.MkdirAll(conf.ImageDir, 0755)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	confb, err := json.Marshal(conf)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(path.Join(conf.ImageDir, ConfFile), confb, 0666)
	if err != nil {
		return nil, fmt.Errorf("error writing configuration: %w", err)
	}

	return &ImageBuilder{
		imageDir: conf.ImageDir,
		confs:    imgConfs,
		children: children,
	}, nil
}

func (ib *ImageBuilder) ImageFilename(image string) (string, error) {
	if _, ok := ib.confs[image]; !ok {
		return "", fmt.Errorf("no configuration for image '%s'", image)
	}

	return ib.imageFilename(image), nil
}

func (ib *ImageBuilder) imageFilename(image string) string {
	return path.Join(ib.imageDir, fmt.Sprintf("%s.%s", image, ImageExt))
}

// getDependencies returns the dependencies of an image
func (ib *ImageBuilder) getDependencies(image string) ([]string, error) {
	var ret []string
	cnf, ok := ib.confs[image]
	if !ok {
		return ret, fmt.Errorf("cannot build dependencies for image %s, because image does not exist ", image)
	}

	parent := cnf.Parent
	for parent != "" {
		// NB: we have checked already that all parents exist
		cnfParent := ib.confs[parent]
		ret = append(ret, parent)
		parent = cnfParent.Parent
	}

	// reverse return slice
	for i, j := 0, len(ret)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret, nil
}

func (ib *ImageBuilder) IsLeafImage(i string) bool {
	_, hasChidren := ib.children[i]
	return !hasChidren
}

func (ib *ImageBuilder) LeafImages() []string {
	ret := make([]string, 0)
	for i, _ := range ib.confs {
		if ib.IsLeafImage(i) {
			ret = append(ret, i)
		}
	}
	return ret
}

func (ib *ImageBuilder) RootImages() []string {
	ret := make([]string, 0)
	for i, cnf := range ib.confs {
		if cnf.Parent == "" {
			ret = append(ret, i)
		}
	}

	return ret
}
