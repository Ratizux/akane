package fakehostfs

import (
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

func (i *FakehostfsInode) Locks() *vfs.FileLocks {
	return &i.locks
}
