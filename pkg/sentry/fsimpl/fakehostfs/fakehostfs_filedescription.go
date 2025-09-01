package fakehostfs

import (
	"path"

	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/usermem"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
)

type FakehostfsFileDescription struct {
	vfs.FileDescriptionDefaultImpl
	vfs.NoLockFD

	fdType FileDescriptionType
	inode *FakehostfsInode
	hostfd int
	hostfdOpen bool
	//filePointer int64 //a.k.a. offset, can be changed using seek()
	// for regular files, should be consistent with host fd
	// for directories, simulate behavior instead. do not track host fd offset
	virtualOffset int64

	direntsCache []vfs.Dirent
	direntsCacheValid bool

	vfsfd vfs.FileDescription
}

type FileDescriptionType int

const (
	FDTYPE_REGULAR FileDescriptionType = iota
	FDTYPE_DIRECTORY
)
// O_RDONLY, O_WRONLY, O_RDWR, O_APPEND, , O_DIRECT, O_DSYNC,
// , , O_SYNC, , and
// O_TRUNC



func (fd *FakehostfsFileDescription) Init(ctx context.Context, opts vfs.OpenOptions) error {
	log.Debugf("Attempt to open FD associated with inode %d",fd.inode.Ino())
	log.Debugf("Flag is %d",opts.Flags)

	nativeFS := fd.inode.fs.nativeFS

	// handle open flags
	flags := opts.Flags

	if flags&linux.O_TMPFILE != 0 {
		// TODO support this
		log.Debugf("O_TMPFILE not implemented")
		return linuxerr.EOPNOTSUPP
	}

	//TODO handle O_NOATIME, O_NOCTTY, O_NONBLOCK

	//flags |= linux.O_NOATIME can improve performance?

	if flags&linux.O_CREAT != 0 {
		log.Debugf("Removed O_CREAT flag, this should be handled by kernfs")
		flags ^= linux.O_CREAT
		if flags&linux.O_EXCL != 0 {
			flags ^= linux.O_EXCL
			log.Debugf("Removed O_EXCL flag, this should be handled by kernfs")
		}
	}

	log.Debugf("Detect file type...")
	metadata, err := nativeFS.GetInoMetadata(fd.inode.Ino())
	if err != nil {
		log.Debugf("Unable to get metadata")
		return linuxerr.EINVAL
	}
	fileType := metadata.Mode&STAT_TYPE_MASK
	switch fileType {
		case linux.S_IFREG:
			fd.fdType = FDTYPE_REGULAR
		case linux.S_IFDIR:
			fd.fdType = FDTYPE_DIRECTORY
		default:
			log.Debugf("Unknown file type %d",fileType)
			return linuxerr.EINVAL
	}


	if fd.fdType == FDTYPE_REGULAR {
		fd.hostfd, err = nativeFS.Open(fd.inode.Ino(), int(flags))
		if err != nil {
			log.Debugf("Failed to open FD")
			return err
		}
		fd.hostfdOpen = true
	} else if fd.fdType == FDTYPE_DIRECTORY {
		directoryPath := path.Join(fd.inode.metadataBasePath,"x"+fd.inode.name)
		if fd.inode.Ino() == 1 {
			directoryPath = fd.inode.metadataBasePath
		}

		fd.hostfd, err = nativeFS.OpenDirectory(directoryPath, int(flags))
		if err != nil {
			log.Debugf("Failed to open FD associated with %s",directoryPath)
			return err
		}
		fd.hostfdOpen = true
	} else {
		return linuxerr.EINVAL
	}



	return nil
}

func (fd *FakehostfsFileDescription) Read(ctx context.Context, dst usermem.IOSequence, opts vfs.ReadOptions) (int64, error) {
	if opts.Flags != 0 {
		return 0, linuxerr.EOPNOTSUPP
	}
	if fd.fdType != FDTYPE_REGULAR {
		panic("FD is not a regular file!")
	}
	bufferSize := dst.NumBytes()
	log.Debugf("Got buffer: size %d",bufferSize)
	buffer := make([]byte,bufferSize)
	bytesRead, err := fd.inode.fs.nativeFS.Read(fd.hostfd,buffer)
	log.Debugf("Bytes read: %d",bytesRead)
	if err != nil {
		log.Debugf("Failure calling host Read(): %s",err.Error())
		return bytesRead, err
	}
	bytesCopied, err := dst.CopyOut(ctx,buffer)
	log.Debugf("Bytes copied: %d",bytesCopied)
	if err != nil {
		log.Debugf("Failure calling usermem CopyOut(): %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	if bytesRead > int64(bytesCopied) {
		log.Debugf("Partial copy")
		return 0, linuxerr.EINVAL
	}
	//fd.filePointer += bytesRead
	return bytesRead, nil
}

func (fd *FakehostfsFileDescription) PRead(ctx context.Context, dst usermem.IOSequence, offset int64, opts vfs.ReadOptions) (int64, error) {
	if opts.Flags != 0 {
		return 0, linuxerr.EOPNOTSUPP
	}
	bufferSize := dst.NumBytes()
	log.Debugf("Got buffer: size %d",bufferSize)
	buffer := make([]byte,bufferSize)
	bytesRead, err := fd.inode.fs.nativeFS.PRead(fd.hostfd,buffer,offset)
	log.Debugf("Bytes read: %d",bytesRead)
	if err != nil {
		log.Debugf("Failure calling host PRead(): %s",err.Error())
		return bytesRead, err
	}
	bytesCopied, err := dst.CopyOut(ctx,buffer)
	log.Debugf("Bytes copied: %d",bytesCopied)
	if err != nil {
		log.Debugf("Failure calling usermem CopyOut(): %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	if bytesRead > int64(bytesCopied) {
		log.Debugf("Partial copy")
		return 0, linuxerr.EINVAL
	}
	return bytesRead, nil
}

func (fd *FakehostfsFileDescription) Write(ctx context.Context, dst usermem.IOSequence, opts vfs.WriteOptions) (int64, error) {
	if fd.fdType != FDTYPE_REGULAR {
		panic("FD is not a regular file!")
	}
	if opts.Flags != 0 {
		return 0, linuxerr.EOPNOTSUPP
	}
	bufferSize := dst.NumBytes()
	log.Debugf("Got buffer: size %d",bufferSize)
	buffer := make([]byte,bufferSize)
	bytesCopied, err := dst.CopyIn(ctx,buffer)
	log.Debugf("Bytes copied: %d",bytesCopied)
	if err != nil {
		log.Debugf("Failure calling usermem CopyIn(): %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	if bufferSize > int64(bytesCopied) {
		log.Debugf("Partial copy")
		return 0, linuxerr.EINVAL
	}
	bytesWritten, err := fd.inode.fs.nativeFS.Write(fd.hostfd,buffer)
	log.Debugf("Bytes written: %d",bytesWritten)
	if err != nil {
		log.Debugf("Failure calling host Write(): %s",err.Error())
		return bytesWritten, err
	}
	//fd.filePointer += bytesWritten
	return bytesWritten, nil
}

func (fd *FakehostfsFileDescription) Seek(ctx context.Context, offset int64, whence int32) (int64, error) {
	if fd.fdType != FDTYPE_REGULAR {
		return 0, linuxerr.EINVAL
	}
	return fd.inode.fs.nativeFS.Seek(fd.hostfd,offset,int(whence))
}

func (fd *FakehostfsFileDescription) PWrite(ctx context.Context, dst usermem.IOSequence, offset int64, opts vfs.WriteOptions) (int64, error) {
	if fd.fdType != FDTYPE_REGULAR {
		panic("FD is not a regular file!")
	}
	if opts.Flags != 0 {
		return 0, linuxerr.EOPNOTSUPP
	}
	bufferSize := dst.NumBytes()
	log.Debugf("Got buffer: size %d",bufferSize)
	buffer := make([]byte,bufferSize)

	bytesCopied, err := dst.CopyIn(ctx,buffer)
	log.Debugf("Bytes copied: %d",bytesCopied)
	if err != nil {
		log.Debugf("Failure calling usermem CopyIn(): %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	if bufferSize > int64(bytesCopied) {
		log.Debugf("Partial copy")
		return 0, linuxerr.EINVAL
	}

	bytesWritten, err := fd.inode.fs.nativeFS.PWrite(fd.hostfd,buffer,offset)
	log.Debugf("Bytes written: %d",bytesWritten)
	if err != nil {
		log.Debugf("Failure calling host PWrite(): %s",err.Error())
		return bytesWritten, err
	}

	return bytesWritten, nil
}

func (fd *FakehostfsFileDescription) Release(ctx context.Context) {
	if fd.hostfdOpen == false {
		return
	}
	err := fd.inode.fs.nativeFS.Close(fd.hostfd)
	log.Debugf("Closing FD associated with base path %s",fd.inode.metadataBasePath)
	if err != nil {
		panic("Unable to close hostfd")
	}
}

func (fd *FakehostfsFileDescription) SetStat(ctx context.Context, opts vfs.SetStatOptions) error {
	log.Debugf("SetStat() called on FD, inode: %d",fd.inode.Ino())
	return fd.inode.SetStatPrivate(ctx,fd.inode.fs.VFSFilesystem(),opts)
}

func (fd *FakehostfsFileDescription) Stat(ctx context.Context, opts vfs.StatOptions) (linux.Statx, error) {
	log.Debugf("Stat() called on FD, inode: %d",fd.inode.Ino())
	return fd.inode.Stat(ctx,fd.inode.fs.VFSFilesystem(),opts)
}

func (fd *FakehostfsFileDescription) UpdateDirentsCache() error {
	nativeFS := fd.inode.fs.nativeFS
	newOffset,err := nativeFS.Seek(fd.hostfd,0,SEEK_SET)
	if err != nil {
		return err
	}
	if newOffset != 0 {
		log.Debugf("File pointer of directory is %d, should be 0...",newOffset)
		return linuxerr.EINVAL
	}
	workdir := path.Join(fd.inode.metadataBasePath, "x"+fd.inode.name)
	if fd.inode.Ino() == 1 {
		workdir = fd.inode.metadataBasePath
	}
	fd.direntsCache, err = nativeFS.GetInnerDirents(fd.hostfd, workdir)
	if err != nil {
		return err
	}
	fd.direntsCacheValid = true
	return nil
}

func (fd *FakehostfsFileDescription) IterDirents(ctx context.Context, cb vfs.IterDirentsCallback) error {
	dirents := []vfs.Dirent{}
	log.Debugf("IterDirents() called on FD")
	if fd.direntsCacheValid == false {
		log.Debugf("direntsCache miss")
		err := fd.UpdateDirentsCache()
		if err != nil {
			return err
		}
	} else {
		log.Debugf("direntsCache hit")
	}
	currentNodeID := fd.inode.Ino()
	// handle current directory
	if fd.virtualOffset == 0 {
		dirents = append(dirents,vfs.Dirent{
			Name: ".",
			Type: linux.DT_DIR,
			Ino: currentNodeID,
			NextOff: 1,
		})
		fd.virtualOffset++
	}
	// handle parent directory
	if fd.virtualOffset == 1 {
		if fd.inode.fs.rootNodeID == currentNodeID {
			// current directory is filesystem root
			dirents = append(dirents,vfs.Dirent{
				Name: "..",
				Type: linux.DT_DIR,
				Ino: currentNodeID,
				NextOff: 2,
			})
		} else {
			// current directory is not filesystem root
			parentDentry := fd.inode.dentry.Parent()
			if parentDentry == nil {
				return linuxerr.EINVAL
			}
			parentInode, ok := parentDentry.Inode().(*FakehostfsInode)
			if !ok {
				return linuxerr.EINVAL
			}
			dirents = append(dirents,vfs.Dirent{
				Name: "..",
				Type: linux.DT_DIR,
				Ino: parentInode.Ino(),
				NextOff: 2,
			})
		}
		fd.virtualOffset++
	}
	// add
	for _,value := range fd.direntsCache {
		if fd.virtualOffset+1 == value.NextOff {
			dirents = append(dirents,value)
		}
		fd.virtualOffset++
	}
	//dirents = append(dirents,realDirents...)
	for index,value := range dirents {
		err := cb.Handle(value)
		if err != nil {
			fd.virtualOffset = int64(index)
			return err
		}
	}
	return nil
}

func (fd *FakehostfsFileDescription) Sync(ctx context.Context) error {
	return nil
}
