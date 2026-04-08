package fakehostfs

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/atomicbitops"
)
//"path"
//"strings"

//"golang.org/x/sys/unix"
//"gvisor.dev/gvisor/pkg/abi/linux"
//"gvisor.dev/gvisor/pkg/errors/linuxerr"
//"gvisor.dev/gvisor/pkg/log"
//"gvisor.dev/gvisor/pkg/fsutil"
//"gvisor.dev/gvisor/pkg/sentry/vfs"

type EntryType int

const (
	ENTRY_REGULAR EntryType = iota
	ENTRY_DIRECTORY
	ENTRY_SYMLINK
)

type logicalInode struct {
	hostfd       int
	hostfdOpen   bool
	mappingCount int32

	inodeType EntryType

	//inode attrs
	devMajor  uint32
	devMinor  uint32
	ino       atomicbitops.Uint64
	mode      atomicbitops.Uint32
	uid       atomicbitops.Uint32
	gid       atomicbitops.Uint32
	nlink     atomicbitops.Uint32
	blockSize atomicbitops.Uint32
	// Timestamps, all nsecs from the Unix epoch.
	atime atomicbitops.Int64
	mtime atomicbitops.Int64
	ctime atomicbitops.Int64
}

func (li *logicalInode) Init(now int64, devMajor uint32, devMinor uint32, ino uint64, inodeMetadata InodeMetadata) error {
	mode := linux.FileMode(inodeMetadata.Mode)
	if mode.FileType() == 0 {
		panic(fmt.Sprintf("No file type specified in 'mode' for FakehostfsInode.Init(): mode=0%o", mode))
	}
	fileType := inodeMetadata.Mode & STAT_TYPE_MASK
	switch fileType {
	case linux.S_IFREG:
		li.inodeType = ENTRY_REGULAR
	case linux.S_IFDIR:
		li.inodeType = ENTRY_DIRECTORY
	case linux.S_IFLNK:
		li.inodeType = ENTRY_SYMLINK
	default:
		log.Debugf("Unknown file type %d", fileType)
		return linuxerr.EINVAL
	}

	nlink := uint32(inodeMetadata.ReferenceCount)
	if mode.FileType() == linux.ModeDirectory {
		nlink = 2
	}
	li.devMajor = devMajor
	li.devMinor = devMinor
	li.ino.Store(ino)
	li.mode.Store(uint32(mode))
	li.uid.Store(uint32(inodeMetadata.UID))
	li.gid.Store(uint32(inodeMetadata.GID))
	li.nlink.Store(nlink)
	li.blockSize.Store(4096)
	li.mtime.Store(inodeMetadata.MTime)
	li.ctime.Store(inodeMetadata.CTime)
	li.atime.Store(now)

	return nil
}
