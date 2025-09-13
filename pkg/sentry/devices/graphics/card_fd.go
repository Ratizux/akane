package graphics

import (
	"io"

	//"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/usermem"
	"gvisor.dev/gvisor/pkg/sentry/arch"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/sentry/memmap"
	// "gvisor.dev/gvisor/pkg/sync"
)

func NewCardFD(ctx context.Context, mnt *vfs.Mount, vfsd *vfs.Dentry, opts vfs.OpenOptions, deviceState *CardDeviceState) (*vfs.FileDescription, error) {
	fd := &CardFD{}
	if deviceState == nil {
		return nil, linuxerr.EINVAL
	}
	fd.state = deviceState
	// TODO adjust fd mmapfile size

	if err := fd.vfsfd.Init(fd, opts.Flags, mnt, vfsd, &vfs.FileDescriptionOptions{
		UseDentryMetadata: true,
	}); err != nil {
		return nil, err
	}
	return &fd.vfsfd, nil
}

type CardFD struct {
	vfsfd vfs.FileDescription
	vfs.FileDescriptionDefaultImpl
	vfs.DentryMetadataFileDescriptionImpl
	vfs.NoLockFD

	mappings memmap.MappingSet
	state *CardDeviceState
}

// Release implements vfs.FileDescriptionImpl.Release.
func (fd *CardFD) Release(context.Context) {
	// noop
}

// PRead implements vfs.FileDescriptionImpl.PRead.
func (fd *CardFD) PRead(ctx context.Context, dst usermem.IOSequence, offset int64, opts vfs.ReadOptions) (int64, error) {
	return 0, io.EOF
}

// Read implements vfs.FileDescriptionImpl.Read.
func (fd *CardFD) Read(ctx context.Context, dst usermem.IOSequence, opts vfs.ReadOptions) (int64, error) {
	return 0, io.EOF
}

// PWrite implements vfs.FileDescriptionImpl.PWrite.
func (fd *CardFD) PWrite(ctx context.Context, src usermem.IOSequence, offset int64, opts vfs.WriteOptions) (int64, error) {
	return src.NumBytes(), nil
}

// Write implements vfs.FileDescriptionImpl.Write.
func (fd *CardFD) Write(ctx context.Context, src usermem.IOSequence, opts vfs.WriteOptions) (int64, error) {
	return src.NumBytes(), nil
}

// Seek implements vfs.FileDescriptionImpl.Seek.
func (fd *CardFD) Seek(ctx context.Context, offset int64, whence int32) (int64, error) {
	return 0, nil
}

func (fd *CardFD) Ioctl(ctx context.Context, uio usermem.IO, sysno uintptr, args arch.SyscallArguments) (uintptr, error) {
	request := args[1].Uint()
	data := args[2].Pointer()

	t := kernel.TaskFromContext(ctx)
	if t == nil {
		panic("Ioctl should be called from a task context")
	}

	log.Debugf("request is %d",request)
	switch request {
	case DRM_IOCTL_VERSION:
		return DrmIoctlVersionHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_GETRESOURCES:
		return DrmIoctlModeGetResourcesHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_GETCONNECTOR:
		return DrmIoctlModeGetConnectorHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_CREATE_DUMB:
		return DrmModeCreateDumbHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_ADDFB:
		return DrmIoctlModeAddFbHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_ADDFB2:
		return DrmIoctlModeAddFb2Handler(fd.state, t, data)

	case DRM_IOCTL_MODE_GETENCODER:
		return DrmIoctlModeGetEncoderHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_MAP_DUMB:
		return DrmIoctlModeMapDumbHandler(fd.state, t, data)

	case DRM_IOCTL_SET_MASTER:
		return DrmIoctlSetMasterHandler(fd.state, t, data)

	case DRM_IOCTL_MODE_GETCRTC:
		return DrmIoctlModeGetCRTC_Handler(fd.state, t, data)

	case DRM_IOCTL_MODE_SETCRTC:
		return DrmIoctlModeSetCRTC_Handler(fd.state, t, data)

	case DRM_IOCTL_MODE_PAGE_FLIP:
		return DrmIoctlModeCRTC_PageFlipHandler(fd.state, t, data)

	case DRM_IOCTL_GET_CAP:
		return DrmIoctlGetCapHandler(fd.state, t, data)

	default:
		if fd.state.mf != nil && fd.state.mf.buffer != nil {
			for _, n := range(fd.state.mf.buffer[0:10]) {
				log.Debugf("%08b", n)
			}
		}
		return 0, linuxerr.EINVAL
	}
}
