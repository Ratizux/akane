package fakehostfs

import (
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/ktime"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/log"
)

// DevMajor returns the device major number.
func (i *FakehostfsInode) DevMajor() uint32 {
	return i.logical.devMajor
}

// DevMinor returns the device minor number.
func (i *FakehostfsInode) DevMinor() uint32 {
	return i.logical.devMinor
}

// Ino returns the inode id.
func (i *FakehostfsInode) Ino() uint64 {
	return i.logical.ino.Load()
}

// UID implements Inode.UID.
func (i *FakehostfsInode) UID() auth.KUID {
	return auth.KUID(i.logical.uid.Load())
}

// GID implements Inode.GID.
func (i *FakehostfsInode) GID() auth.KGID {
	return auth.KGID(i.logical.gid.Load())
}

// Mode implements Inode.Mode.
func (i *FakehostfsInode) Mode() linux.FileMode {
	return linux.FileMode(i.logical.mode.Load())
}

// Links returns the link count.
func (i *FakehostfsInode) Links() uint32 {
	return i.logical.nlink.Load()
}

func (i *FakehostfsInode) CheckPermissions(_ context.Context, creds *auth.Credentials, ats vfs.AccessTypes) error {
	err := vfs.GenericCheckPermissions(
		creds,
		ats,
		i.Mode(),
		auth.KUID(i.logical.uid.Load()),
		auth.KGID(i.logical.gid.Load()),
	)
	if err == nil {
		log.Debugf("CheckPermissions() returned nil")
	} else {
		log.Debugf("CheckPermissions() returned %s", err.Error())
	}
	return err
}

// Init inode and create new logical-inode
func (i *FakehostfsInode) Init(ctx context.Context, devMajor uint32, devMinor uint32, ino uint64) error {
	inodeMetadata, err := i.fs.nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}

	i.logical = &logicalInode{}
	now := ktime.NowFromContext(ctx).Nanoseconds()

	err = i.logical.Init(now, devMajor, devMinor, ino, inodeMetadata)
	if err != nil {
		return linuxerr.EIO
	}

	i.valid = true
	return nil
}

// Sync stat of current inode to disk
func (i *FakehostfsInode) FlushStat() error {
	inodeMetadata := InodeMetadata {
		Mode: uint16(i.logical.mode.Load()),
		UID: i.logical.uid.Load(),
		GID: i.logical.gid.Load(),
		CTime: i.logical.ctime.Load(),
		MTime: i.logical.mtime.Load(),
	}
	if i.logical.inodeType == ENTRY_DIRECTORY {
		inodeMetadata.ReferenceCount = 1
	} else {
		inodeMetadata.ReferenceCount = uint16(i.logical.nlink.Load())
	}
	err := i.fs.nativeFS.SetInoMetadata(i.Ino(), inodeMetadata)
	if err != nil {
		return linuxerr.EIO
	}
	return nil
}

// SetStat implements Inode.SetStat.
func (i *FakehostfsInode) SetStatPrivate(ctx context.Context, fs *vfs.Filesystem, opts vfs.SetStatOptions) error {
	clearSID := false
	stat := opts.Stat
	if stat.Mask&linux.STATX_UID != 0 {
		i.logical.uid.Store(stat.UID)
		clearSID = true
	}
	if stat.Mask&linux.STATX_GID != 0 {
		i.logical.gid.Store(stat.GID)
		clearSID = true
	}
	if stat.Mask&linux.STATX_MODE != 0 {
		for {
			old := i.logical.mode.Load()
			ft := old & linux.S_IFMT
			newMode := ft | uint32(stat.Mode & ^uint16(linux.S_IFMT))
			if clearSID {
				newMode = vfs.ClearSUIDAndSGID(newMode)
			}
			if swapped := i.logical.mode.CompareAndSwap(old, newMode); swapped {
				clearSID = false
				break
			}
		}
	}

	// We may have to clear the SUID/SGID bits, but didn't do so as part of
	// STATX_MODE.
	if clearSID {
		for {
			old := i.logical.mode.Load()
			newMode := vfs.ClearSUIDAndSGID(old)
			if swapped := i.logical.mode.CompareAndSwap(old, newMode); swapped {
				break
			}
		}
	}

	now := ktime.NowFromContext(ctx).Nanoseconds()
	if stat.Mask&linux.STATX_ATIME != 0 {
		if stat.Atime.Nsec == linux.UTIME_NOW {
			stat.Atime = linux.NsecToStatxTimestamp(now)
		}
		i.logical.atime.Store(stat.Atime.ToNsec())
	}
	if stat.Mask&linux.STATX_MTIME != 0 {
		if stat.Mtime.Nsec == linux.UTIME_NOW {
			stat.Mtime = linux.NsecToStatxTimestamp(now)
		}
		i.logical.mtime.Store(stat.Mtime.ToNsec())
	}

	return i.FlushStat()
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
	if err := vfs.CheckSetStat(ctx, creds, &opts, i.Mode(), auth.KUID(i.logical.uid.Load()), auth.KGID(i.logical.gid.Load())); err != nil {
		return err
	}

	return i.SetStatPrivate(ctx, fs, opts)
}

func (i *FakehostfsInode) Stat(context.Context, *vfs.Filesystem, vfs.StatOptions) (linux.Statx, error) {
	/*inodeMetadata, err := i.fs.nativeFS.GetInoMetadata(i.Ino())
	if err != nil {
		return linux.Statx{}, err
	}*/
	stat := linux.Statx{}
	stat.Mask = linux.STATX_TYPE | linux.STATX_MODE | linux.STATX_UID | linux.STATX_GID | linux.STATX_INO | linux.STATX_NLINK | linux.STATX_ATIME | linux.STATX_MTIME | linux.STATX_CTIME
	if i.logical.inodeType == ENTRY_REGULAR {
		stat.Mask |= linux.STATX_SIZE
		objectSize, err := i.fs.nativeFS.InodeObjectSize(i.Ino())
		if err != nil {
			return linux.Statx{}, err
		}
		stat.Size = objectSize
	}
	stat.DevMajor = i.logical.devMajor
	stat.DevMinor = i.logical.devMinor
	stat.Ino = i.logical.ino.Load()
	stat.Mode = uint16(i.Mode())
	stat.UID = i.logical.uid.Load()
	stat.GID = i.logical.gid.Load()
	stat.Nlink = i.logical.nlink.Load()
	stat.Blksize = i.logical.blockSize.Load()
	stat.Atime = linux.NsecToStatxTimestamp(i.logical.atime.Load())
	stat.Mtime = linux.NsecToStatxTimestamp(i.logical.mtime.Load())
	stat.Ctime = linux.NsecToStatxTimestamp(i.logical.ctime.Load())
	return stat, nil
}
