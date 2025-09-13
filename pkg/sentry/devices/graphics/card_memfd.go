package graphics

import (
	"path"

	"golang.org/x/sys/unix"
)

func MemfdImplNative(name string, flags int) (int, error) {
	return unix.MemfdCreate(name, flags)
}

func MemfdImplDevShm(name string, flags int) (int, error) {
	realName := path.Join("/dev", "shm", name)
	return unix.Open(realName, flags|unix.O_CREAT|unix.O_RDWR|unix.O_LARGEFILE, 0o644)
}
