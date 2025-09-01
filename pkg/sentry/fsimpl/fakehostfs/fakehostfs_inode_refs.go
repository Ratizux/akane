package fakehostfs

import (
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
)

func (i *FakehostfsInode) IncRef() {
	log.Debugf("IncRef() called on inode %d",i.Ino())
	return
}

func (i *FakehostfsInode) DecRef(ctx context.Context) {
	log.Debugf("DecRef() called on inode %d",i.Ino())
	return
}

func (i *FakehostfsInode) TryIncRef() bool {
	log.Debugf("TryIncRef() called on inode %d",i.Ino())
	return true
}
