package kernels

type KernelsDir struct {
	Dir  string
	Conf Conf
}

func (kd *KernelsDir) KernelConfig(name string) *KernelConf {
	for i := range kd.Conf.Kernels {
		if kd.Conf.Kernels[i].Name == name {
			return &kd.Conf.Kernels[i]
		}
	}

	return nil
}
