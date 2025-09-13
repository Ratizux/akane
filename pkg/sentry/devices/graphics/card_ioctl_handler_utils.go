package graphics

import (
	"gvisor.dev/gvisor/pkg/log"
	//"gvisor.dev/gvisor/pkg/gohacks"
	//"gvisor.dev/gvisor/pkg/usermem"
	//"gvisor.dev/gvisor/pkg/sentry/arch"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/hostarch"
)

func SetUint32ByPointer(t *kernel.Task, pointer uint64, value uint32) {
	if pointer == 0 {
		return
	}
	valueBuffer := make([]byte, 4)
	hostarch.ByteOrder.PutUint32(valueBuffer, value)
	length, err := t.CopyOutBytes(hostarch.Addr(pointer), valueBuffer)
	if err != nil {
		log.Debugf("Failed to copy to string")
	}
	log.Debugf("Copied %d bytes to string", length)
}

func SetUint32ArrayByPointer(t *kernel.Task, pointer uint64, values []uint32) {
	if pointer == 0 {
		return
	}
	valueBuffer := make([]byte, 4*len(values))
	for index, value := range(values) {
		start := index*4
		end := (index+1)*4
		hostarch.ByteOrder.PutUint32(valueBuffer[start:end], value)
	}
	length, err := t.CopyOutBytes(hostarch.Addr(pointer), valueBuffer)
	if err != nil {
		log.Debugf("Failed to copy to string")
	}
	log.Debugf("Copied %d bytes to string", length)
}

func SetModeinfoArrayByPointer(t *kernel.Task, pointer uint64, modes []Mode) {
	if pointer == 0 {
		return
	}
	modeinfoSize := (&DrmModeModeinfo{}).SizeBytes()
	valueBuffer := make([]byte, modeinfoSize*len(modes))
	for index, mode := range(modes) {
		modeinfo := DrmModeModeinfo {
			hdisplay: mode.height,
			vdisplay: mode.width,
			vrefresh: mode.refreshRate,
		}
		copy(modeinfo.name[:], mode.name)
		start := index*modeinfoSize
		end := (index+1)*modeinfoSize
		modeinfo.MarshalBytes(valueBuffer[start:end])
	}
	length, err := t.CopyOutBytes(hostarch.Addr(pointer), valueBuffer)
	if err != nil {
		log.Debugf("Failed to copy to string")
	}
	log.Debugf("Copied %d bytes to string", length)
}
