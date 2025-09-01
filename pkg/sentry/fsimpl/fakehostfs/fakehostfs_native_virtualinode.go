package fakehostfs

import (
	"os"
	"path"
	"strconv"


	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/log"
)

func AppendPath(basePath string, num int64) string {
	return path.Join(basePath, strconv.FormatInt(num, 10))
}

func (nativeFS *nativeFilesystem) GetFreeInode() (uint64, error) {
	RegularFileExists := func (targetPath string) (bool, error) {
		unixStat := &unix.Stat_t{}
		if err := unix.Stat(targetPath, unixStat); err != nil && err != unix.ENOENT {
			log.Debugf("Unable to stat() %s: %s", targetPath, err.Error())
			return false, err
		}
		if unixStat.Mode&STAT_TYPE_MASK == unix.S_IFREG {
			return true, nil
		}
		return false, nil
	}

	for part1 := range 100 {
		targetPath := AppendPath(nativeFS.objectsPath,int64(part1))
		err := CreatePathIfNotExist(targetPath)
		if err != nil {
			return 0, err
		}
		for part2 := range 100 {
			targetPath := AppendPath(targetPath,int64(part2))
			err := CreatePathIfNotExist(targetPath)
			if err != nil {
				return 0, err
			}
			for part3 := range 100 {
				targetPath := AppendPath(targetPath,int64(part3))
				err := CreatePathIfNotExist(targetPath)
				if err != nil {
					return 0, err
				}
				for part4 := range 100 {
					targetPath := path.Join(targetPath,"m"+strconv.FormatInt(int64(part4), 10))
					exists, err := RegularFileExists(targetPath)
					if err != nil {
						return 0, err
					}
					if !exists {
						// found
						return uint64(part1*1000000 + part2*10000 + part3*100 + part4), nil
					}
				}
			}
		}
	}
	return 0, linuxerr.EINVAL
}

func (nativeFS *nativeFilesystem) FindAndRegisterInode(inodeMetadata InodeMetadata) (uint64, error) {
	ino, err := nativeFS.GetFreeInode()
	if err != nil {
		return 0, linuxerr.EINVAL
	}
	err = nativeFS.RegisterInode(ino, inodeMetadata)
	if err != nil {
		return 0, linuxerr.EINVAL
	}
	return ino, nil
}

func (nativeFS *nativeFilesystem) CleanInode(ino uint64) error {
	objectPath, metadataPath, err := nativeFS.GetInodePaths(ino)
	inodeMetadata, err := nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}
	if inodeMetadata.ReferenceCount != 0 {
		return nil
	}
	if inodeMetadata.Mode&S_IFREG != 0 {
		err := os.Remove(objectPath)
		if err != nil {
			return linuxerr.EINVAL
		}
	}
	err = os.Remove(metadataPath)
	if err != nil {
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) IncreaseInodeReferenceCount(ino uint64) error {
	inodeMetadata, err := nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}
	inodeMetadata.ReferenceCount++
	err = nativeFS.SetInoMetadata(ino, inodeMetadata)
	if err != nil {
		return err
	}
	return nil
}

func (nativeFS *nativeFilesystem) DecreaseInodeReferenceCount(ino uint64) error {
	inodeMetadata, err := nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}
	inodeMetadata.ReferenceCount--
	err = nativeFS.SetInoMetadata(ino, inodeMetadata)
	if err != nil {
		return err
	}
	if inodeMetadata.ReferenceCount == 0 {
		nativeFS.CleanInode(ino)
	}
	return nil
}

func (nativeFS *nativeFilesystem) RegisterInodePrivate(ino uint64, inodeMetadata InodeMetadata, allowZero bool) error {
	if ino > maxInode {
		return linuxerr.EINVAL
	}
	if ino < 1 && !allowZero {
		return linuxerr.EINVAL
	}
	number := int64(ino)
	part4 := number%100
	number /= 100
	part3 := number%100
	number /= 100
	part2 := number%100
	number /= 100
	part1 := number

	targetPath := nativeFS.objectsPath
	for _, value := range []int64{part1,part2,part3} {
		targetPath = AppendPath(targetPath,value)
		err := CreatePathIfNotExist(targetPath)
		if err != nil {
			log.Debugf("Failure creating parent directories: %s",err.Error())
			return err
		}
	}
	objectPath := path.Join(targetPath,"o"+strconv.FormatInt(int64(part4), 10))
	metadataPath := path.Join(targetPath,"m"+strconv.FormatInt(int64(part4), 10))

	log.Debugf("Creating file %s",metadataPath)
	err := CreateFile(metadataPath)
	if err != nil {
		log.Debugf("Failure creating file: %s",err.Error())
		return linuxerr.EINVAL
	}

	if ino == 0 {
		return nil
	}

	if inodeMetadata.Mode&S_IFREG != 0 {
		err = CreateFile(objectPath)
		if err != nil {
			log.Debugf("Failure creating file: %s",err.Error())
			return linuxerr.EINVAL
		}
	}

	err = nativeFS.SetInoMetadata(ino, inodeMetadata)
	if err != nil {
		log.Debugf("Failure setting inode metadata: %s",err.Error())
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) RegisterInode(ino uint64, inodeMetadata InodeMetadata) error {
	return nativeFS.RegisterInodePrivate(ino, inodeMetadata, false)
}

func (nativeFS *nativeFilesystem) GetInodePaths(ino uint64) (string, string, error) {
	if ino > maxInode || ino < 1 {
		return "", "", linuxerr.EINVAL
	}
	number := int64(ino)
	part4 := strconv.FormatInt(number%100, 10)
	number /= 100
	part3 := strconv.FormatInt(number%100, 10)
	number /= 100
	part2 := strconv.FormatInt(number%100, 10)
	number /= 100
	part1 := strconv.FormatInt(number, 10)
	objectPath := path.Join(nativeFS.objectsPath,part1,part2,part3,"o"+part4)
	metadataPath := path.Join(nativeFS.objectsPath,part1,part2,part3,"m"+part4)
	return objectPath, metadataPath, nil
}

func (nativeFS *nativeFilesystem) InodeObjectSize(ino uint64) (uint64, error) {
	objectPath, _,err := nativeFS.GetInodePaths(ino)
	if err != nil {
		return 0, err
	}
	inodeMetadata, err := nativeFS.GetInoMetadata(ino)
	if err != nil {
		return 0, err
	}
	if inodeMetadata.Mode&S_IFDIR != 0 {
		return 1, nil
	}
	unixStat := &unix.Stat_t{}
	if err := unix.Stat(objectPath, unixStat); err != nil {
		log.Debugf("Unable to stat() %s: %s",objectPath,err.Error())
		return 0, err
	}
	return uint64(unixStat.Size), nil
}

func (nativeFS *nativeFilesystem) InodeValid(ino uint64) bool {
	_, metadataPath, err := nativeFS.GetInodePaths(ino)
	if err != nil {
		return false
	}
	unixStat := &unix.Stat_t{}
	if err := unix.Stat(metadataPath, unixStat); err != nil {
		log.Debugf("Unable to stat() %s: %s",metadataPath,err.Error())
		return false
	}
	if unixStat.Mode&STAT_TYPE_MASK == unix.S_IFREG {
		return true
	}
	return false
}
