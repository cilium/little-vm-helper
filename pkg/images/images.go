package images

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var (
	ConfFile = "conf.json"
)

// code for creating images

// ImageConf describes the configuration of an image
type ImageConf struct {
	Name     string
	Parent   string
	Packages []string
	Actions  []Action
}

// ImageBuilder can be used to build images
type ImageBuilder struct {
	confs    map[string]*ImageConf
	imageDir string
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

	for _, parent := range imageParent {
		if _, ok := imgConfs[parent]; !ok {
			return nil, fmt.Errorf("image '%s' specified as parent, but it is not defined", parent)
		}
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
		confs:    imgConfs,
		imageDir: conf.ImageDir,
	}, nil
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
