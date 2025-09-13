package graphics

import (
	"gvisor.dev/gvisor/pkg/sentry/memmap"
	//"gvisor.dev/gvisor/pkg/usermem"
	//"gvisor.dev/gvisor/pkg/sentry/arch"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/safemem"
	//"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/log"
	"golang.org/x/sys/unix"
)

type MemmapFile struct {
	memmap.DefaultMemoryType
	memmap.NoBufferedIOFallback

	buffer []byte
	hostfd int
}

func (mf *MemmapFile) Init(size uint64) {
	log.Debugf("MemmapFile: ---> Init()")
	defer log.Debugf("MemmapFile: <--- Init()")
	var err error
	mf.hostfd, err = MemfdImplDevShm("framebuffer", 0)
	if err != nil {
		panic("failed to call memfd: "+err.Error())
	}
	err = unix.Ftruncate(mf.hostfd, 1280*720*4)
	if err != nil {
		panic("failed to change memfd size")
	}
	mf.buffer, err = unix.Mmap(mf.hostfd, 0, 1280*720*4, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic("failed to mmap: "+err.Error())
	}
}

func (mf *MemmapFile) MapInternal(fr memmap.FileRange, at hostarch.AccessType) (safemem.BlockSeq, error) {
	log.Debugf("MapInternal: ---> Init()")
	defer log.Debugf("MapInternal: <--- Init()")
	bytes := mf.buffer[fr.Start:fr.Start+fr.Length()]
	block := safemem.BlockFromSafeSlice(bytes)
	log.Debugf("Pointer is %p",block.Addr())
	return safemem.BlockSeqOf(block), nil
}

func (mf *MemmapFile) DataFD(fr memmap.FileRange) (int, error) {
	return mf.FD(), nil
}

func (mf *MemmapFile) FD() int {
	return mf.hostfd
}

func (mf *MemmapFile) IncRef(fr memmap.FileRange, memCgID uint32) {}

func (mf *MemmapFile) DecRef(fr memmap.FileRange) {}
