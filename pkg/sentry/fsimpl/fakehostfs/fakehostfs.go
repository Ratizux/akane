// Package fakehostfs implements fakehostfs.
package fakehostfs

import (
	//"strconv"

	//"gvisor.dev/gvisor/pkg/errors/linuxerr"
	//"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	//"gvisor.dev/gvisor/pkg/sentry/vfs"
	//"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	//"gvisor.dev/gvisor/pkg/context"
)

const (
	// Name is the default filesystem name.
	Name                     = "fakehostfs"
	defaultMaxCachedDentries = uint64(1000)
)

// vfs.FilesystemImpl
// vfs.FilesystemType
