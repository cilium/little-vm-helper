package kernels

type UrlExample struct {
	Name string
	URL  string

	// NB: used for testing
	expectedKernelURL KernelURL
}

var UrlExamples = []UrlExample{
	{
		Name: "bpf-next",
		URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git",
		expectedKernelURL: &GitURL{
			Repo:   "git://git.kernel.org/pub/scm/linux/kernel/git/bpf/bpf-next.git",
			Branch: "master",
		},
	}, {
		Name: "5.18",
		URL:  "git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git#linux-5.18.y",
		expectedKernelURL: &GitURL{
			Repo:   "git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git",
			Branch: "linux-5.18.y",
		},
	},
}
