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

func (nativeFS *nativeFilesystem) GetFreeInode(lastFreeInode uint64) (uint64, error) {
	RegularFileExists := func(targetPath string) (bool, error) {
		log.Debugf("Attempt to stat %s", targetPath)

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

	part1 := lastFreeInode / 1000000
	part2 := lastFreeInode / 10000 % 100
	part3 := lastFreeInode / 100 % 100
	part4 := lastFreeInode % 100

	lastFreeObjectPath := path.Join(nativeFS.objectsPath,
		strconv.FormatInt(int64(part1), 10),
		strconv.FormatInt(int64(part2), 10),
		strconv.FormatInt(int64(part3), 10),
		"m"+strconv.FormatInt(int64(part4), 10))
	exists, err := RegularFileExists(lastFreeObjectPath)
	if err != nil {
		return 0, err
	}
	if !exists {
		// nativeFS.lastFreeInode is not used

		return lastFreeInode, nil
	}

	// nativeFS.lastFreeInode is used, find another
	for i := part1; i < 100; i++ {
		targetPath := AppendPath(nativeFS.objectsPath, int64(i))
		err := CreatePathIfNotExist(targetPath)
		if err != nil {
			return 0, err
		}

		for j := part2; j < 100; j++ {
			targetPath := AppendPath(targetPath, int64(j))
			err := CreatePathIfNotExist(targetPath)
			if err != nil {
				return 0, err
			}

			for k := part3; k < 100; k++ {
				targetPath := AppendPath(targetPath, int64(k))
				err := CreatePathIfNotExist(targetPath)
				if err != nil {
					return 0, err
				}

				entries, err := os.ReadDir(targetPath)
				if err != nil {
					return 0, err
				}

				var used_entry_map [100]bool

				for index, value := range entries {
					_ = index
					entry_id, err := strconv.Atoi(value.Name()[1:])
					if err != nil {
						return 0, err
					}
					// TODO index range check?
					used_entry_map[entry_id] = true
				}

				var entry_id uint64
				for entry_id = 0; entry_id < 100; entry_id++ {
					if used_entry_map[entry_id] == true {
						continue
					}
					free_inode := uint64(i*1000000 + j*10000 + k*100 + entry_id)
					return free_inode, nil
				}
			}
		}
	}
	return 0, linuxerr.EINVAL
}

func (nativeFS *nativeFilesystem) FindAndRegisterInode(inodeMetadata InodeMetadata, hasObject bool) (uint64, error) {
	ino, err := nativeFS.GetFreeInode(nativeFS.lastFreeInode)
	if err != nil {
		return 0, linuxerr.EINVAL
	}
	log.Debugf("FindAndRegisterInode(): got free inode %d", ino)
	err = nativeFS.RegisterInode(ino, inodeMetadata, hasObject)
	if err != nil {
		nativeFS.lastFreeInode = ino
		return 0, linuxerr.EINVAL
	}
	nativeFS.lastFreeInode = ino + 1
	return ino, nil
}

func (nativeFS *nativeFilesystem) CleanInode(ino uint64, hasObject bool) error {
	objectPath, metadataPath, err := nativeFS.GetInodePaths(ino)
	inodeMetadata, err := nativeFS.GetInoMetadata(ino)
	if err != nil {
		return err
	}
	if inodeMetadata.ReferenceCount != 0 {
		return nil
	}
	if hasObject {
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

func (nativeFS *nativeFilesystem) DecreaseInodeReferenceCount(ino uint64, hasObject bool) error {
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
		nativeFS.CleanInode(ino, hasObject)
	}
	return nil
}

func (nativeFS *nativeFilesystem) RegisterInodePrivate(ino uint64, inodeMetadata InodeMetadata, hasObject bool, allowZero bool) error {
	if ino > maxInode {
		return linuxerr.EINVAL
	}
	if ino < 1 && !allowZero {
		return linuxerr.EINVAL
	}
	number := int64(ino)
	part4 := number % 100
	number /= 100
	part3 := number % 100
	number /= 100
	part2 := number % 100
	number /= 100
	part1 := number

	targetPath := nativeFS.objectsPath
	for _, value := range []int64{part1, part2, part3} {
		targetPath = AppendPath(targetPath, value)
		err := CreatePathIfNotExist(targetPath)
		if err != nil {
			log.Debugf("Failure creating parent directories: %s", err.Error())
			return err
		}
	}
	objectPath := path.Join(targetPath, "o"+strconv.FormatInt(int64(part4), 10))
	metadataPath := path.Join(targetPath, "m"+strconv.FormatInt(int64(part4), 10))

	log.Debugf("Creating file %s", metadataPath)
	err := CreateFile(metadataPath)
	if err != nil {
		log.Debugf("Failure creating file: %s", err.Error())
		return linuxerr.EINVAL
	}

	if ino == 0 {
		return nil
	}

	if hasObject {
		err = CreateFile(objectPath)
		if err != nil {
			log.Debugf("Failure creating file: %s", err.Error())
			return linuxerr.EINVAL
		}
	}

	err = nativeFS.SetInoMetadata(ino, inodeMetadata)
	if err != nil {
		log.Debugf("Failure setting inode metadata: %s", err.Error())
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) RegisterInode(ino uint64, inodeMetadata InodeMetadata, hasObject bool) error {
	return nativeFS.RegisterInodePrivate(ino, inodeMetadata, hasObject, false)
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
	objectPath := path.Join(nativeFS.objectsPath, part1, part2, part3, "o"+part4)
	metadataPath := path.Join(nativeFS.objectsPath, part1, part2, part3, "m"+part4)
	return objectPath, metadataPath, nil
}

func (nativeFS *nativeFilesystem) InodeObjectSize(ino uint64) (uint64, error) {
	objectPath, _, err := nativeFS.GetInodePaths(ino)
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
		log.Debugf("Unable to stat() %s: %s", objectPath, err.Error())
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
		log.Debugf("Unable to stat() %s: %s", metadataPath, err.Error())
		return false
	}
	if unixStat.Mode&STAT_TYPE_MASK == unix.S_IFREG {
		return true
	}
	return false
}
