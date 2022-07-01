package images

import "github.com/sirupsen/logrus"

// StepConf is common step configuration
type StepConf struct {
	imageDir string
	imgCnf   *ImgConf
	log      *logrus.Logger
}
