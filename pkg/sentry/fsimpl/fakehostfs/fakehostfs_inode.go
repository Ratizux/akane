package fakehostfs

import (
	"path"

	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/context"
	//"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	"gvisor.dev/gvisor/pkg/sentry/ktime"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/atomicbitops"
)

type EntryType int

const (
	ENTRY_REGULAR EntryType = iota
	ENTRY_DIRECTORY
	ENTRY_SYMLINK
)

type FakehostfsInode struct {
	// fs/metadataBasePath/name must be set during initialization
	fs *FakehostfsImpl

	// these two fields are only applicable for directories
	metadataBasePath string
	name string

	inodeType EntryType

	dentry *kernfs.Dentry

	kernfs.InodeNotAnonymous

	kernfs.InodeWatches

	//lock
	locks vfs.FileLocks

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

func (i *FakehostfsInode) Readlink(ctx context.Context, mnt *vfs.Mount) (string, error) {
	log.Debugf("fakehostfs: ---> Readlink(): %d, %s", i.Ino(), i.name)
	defer log.Debugf("fakehostfs: <--- Readlink(): %d, %s", i.Ino(), i.name)
	target, err := i.fs.nativeFS.ReadSymlink(i.metadataBasePath, i.name)
	if err != nil {
		log.Debugf("Failed to get symlink")
		return "", linuxerr.EINVAL
	}
	return target, nil
}

func (i *FakehostfsInode) Getlink(ctx context.Context, mnt *vfs.Mount) (vfs.VirtualDentry, string, error) {
	log.Debugf("fakehostfs: ---> Getlink(): %d, %s", i.Ino(), i.name)
	defer log.Debugf("fakehostfs: <--- Getlink(): %d, %s", i.Ino(), i.name)
	//TODO support VirtualDentry
	target, err := i.Readlink(ctx, mnt)
	return vfs.VirtualDentry{}, target, err
}

func (i *FakehostfsInode) NewSymlink(ctx context.Context, name string, target string) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> NewSymlink(): %s", name)
	defer log.Debugf("fakehostfs: <--- NewSymlink(): %s", name)
	nativeFS := i.fs.nativeFS
	now := ktime.NowFromContext(ctx).Nanoseconds()
	inodeMetadata := InodeMetadata{
		Mode: uint16(0o777|S_IFLNK),
		ReferenceCount: 1,
		CTime: now,
		MTime: now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, false)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterSymlink(i.metadataBasePath,i.name,name,target,newIno,i.Ino()==1)
	if err != nil {
		return nil, err
	}
	inode, err := i.Lookup(ctx, name)
	if err != nil {
		return nil, err
	}
	return inode, nil
}

func (i *FakehostfsInode) NewNode(ctx context.Context, name string, opts vfs.MknodOptions) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> NewNode(): %s", name)
	defer log.Debugf("fakehostfs: <--- NewNode(): %s", name)
	return nil, linuxerr.ENOSYS
}

func (i *FakehostfsInode) NewFile(ctx context.Context, name string, opts vfs.OpenOptions) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> NewFile(): %s", name)
	defer log.Debugf("fakehostfs: <--- NewFile(): %s", name)
	nativeFS := i.fs.nativeFS
	now := ktime.NowFromContext(ctx).Nanoseconds()
	inodeMetadata := InodeMetadata{
		Mode: uint16(opts.Mode|S_IFREG),
		ReferenceCount: 1,
		CTime: now,
		MTime: now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, true)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterFile(i.metadataBasePath,i.name,name,newIno,i.Ino()==1)
	if err != nil {
		return nil, err
	}
	inode, err := i.Lookup(ctx, name)
	if err != nil {
		return nil, err
	}
	return inode, nil
}

func (i *FakehostfsInode) NewDir(ctx context.Context, name string, opts vfs.MkdirOptions) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> NewDir(): %s", name)
	defer log.Debugf("fakehostfs: <--- NewDir(): %s", name)
	nativeFS := i.fs.nativeFS
	now := ktime.NowFromContext(ctx).Nanoseconds()
	inodeMetadata := InodeMetadata{
		Mode: uint16(opts.Mode|S_IFDIR),
		ReferenceCount: 1,
		CTime: now,
		MTime: now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, false)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterDirectory(i.metadataBasePath,i.name,name,newIno,i.Ino()==1)
	if err != nil {
		return nil, err
	}
	inode, err := i.Lookup(ctx, name)
	if err != nil {
		return nil, err
	}
	return inode, nil
}

func (i *FakehostfsInode) NewLink(ctx context.Context, name string, opts kernfs.Inode) (kernfs.Inode, error) {
	return nil, linuxerr.ENOSYS
}

func (i *FakehostfsInode) Open(ctx context.Context, rp *vfs.ResolvingPath, d *kernfs.Dentry, opts vfs.OpenOptions) (*vfs.FileDescription, error) {
	log.Debugf("fakehostfs: ---> Open(): %d",i.Ino())
	defer log.Debugf("fakehostfs: <--- Open(): %d",i.Ino())
	fd := &FakehostfsFileDescription{inode: i}
	if err := fd.Init(ctx, opts); err != nil {
		log.Debugf("Failed attempt of fd.Init()")
		return nil, err
	}
	if err := fd.vfsfd.Init(fd, opts.Flags, rp.Mount(), d.VFSDentry(), &vfs.FileDescriptionOptions{}); err != nil {
		log.Debugf("Failed attempt of fd.vfsfd.Init()")
		return nil, err
	}
	log.Debugf("Inode %d initialized successfully",i.Ino())
	return &fd.vfsfd, nil
}

func (i *FakehostfsInode) StatFS(ctx context.Context, fs *vfs.Filesystem) (linux.Statfs, error) {
	log.Debugf("StatFS() called on inode: %d",i.Ino())
	statfs := linux.Statfs{}
	return statfs, linuxerr.EINVAL
}

func (i *FakehostfsInode) Keep() bool {
	log.Debugf("Keep() called on inode: %d",i.Ino())
	return true
}

func (i *FakehostfsInode) Valid(ctx context.Context, parent *kernfs.Dentry, name string) bool {
	log.Debugf("Valid() called on inode: %d",i.Ino())
	// TODO
	// figure out mechanism of inode invalidation
	return true
}

func (i *FakehostfsInode) RegisterDentry(d *kernfs.Dentry) {
	i.dentry = d
	log.Debugf("RegisterDentry() called on inode: %d",i.Ino())
}

func (i *FakehostfsInode) UnregisterDentry(d *kernfs.Dentry) {
	i.dentry = nil
	log.Debugf("UnregisterDentry() called on inode: %d",i.Ino())
}

func (i *FakehostfsInode) HasChildren() bool {
	log.Debugf("HasChildren() called on inode: %d",i.Ino())
	return false
}

func (i *FakehostfsInode) IterDirents(ctx context.Context, mnt *vfs.Mount, callback vfs.IterDirentsCallback, offset, relOffset int64) (newOffset int64, err error) {
	log.Debugf("IterDirents() called on inode: %d",i.Ino())
	return offset, linuxerr.EINVAL
}

func (i *FakehostfsInode) Lookup(ctx context.Context, name string) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> Lookup(): %d, %s",i.Ino(),name)
	defer log.Debugf("fakehostfs: <--- Lookup(): %d, %s",i.Ino(),name)
	nativeFS := i.fs.nativeFS
	// regular files may not have existence in filesystem, check metadata instead
	childMetadataPath := path.Join(i.metadataBasePath,"x"+i.name,"i"+name)
	childMetadataBasePath := path.Join(i.metadataBasePath,"x"+i.name)
	if i.Ino() == 1 {
		childMetadataPath = path.Join(i.metadataBasePath,"i"+name)
		childMetadataBasePath = path.Join(i.metadataBasePath)
	}

	log.Debugf("child metadata path path is %s",childMetadataPath)
	childIno, err := nativeFS.GetInoFromPath(childMetadataPath)
	log.Debugf("Inode: %d",childIno)
	if err != nil {
		return nil, err
	}
	dentry := &FakehostfsDentry{}
	inode := FakehostfsInode{
		fs: i.fs,
		metadataBasePath: childMetadataBasePath,
		name: name,
	}
	err = inode.Init(ctx, i.fs.devMajor, i.fs.devMinor, childIno)
	if err != nil {
		return nil, err
	}
	dentry.Init(&i.fs.Filesystem, &inode)
	log.Debugf("initialized inode %d",childIno)
	return &inode, nil
}

func (i *FakehostfsInode) Rename(ctx context.Context, oldname string, newname string, child, dstDir kernfs.Inode) error {
	//dstDir
	dstInode, ok := dstDir.(*FakehostfsInode)
	if !ok {
		return linuxerr.EINVAL
	}
	//check if src exist
	nativeFS := i.fs.nativeFS
	//TODO invalidate old inode
	srcIno, err := nativeFS.GetIno(i.metadataBasePath,i.name,oldname,i.Ino()==1)
	if err != nil {
		return linuxerr.EINVAL
	}
	srcMetadata, err := nativeFS.GetInoMetadata(srcIno)
	if err != nil {
		return linuxerr.ENOENT
	}
	//check if dest exist
	_, err = nativeFS.GetIno(dstInode.metadataBasePath,dstInode.name,newname,dstInode.Ino()==1)
	if err != nil {
		if err != linuxerr.ENOENT {
			return linuxerr.EEXIST
		}
	}
	//move
	if srcMetadata.Mode&S_IFDIR != 0 {
		err = nativeFS.RenameDirectory(i.metadataBasePath,i.name,oldname,i.Ino()==1,dstInode.metadataBasePath,dstInode.name,newname,dstInode.Ino()==1)
	} else if srcMetadata.Mode&S_IFREG != 0 {
		err = nativeFS.RenameNode(i.metadataBasePath,i.name,oldname,i.Ino()==1,dstInode.metadataBasePath,dstInode.name,newname,dstInode.Ino()==1)
	} else {
		return linuxerr.EINVAL
	}
	if err != nil {
		return linuxerr.EINVAL
	}
	return nil
}

func (i *FakehostfsInode) RmDir(ctx context.Context, name string, child kernfs.Inode) error {
	return linuxerr.EINVAL
}

func (i *FakehostfsInode) Unlink(ctx context.Context, name string, child kernfs.Inode) error {
	log.Debugf("Delete file: %s, parent Ino is %d",name, i.Ino())
	nativeFS := i.fs.nativeFS
	childIno, err := nativeFS.GetIno(i.metadataBasePath,i.name,name,i.Ino()==1)
	if err != nil {
		return linuxerr.EINVAL
	}
	inodeMetadata, err := nativeFS.GetInoMetadata(childIno)
	if err != nil {
		return linuxerr.EINVAL
	}
	var childType EntryType
	// FILE TYPE
	if inodeMetadata.Mode&S_IFREG != 0 {
		childType = ENTRY_REGULAR
		err = nativeFS.DeleteFile(i.metadataBasePath,i.name,name,i.Ino()==1)
	} else if inodeMetadata.Mode&S_IFDIR != 0 {
		childType = ENTRY_DIRECTORY
		err = nativeFS.DeleteDirectory(i.metadataBasePath,i.name,name,i.Ino()==1)
	} else if inodeMetadata.Mode&S_IFLNK != 0 {
		childType = ENTRY_SYMLINK
		err = nativeFS.DeleteSymlink(i.metadataBasePath,i.name,name,i.Ino()==1)
	} else {
		return linuxerr.EINVAL
	}
	if err != nil {
		return linuxerr.EINVAL
	}
	if childType == ENTRY_REGULAR {
		err = nativeFS.DecreaseInodeReferenceCount(childIno, true)
	} else {
		err = nativeFS.DecreaseInodeReferenceCount(childIno, false)
	}
	if err != nil {
		return linuxerr.EINVAL
	}
	return nil
}
