// This file implements methods related to vfs.FilesystemType
package fakehostfs

import (
	//"strconv"

	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/context"
)

// FakehostfsType implements vfs.FilesystemType.
type FakehostfsType struct{
	fsImpl *FakehostfsImpl
}

// Name implements vfs.FilesystemType.Name.
func (fsType FakehostfsType) Name() string {
	return Name
}

// Release implements vfs.FilesystemType.Release.
func (fsType FakehostfsType) Release(ctx context.Context) {}

// GetFilesystem implements vfs.FilesystemType.GetFilesystem.
func (fsType FakehostfsType) GetFilesystem(ctx context.Context, vfsObj *vfs.VirtualFilesystem, creds *auth.Credentials, source string, opts vfs.GetFilesystemOptions) (*vfs.Filesystem, *vfs.Dentry, error) {
	ctx.Debugf("fakehostfs.GetFilesystem is called!\n");
	devMinor, err := vfsObj.GetAnonBlockDevMinor()
	if err != nil {
		ctx.Debugf("GetAnonBlockDevMinor failed\n");
		return nil, nil, err
	}

	var realSource string
	mopts := vfs.GenericParseMountOptions(opts.Data)
	if str, ok := mopts["source"]; ok {
		log.Debugf("source is %s", str)
		realSource = str
	} else {
		return nil, nil, linuxerr.EINVAL
	}
	/*maxCachedDentries := defaultMaxCachedDentries
	if str, ok := mopts["dentry_cache_limit"]; ok {
		delete(mopts, "dentry_cache_limit")
		maxCachedDentries, err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			ctx.Warningf("fakehostfs.FakehostfsType.GetFilesystem: invalid dentry cache limit: dentry_cache_limit=%s", str)
			return nil, nil, linuxerr.EINVAL
		}
	}*/

	fs := &FakehostfsImpl{
		devMajor:1,
		devMinor: devMinor,
		nativeFS: &nativeFilesystem{},
	}
	var rootNodeID uint64 = 1
	fs.nativeFS.Init(realSource)
	fs.rootNodeID = rootNodeID

	fsType.fsImpl=fs

	fs.VFSFilesystem().Init(vfsObj, fsType, fs)

	fs.root = &FakehostfsDentry{}

	inode := &FakehostfsInode{
		fs: fs,
		metadataBasePath: "/",
		name: "",
	}
	inode.Init(ctx,fs.devMajor,fs.devMinor,rootNodeID)

	fs.root.InitRoot(&fs.Filesystem, inode)

	ctx.Debugf("Fakehostfs initialized successfully.\n");
	return fs.VFSFilesystem(), fs.root.VFSDentry(), nil
}
