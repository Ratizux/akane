package fakehostfs

import (
	//"gvisor.dev/gvisor/pkg/context"
	//"gvisor.dev/gvisor/pkg/sentry/vfs"
	//"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	//"gvisor.dev/gvisor/pkg/errors/linuxerr"
	//"gvisor.dev/gvisor/pkg/abi/linux"
	//"gvisor.dev/gvisor/pkg/sentry/socket/unix/transport"
	//"gvisor.dev/gvisor/pkg/fspath"
)

//see sentry/vfs/filesytem.go: FilesystemImpl interface

// FakehostfsImpl implements vfs.FilesystemImpl.
type FakehostfsImpl struct{
	kernfs.Filesystem

	root *FakehostfsDentry
	rootNodeID uint64
	nativeFS nativeFilesystem

	devMajor uint32
	devMinor uint32
}

// MountOptions implements vfs.FilesystemImpl.MountOptions.
func (fs *FakehostfsImpl) MountOptions() string {
	return ""
}
