package images

var (
	// Binaries used
	// Debootstrap = "debootstrap"
	Mmdebstrap    = "mmdebstrap"
	QemuImg       = "qemu-img"
	VirtCustomize = "virt-customize"
	GuestFish     = "guestfish"

	Binaries = []string{
		Mmdebstrap,
		QemuImg,
		VirtCustomize,
		GuestFish,
	}

	DefaultImagesDir = "images"
	DefaultConfFile  = "images.json"
	DefaultImageExt  = "img"
	DefaultImageSize = "8G"
)
