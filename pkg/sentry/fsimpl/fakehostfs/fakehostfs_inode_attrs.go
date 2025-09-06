package fakehostfs

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/sentry/ktime"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"

)

// DevMajor returns the device major number.
func (i *FakehostfsInode) DevMajor() uint32 {
	return i.devMajor
}

// DevMinor returns the device minor number.
func (i *FakehostfsInode) DevMinor() uint32 {
	return i.devMinor
}

// Ino returns the inode id.
func (i *FakehostfsInode) Ino() uint64 {
	return i.ino.Load()
}

// UID implements Inode.UID.
func (i *FakehostfsInode) UID() auth.KUID {
	return auth.KUID(i.uid.Load())
}

// GID implements Inode.GID.
func (i *FakehostfsInode) GID() auth.KGID {
	return auth.KGID(i.gid.Load())
}

// Mode implements Inode.Mode.
func (i *FakehostfsInode) Mode() linux.FileMode {
	return linux.FileMode(i.mode.Load())
}

// Links returns the link count.
func (i *FakehostfsInode) Links() uint32 {
	return i.nlink.Load()
}

func (i *FakehostfsInode) CheckPermissions(_ context.Context, creds *auth.Credentials, ats vfs.AccessTypes) error {
	return vfs.GenericCheckPermissions(
		creds,
		ats,
		i.Mode(),
		auth.KUID(i.uid.Load()),
		auth.KGID(i.gid.Load()),
	)
}

func (i *FakehostfsInode) Init(ctx context.Context, devMajor uint32, devMinor uint32, ino uint64) error {
	inodeMetadata, err := i.fs.nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}
	mode := linux.FileMode(inodeMetadata.Mode)
	if mode.FileType() == 0 {
		panic(fmt.Sprintf("No file type specified in 'mode' for FakehostfsInode.Init(): mode=0%o", mode))
	}
	fileType := inodeMetadata.Mode&STAT_TYPE_MASK
	switch fileType {
		case linux.S_IFREG:
			i.inodeType = ENTRY_REGULAR
		case linux.S_IFDIR:
			i.inodeType = ENTRY_DIRECTORY
		case linux.S_IFLNK:
			i.inodeType = ENTRY_SYMLINK
		default:
			log.Debugf("Unknown file type %d",fileType)
			return linuxerr.EINVAL
	}

	nlink := uint32(inodeMetadata.ReferenceCount)
	if mode.FileType() == linux.ModeDirectory {
		nlink = 2
	}
	i.devMajor = devMajor
	i.devMinor = devMinor
	i.ino.Store(ino)
	i.mode.Store(uint32(mode))
	i.uid.Store(uint32(inodeMetadata.UID))
	i.gid.Store(uint32(inodeMetadata.GID))
	i.nlink.Store(nlink)
	i.blockSize.Store(hostarch.PageSize)
	i.mtime.Store(inodeMetadata.MTime)
	i.ctime.Store(inodeMetadata.CTime)
	now := ktime.NowFromContext(ctx).Nanoseconds()
	i.atime.Store(now)
	return nil
}

// SetStat implements Inode.SetStat.
func (i *FakehostfsInode) SetStatPrivate(ctx context.Context, fs *vfs.Filesystem, opts vfs.SetStatOptions) error {
	//TODO sync changes to root

	clearSID := false
	stat := opts.Stat
	if stat.Mask&linux.STATX_UID != 0 {
		i.uid.Store(stat.UID)
		clearSID = true
	}
	if stat.Mask&linux.STATX_GID != 0 {
		i.gid.Store(stat.GID)
		clearSID = true
	}
	if stat.Mask&linux.STATX_MODE != 0 {
		for {
			old := i.mode.Load()
			ft := old & linux.S_IFMT
			newMode := ft | uint32(stat.Mode & ^uint16(linux.S_IFMT))
			if clearSID {
				newMode = vfs.ClearSUIDAndSGID(newMode)
			}
			if swapped := i.mode.CompareAndSwap(old, newMode); swapped {
				clearSID = false
				break
			}
		}
	}

	// We may have to clear the SUID/SGID bits, but didn't do so as part of
	// STATX_MODE.
	if clearSID {
		for {
			old := i.mode.Load()
			newMode := vfs.ClearSUIDAndSGID(old)
			if swapped := i.mode.CompareAndSwap(old, newMode); swapped {
				break
			}
		}
	}

	now := ktime.NowFromContext(ctx).Nanoseconds()
	if stat.Mask&linux.STATX_ATIME != 0 {
		if stat.Atime.Nsec == linux.UTIME_NOW {
			stat.Atime = linux.NsecToStatxTimestamp(now)
		}
		i.atime.Store(stat.Atime.ToNsec())
	}
	if stat.Mask&linux.STATX_MTIME != 0 {
		if stat.Mtime.Nsec == linux.UTIME_NOW {
			stat.Mtime = linux.NsecToStatxTimestamp(now)
		}
		i.mtime.Store(stat.Mtime.ToNsec())
	}

	return nil
}

// SetStat implements Inode.SetStat.
func (i *FakehostfsInode) SetStat(ctx context.Context, fs *vfs.Filesystem, creds *auth.Credentials, opts vfs.SetStatOptions) error {
	if opts.Stat.Mask == 0 {
		return nil
	}

	// Note that not all fields are modifiable. For example, the file type and
	// inode numbers are immutable after node creation. Setting the size is often
	// allowed by kernfs files but does not do anything. If some other behavior is
	// needed, the embedder should consider extending SetStat.
	if opts.Stat.Mask&^(linux.STATX_MODE|linux.STATX_UID|linux.STATX_GID|linux.STATX_ATIME|linux.STATX_MTIME|linux.STATX_SIZE) != 0 {
		return linuxerr.EPERM
	}
	if opts.Stat.Mask&linux.STATX_SIZE != 0 && i.Mode().IsDir() {
		return linuxerr.EISDIR
	}
	if err := vfs.CheckSetStat(ctx, creds, &opts, i.Mode(), auth.KUID(i.uid.Load()), auth.KGID(i.gid.Load())); err != nil {
		return err
	}

	return i.SetStatPrivate(ctx,fs,opts)
}

func (i *FakehostfsInode) Stat(context.Context, *vfs.Filesystem, vfs.StatOptions) (linux.Statx, error) {
	/*inodeMetadata, err := i.fs.nativeFS.GetInoMetadata(i.Ino())
	if err != nil {
		return linux.Statx{}, err
	}*/
	stat := linux.Statx{}
	stat.Mask = linux.STATX_TYPE | linux.STATX_MODE | linux.STATX_UID | linux.STATX_GID | linux.STATX_INO | linux.STATX_NLINK | linux.STATX_ATIME | linux.STATX_MTIME | linux.STATX_CTIME
	if i.inodeType == ENTRY_REGULAR {
		stat.Mask |= linux.STATX_SIZE
		objectSize, err := i.fs.nativeFS.InodeObjectSize(i.Ino())
		if err != nil {
			return linux.Statx{}, err
		}
		stat.Size = objectSize
	}
	stat.DevMajor = i.devMajor
	stat.DevMinor = i.devMinor
	stat.Ino = i.ino.Load()
	stat.Mode = uint16(i.Mode())
	stat.UID = i.uid.Load()
	stat.GID = i.gid.Load()
	stat.Nlink = i.nlink.Load()
	stat.Blksize = i.blockSize.Load()
	stat.Atime = linux.NsecToStatxTimestamp(i.atime.Load())
	stat.Mtime = linux.NsecToStatxTimestamp(i.mtime.Load())
	stat.Ctime = linux.NsecToStatxTimestamp(i.ctime.Load())
	return stat, nil
}
