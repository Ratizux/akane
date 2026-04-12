package fakehostfs

import (
	"path"

	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	//"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	"gvisor.dev/gvisor/pkg/sentry/ktime"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

type FakehostfsInode struct {
	// fs/metadataBasePath/name must be set during initialization
	fs *FakehostfsImpl

	name             string

	valid bool

	dentry *kernfs.Dentry

	kernfs.InodeNotAnonymous

	kernfs.InodeWatches

	//lock
	locks vfs.FileLocks

	logical *logicalInode

	isRoot bool
}

func (i *FakehostfsInode) Readlink(ctx context.Context, mnt *vfs.Mount) (string, error) {
	log.Debugf("fakehostfs: ---> Readlink(): %d, %s", i.Ino(), i.name)
	defer log.Debugf("fakehostfs: <--- Readlink(): %d, %s", i.Ino(), i.name)
	target, err := i.fs.nativeFS.ReadSymlink(i.MetadataBasePath(), i.name)
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
		Mode:           uint16(0o777 | S_IFLNK),
		ReferenceCount: 1,
		CTime:          now,
		MTime:          now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, false)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterSymlink(i.MetadataBasePath(), i.name, name, target, newIno, i.Ino() == 1)
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
		Mode:           uint16(opts.Mode | S_IFREG),
		ReferenceCount: 1,
		CTime:          now,
		MTime:          now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, true)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterFile(i.MetadataBasePath(), i.name, name, newIno, i.Ino() == 1)
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
		Mode:           uint16(opts.Mode | S_IFDIR),
		ReferenceCount: 1,
		CTime:          now,
		MTime:          now,
	}
	newIno, err := nativeFS.FindAndRegisterInode(inodeMetadata, false)
	if err != nil {
		return nil, err
	}
	err = nativeFS.RegisterDirectory(i.MetadataBasePath(), i.name, name, newIno, i.Ino() == 1)
	if err != nil {
		return nil, err
	}
	inode, err := i.Lookup(ctx, name)
	if err != nil {
		return nil, err
	}
	return inode, nil
}

func (i *FakehostfsInode) NewLink(ctx context.Context, name string, target kernfs.Inode) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> NewLink(): %d", i.Ino())
	defer log.Debugf("fakehostfs: <--- NewLink(): %d", i.Ino())

	targetInode, ok := target.(*FakehostfsInode)
	if !ok || targetInode == nil {
		return nil, linuxerr.EIO
	}

	log.Debugf("NewLink: %s want to link child %s(%d) with a new name %s", i.name, targetInode.name, targetInode.Ino(), name)
	//return nil, linuxerr.EPERM
	err := targetInode.fs.nativeFS.RegisterNode(i.MetadataBasePath(),
						    i.name,
						    name,
						    targetInode.Ino(),
						    i.Ino() == 1)
	if err != nil {
		return nil, linuxerr.EIO
	}
	targetInode.logical.nlink.Add(1)

	err = targetInode.FlushStat()
	if err != nil {
		return nil, linuxerr.EIO
	}

	newInode := FakehostfsInode{
		fs:               i.fs,
		// metadataBasePath: i.metadataBasePath,
		name:             name,
		logical:          targetInode.logical,
		valid:            true,
	}
	return &newInode, nil
}

func (i *FakehostfsInode) Open(ctx context.Context, rp *vfs.ResolvingPath, d *kernfs.Dentry, opts vfs.OpenOptions) (*vfs.FileDescription, error) {
	log.Debugf("fakehostfs: ---> Open(): %d", i.Ino())
	defer log.Debugf("fakehostfs: <--- Open(): %d", i.Ino())
	fd := &FakehostfsFileDescription{
		inode: i,
		logical: i.logical,
	}
	if err := fd.Init(ctx, opts); err != nil {
		log.Debugf("Failed attempt of fd.Init()")
		return nil, err
	}
	if err := fd.vfsfd.Init(fd, opts.Flags, rp.Mount(), d.VFSDentry(), &vfs.FileDescriptionOptions{}); err != nil {
		log.Debugf("Failed attempt of fd.vfsfd.Init()")
		return nil, err
	}
	log.Debugf("Inode %d initialized successfully", i.Ino())
	return &fd.vfsfd, nil
}

func (i *FakehostfsInode) StatFS(ctx context.Context, fs *vfs.Filesystem) (linux.Statfs, error) {
	log.Debugf("StatFS() called on inode: %d", i.Ino())
	// TODO currently return fake statfs. should make sure if real statfs impl is feasible
	statfs := linux.Statfs {
		Type: linux.EXT_SUPER_MAGIC,
		BlockSize: 4096,
		FragmentSize: 4096,
		Blocks: 134742016,
		BlocksFree: 104857600,
		BlocksAvailable: 104857600,
		Files: 99999999,
		FilesFree: 114514,
		NameLength: 192,
	}
	return statfs, nil
}

func (i *FakehostfsInode) Keep() bool {
	log.Debugf("Keep() called on inode: %d", i.Ino())
	return true
}

func (i *FakehostfsInode) Valid(ctx context.Context, parent *kernfs.Dentry, name string) bool {
	log.Debugf("Valid() called on inode: %d", i.Ino())
	// TODO
	// figure out mechanism of inode invalidation
	return i.valid
}

func (i *FakehostfsInode) RegisterDentry(d *kernfs.Dentry) {
	i.dentry = d
	log.Debugf("RegisterDentry() called on inode: %d", i.Ino())
}

func (i *FakehostfsInode) UnregisterDentry(d *kernfs.Dentry) {
	i.dentry = nil
	log.Debugf("UnregisterDentry() called on inode: %d", i.Ino())
}

func (i *FakehostfsInode) HasChildren() bool {
	log.Debugf("HasChildren() called on inode: %d", i.Ino())
	return false
}

func (i *FakehostfsInode) IterDirents(ctx context.Context, mnt *vfs.Mount, callback vfs.IterDirentsCallback, offset, relOffset int64) (newOffset int64, err error) {
	log.Debugf("IterDirents() called on inode: %d", i.Ino())
	return offset, linuxerr.EINVAL
}

func (i *FakehostfsInode) Lookup(ctx context.Context, name string) (kernfs.Inode, error) {
	log.Debugf("fakehostfs: ---> Lookup(): %d, %s", i.Ino(), name)
	defer log.Debugf("fakehostfs: <--- Lookup(): %d, %s", i.Ino(), name)
	nativeFS := i.fs.nativeFS
	// regular files may not have existence in filesystem, check metadata instead

	var childMetadataPath string
	if i.isRoot {
		childMetadataPath = path.Join(i.MetadataBasePath(), "i"+name)
	} else {
		childMetadataPath = path.Join(i.MetadataBasePath(), "x"+i.name, "i"+name)
	}

	log.Debugf("child metadata path path is %s", childMetadataPath)
	childIno, err := nativeFS.GetInoFromPath(childMetadataPath)
	log.Debugf("Inode: %d", childIno)
	if err != nil {
		return nil, err
	}
	dentry := &FakehostfsDentry{}
	inode := FakehostfsInode{
		fs:               i.fs,
		name:             name,
	}
	err = inode.Init(ctx, i.fs.devMajor, i.fs.devMinor, childIno)
	if err != nil {
		return nil, err
	}
	dentry.Init(&i.fs.Filesystem, &inode)
	log.Debugf("initialized inode %d", childIno)
	return &inode, nil
}

func (i *FakehostfsInode) Rename(ctx context.Context, oldname string, newname string, child kernfs.Inode, dstDir kernfs.Inode) error {
	//dstDir
	dstInode, ok := dstDir.(*FakehostfsInode)
	if !ok {
		return linuxerr.EINVAL
	}
	//check if src exist
	nativeFS := i.fs.nativeFS
	//TODO invalidate old inode?

	selfMetadataBasePath := i.MetadataBasePath()
	targetMetadataBasePath := dstInode.MetadataBasePath()

	srcIno, err := nativeFS.GetIno(selfMetadataBasePath, i.name, oldname, i.isRoot)
	if err != nil {
		return linuxerr.EINVAL
	}
	srcMetadata, err := nativeFS.GetInoMetadata(srcIno)
	if err != nil {
		return linuxerr.ENOENT
	}
	//check if dest exist
	_, err = nativeFS.GetIno(targetMetadataBasePath, dstInode.name, newname, dstInode.isRoot)
	if err != nil {
		if err != linuxerr.ENOENT {
			return linuxerr.EEXIST
		}
	}
	//move
	if srcMetadata.Mode&S_IFMT == S_IFDIR {
		err = nativeFS.RenameDirectory(selfMetadataBasePath, i.name, oldname, i.isRoot, targetMetadataBasePath, dstInode.name, newname, dstInode.isRoot)
	} else if srcMetadata.Mode&S_IFMT == S_IFREG {
		err = nativeFS.RenameFile(selfMetadataBasePath, i.name, oldname, i.isRoot, targetMetadataBasePath, dstInode.name, newname, dstInode.isRoot)
	} else if srcMetadata.Mode&S_IFMT == S_IFLNK {
		err = nativeFS.RenameSymlink(selfMetadataBasePath, i.name, oldname, i.isRoot, targetMetadataBasePath, dstInode.name, newname, dstInode.isRoot)
	} else {
		return linuxerr.EINVAL
	}
	if err != nil {
		return linuxerr.EINVAL
	}

	childInode, ok := child.(*FakehostfsInode)
	if ok && childInode != nil {
		childInode.name = newname
		// childInode.valid = false
	} else {
		log.Debugf("Rename(): WARNING child inode is nil or unknown type")
	}

	return nil
}

func (i *FakehostfsInode) RmDir(ctx context.Context, name string, child kernfs.Inode) error {
	// TODO should consider non-empty dir, although kernfs docs claims that there is no need
	return i.Unlink(ctx, name, child)
}

func (i *FakehostfsInode) Unlink(ctx context.Context, name string, child kernfs.Inode) error {
	log.Debugf("Delete file: %s, parent Ino is %d", name, i.Ino())
	nativeFS := i.fs.nativeFS

	selfMetadataBasePath := i.MetadataBasePath()

	childIno, err := nativeFS.GetIno(selfMetadataBasePath, i.name, name, i.Ino() == 1)
	childInode, ok := child.(*FakehostfsInode)
	if !ok || childInode == nil {
		return linuxerr.EIO
	}

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
		err = nativeFS.DeleteFile(selfMetadataBasePath, i.name, name, i.Ino() == 1)
	} else if inodeMetadata.Mode&S_IFDIR != 0 {
		childType = ENTRY_DIRECTORY
		err = nativeFS.DeleteDirectory(selfMetadataBasePath, i.name, name, i.Ino() == 1)
	} else if inodeMetadata.Mode&S_IFLNK != 0 {
		childType = ENTRY_SYMLINK
		err = nativeFS.DeleteSymlink(selfMetadataBasePath, i.name, name, i.Ino() == 1)
	} else {
		return linuxerr.EINVAL
	}
	if err != nil {
		return linuxerr.EINVAL
	}

	newNlinks := childInode.logical.nlink.Load()
	if newNlinks == 0 {
		panic("unexpected nlink")
	}
	newNlinks --
	childInode.logical.nlink.Store(newNlinks)
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
