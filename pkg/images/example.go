package images

var ExampleImagesConf = []ImageConf{
	{
		Name: "base",
		Packages: []string{
			"less",
			"vim",
			"sudo",
			"openssh-server",
			"curl",
		},
		Actions: []Action{{
			Comment: "disable password for root",
			Op: &RunCommand{
				Cmd: "passwd -d root",
			},
		}},
	},
	{
		Name:   "k8s",
		Parent: "base",
		Packages: []string{
			"docker.io",
		},
	},
}
