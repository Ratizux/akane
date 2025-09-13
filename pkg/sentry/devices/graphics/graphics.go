package graphics

import (
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

const (
	CardDevMinor = 0
)

func Register(vfsObj *vfs.VirtualFilesystem) error {
	card0 := &CardDevice{}
	card0.Init()

	for minor, spec := range map[uint32]struct {
		dev      vfs.Device
		pathname string
	}{
		CardDevMinor:    {card0, "dri/card0"},
		// TODO renderD
	} {
		if err := vfsObj.RegisterDevice(vfs.CharDevice, linux.DRM_MAJOR, minor, spec.dev, &vfs.RegisterDeviceOptions{
			GroupName: "drm",
			Pathname:  spec.pathname,
			FilePerms: 0666,
		}); err != nil {
			return err
		}
	}
	return nil
}
