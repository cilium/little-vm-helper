package images

// RunScript runs a script in a path specified by a string
type RunScript string

// CopyFile copies a file from the host inside an image
type CopyFile struct {
	HostPath  string
	ImagePath string
}

// NB:Other potential actions
//  - CopyFileFromDockerImage
//  - ...

type Action interface{}
