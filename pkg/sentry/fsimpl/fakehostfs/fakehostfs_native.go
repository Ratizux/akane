package fakehostfs

import (
	"path"
	"strings"

	"golang.org/x/sys/unix"
	//"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/fsutil"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

type nativeFilesystem struct {
	hostPath string
	objectsPath string
	entriesPath string
}

type Metadata struct {
	Ino uint64
	Mode uint16
}

const (
	STAT_TYPE_MASK = 0170000
	STAT_PERMISSION_MASK = 0007777

	SEEK_SET = unix.SEEK_SET
	SEEK_CUR = unix.SEEK_CUR
	SEEK_END = unix.SEEK_END

	O_RDONLY = unix.O_RDONLY
	O_WRONLY = unix.O_WRONLY
	O_RDWR = unix.O_RDWR

	DT_DIR = unix.DT_DIR
	DT_REG = unix.DT_REG

	maxInode = 99999999

	S_IFREG = unix.S_IFREG
	S_IFDIR = unix.S_IFDIR
)

func CreatePathIfNotExist(targetPath string) error {
	if err := unix.Mkdir(targetPath, 0o700); err != nil {
		if err != unix.EEXIST {
			return err
		}
	}
	return nil
}

func CreateFile(targetPath string) error {
	if err := unix.Mknod(targetPath, unix.S_IFREG|0o700, 0); err != nil {
		return err
	}
	return nil
}

func (nativeFS *nativeFilesystem) Init(targetPath string) error {
	if path.IsAbs(targetPath) == false {
		return linuxerr.EINVAL
	}
	nativeFS.hostPath = targetPath
	nativeFS.objectsPath = path.Join(targetPath,"objects")
	err := unix.Mkdir(nativeFS.objectsPath,0o700)
	if err != nil && err != unix.EEXIST {
		return err
	}
	nativeFS.entriesPath = path.Join(targetPath,"entries")
	err = unix.Mkdir(nativeFS.entriesPath,0o700)
	if err != nil && err != unix.EEXIST {
		return err
	}

	if nativeFS.InodeValid(0) == false {
		// root node is not inialized. inode 0 is simply used as an indicator, though. root inode is 1.
		err := nativeFS.RegisterInodePrivate(0, InodeMetadata{
			Mode: 0,
		}, true)
		if err != nil {
			log.Debugf("Failed to register node 0")
			return err
		}
		err = nativeFS.RegisterInode(1, InodeMetadata{
			Mode: S_IFDIR|0o755,
		})
		if err != nil {
			log.Debugf("Failed to register node 1")
			return err
		}
	}
	/*
	_, reservedInode, err := GetInodePaths(0)
	if err != nil {
		return err
	}
	err := unix.Mknod()
	*/
	return nil
}

func (nativeFS *nativeFilesystem) Open(ino uint64, mode int) (int, error) {
	objectPath, _, err := nativeFS.GetInodePaths(ino)
	if err != nil {
		return -1, linuxerr.EINVAL
	}
	return unix.Open(objectPath, mode, 0)
}

func (nativeFS *nativeFilesystem) OpenDirectory(logicalPath string, mode int) (int, error) {
	log.Debugf("Warning: OpenDirectory mode is %d", mode)
	realPath := path.Join(nativeFS.entriesPath,logicalPath)
	return unix.Open(realPath, mode, 0)
}

func (nativeFS *nativeFilesystem) Close(hostfd int) error {
	return unix.Close(hostfd)
}

func (nativeFS *nativeFilesystem) Seek(hostfd int, offset int64, whence int) (int64,error) {
	return unix.Seek(hostfd,offset,whence)
}

func (nativeFS *nativeFilesystem) GetInnerDirents(hostfd int, workdir string) ([]vfs.Dirent,error) {
	log.Debugf("NativeFilesystem")
	hostDirents := []vfs.Dirent{}
	isDir := map[string]bool{}
	err := fsutil.ForEachDirent(hostfd,func(ino uint64, off int64, ftype uint8, name string, reclen uint16){
		dirent := vfs.Dirent{
			Name: name,
			Type: ftype,
			Ino: ino,
			NextOff: off,
		}
		if strings.HasPrefix(name,"x") {
			isDir[name[1:]] = true
		}
		hostDirents = append(hostDirents, dirent)
		log.Debugf("Ino: %d, Offset: %d, Type: %d, Name: %d",ino,off,ftype,name)
		log.Debugf("RecLen: %d",reclen)
	})
	dirents := []vfs.Dirent{}
	curNextOff := int64(3)
	for _,value := range hostDirents {
		if !strings.HasPrefix(value.Name,"i") {
			continue
		}
		dirent := vfs.Dirent{
			Name: value.Name[1:],
			NextOff: curNextOff,
		}
		if _, exists := isDir[value.Name[1:]]; exists {
			dirent.Type = DT_DIR
		} else {
			dirent.Type = DT_REG
		}
		ino, err := nativeFS.GetInoFromPath(path.Join(workdir, value.Name))
		if err != nil {
			return dirents, err
		}
		dirent.Ino = ino
		dirents = append(dirents, dirent)
		curNextOff++
	}
	if err != nil {
		log.Debugf("Failure getting directory entries")
		return dirents, err
	}
	return dirents, nil
}

// read-write call

func (nativeFS *nativeFilesystem) Read(hostfd int, dst []byte) (int64, error) {
	bytesRead, err := unix.Read(hostfd,dst)
	return int64(bytesRead), err
}

func (nativeFS *nativeFilesystem) PRead(hostfd int, dst []byte, offset int64) (int64, error) {
	bytesRead, err := unix.Pread(hostfd,dst,offset)
	return int64(bytesRead), err
}

func (nativeFS *nativeFilesystem) Write(hostfd int, src []byte) (int64, error) {
	bytesWritten, err := unix.Write(hostfd,src)
	return int64(bytesWritten), err
}

func (nativeFS *nativeFilesystem) PWrite(hostfd int, src []byte, offset int64) (int64, error) {
	bytesWritten, err := unix.Pwrite(hostfd,src,offset)
	return int64(bytesWritten), err
}


