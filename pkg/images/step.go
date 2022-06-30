package images

import "github.com/sirupsen/logrus"

// StepConf is common step configuration
type StepConf struct {
	imageDir string
	imgCnf   *ImageConf
	log      *logrus.Logger
}
