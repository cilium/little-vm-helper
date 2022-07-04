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

func (kd *KernelsDir) RemoveKernelConfig(name string) *KernelConf {
	for i := range kd.Conf.Kernels {
		if kd.Conf.Kernels[i].Name == name {
			ret := &kd.Conf.Kernels[i]
			kd.Conf.Kernels = append(kd.Conf.Kernels[:i], kd.Conf.Kernels[i+1:]...)
			return ret
		}
	}

	return nil
}
