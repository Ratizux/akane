package graphics

import (
	"io"

	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/sentry/memmap"
)

func (fd *CardFD) ConfigureMMap(ctx context.Context, opts *memmap.MMapOpts) error {
	log.Debugf("ConfigureMMap of cardFD is called")
	opts.SentryOwnedContent = true

	return vfs.GenericProxyDeviceConfigureMMap(&fd.vfsfd, fd, opts)
}

func (fd *CardFD) AddMapping(ctx context.Context, ms memmap.MappingSpace, ar hostarch.AddrRange, offset uint64, writable bool) error {
	fd.mappings.AddMapping(ms, ar, offset, writable)
	return nil
}

func (fd *CardFD) RemoveMapping(ctx context.Context, ms memmap.MappingSpace, ar hostarch.AddrRange, offset uint64, writable bool) {
	fd.mappings.RemoveMapping(ms, ar, offset, writable)
}

func (fd *CardFD) CopyMapping(ctx context.Context, ms memmap.MappingSpace, srcAR, dstAR hostarch.AddrRange, offset uint64, writable bool) error {
	return fd.AddMapping(ctx, ms, dstAR, offset, writable)
}

func (fd *CardFD) InvalidateUnsavable(ctx context.Context) error {
	fd.mappings.InvalidateAll(memmap.InvalidateOpts{})
	return nil
}

func (fd *CardFD) Size() uint64 {
	return uint64(len(fd.state.mf.buffer))
}

func (fd *CardFD) Translate(ctx context.Context, required, optional memmap.MappableRange, at hostarch.AccessType) ([]memmap.Translation, error) {
	log.Debugf("cardFD: ---> Translate()")
	defer log.Debugf("cardFD: <--- Translate()")
	if fd.state.mf == nil {
		return nil, &memmap.BusError{io.EOF}
	}

	pgend, _ := hostarch.PageRoundUp(fd.Size())
	if required.End > pgend {
		if required.Start >= pgend {
			return nil, &memmap.BusError{io.EOF}
		}
		required.End = pgend
	}
	if optional.End > pgend {
		optional.End = pgend
	}
	mr := optional
	log.Debugf("required: %d-%d", required.Start, required.End)
	log.Debugf("optional: %d-%d", optional.Start, optional.End)
	return []memmap.Translation{
		{
			Source: mr,
			File:   fd.state.mf,
			Offset: mr.Start,
			Perms:  hostarch.ReadWrite,
		},
	}, nil
}
