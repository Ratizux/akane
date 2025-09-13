package graphics

import (
	"reflect"

	"gvisor.dev/gvisor/pkg/abi/linux"
)

func DRM_IO(nr uint32) uint32 {
	return linux.IO('d', nr)
}

func DRM_IOWR(nr uint32, target any) uint32 {
	t := reflect.TypeOf(target)
	return linux.IOWR('d', nr, uint32(t.Size()))
}
