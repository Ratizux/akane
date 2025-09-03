package fakehostfs

import (
	//"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	//"gvisor.dev/gvisor/pkg/atomicbitops"
	//"gvisor.dev/gvisor/pkg/context"
)

// FakehostfsDentry implements vfs.DentryImpl.
type FakehostfsDentry struct {
	kernfs.Dentry
}
