package fakehostfs

import (
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/hostarch"
	"gvisor.dev/gvisor/pkg/sentry/memmap"
)

// AddMapping implements memmap.Mappable.AddMapping.
func (fd *FakehostfsFileDescription) AddMapping(ctx context.Context, ms memmap.MappingSpace, ar hostarch.AddrRange, offset uint64, writable bool) error {
	log.Debugf("fakehostfs: ---> AddMapping(): %d",fd.inode.Ino())
	fd.mappingCount++
	log.Debugf("Mapping Count %d",fd.mappingCount)
	defer log.Debugf("fakehostfs: <--- AddMapping(): %d",fd.inode.Ino())
	return fd.CachedMappable.AddMapping(ctx, ms, ar, offset, writable)
}

// RemoveMapping implements memmap.Mappable.RemoveMapping.
func (fd *FakehostfsFileDescription) RemoveMapping(ctx context.Context, ms memmap.MappingSpace, ar hostarch.AddrRange, offset uint64, writable bool) {
	log.Debugf("fakehostfs: ---> RemoveMapping(): %d",fd.inode.Ino())
	fd.mappingCount--
	log.Debugf("Mapping Count %d",fd.mappingCount)
	defer log.Debugf("fakehostfs: <--- RemoveMapping(): %d",fd.inode.Ino())
	fd.CachedMappable.RemoveMapping(ctx, ms, ar, offset, writable)
}

// CopyMapping implements memmap.Mappable.CopyMapping.
func (fd *FakehostfsFileDescription) CopyMapping(ctx context.Context, ms memmap.MappingSpace, srcAR, dstAR hostarch.AddrRange, offset uint64, writable bool) error {
	log.Debugf("fakehostfs: ---> CopyMapping(): %d",fd.inode.Ino())
	defer log.Debugf("fakehostfs: <--- CopyMapping(): %d",fd.inode.Ino())
	return fd.CachedMappable.CopyMapping(ctx, ms, srcAR, dstAR, offset, writable)
}

// Translate implements memmap.Mappable.Translate.
func (fd *FakehostfsFileDescription) Translate(ctx context.Context, required, optional memmap.MappableRange, at hostarch.AccessType) ([]memmap.Translation, error) {
	log.Debugf("fakehostfs: ---> Translate(): %d",fd.inode.Ino())
	defer log.Debugf("fakehostfs: <--- Translate(): %d",fd.inode.Ino())
	return fd.CachedMappable.Translate(ctx, required, optional, at)
}

// InvalidateUnsavable implements memmap.Mappable.InvalidateUnsavable.
func (fd *FakehostfsFileDescription) InvalidateUnsavable(ctx context.Context) error {
	log.Debugf("fakehostfs: ---> InvalidateUnsavable(): %d",fd.inode.Ino())
	defer log.Debugf("fakehostfs: <--- InvalidateUnsavable(): %d",fd.inode.Ino())
	return fd.CachedMappable.InvalidateUnsavable(ctx)
}

// InvalidateRange invalidates the passed range on i.mappings.
func (fd *FakehostfsFileDescription) InvalidateRange(r memmap.MappableRange) {
	log.Debugf("fakehostfs: ---> InvalidateRange(): %d",fd.inode.Ino())
	defer log.Debugf("fakehostfs: <--- InvalidateRange(): %d",fd.inode.Ino())
	fd.CachedMappable.InvalidateRange(r)
}

// InitFileMapperOnce initializes the host file mapper. It ensures that the
// file mapper is initialized just once.
func (fd *FakehostfsFileDescription) InitFileMapperOnce() {
	log.Debugf("fakehostfs: ---> InitFileMapperOnce(): %d",fd.inode.Ino())
	defer log.Debugf("fakehostfs: <--- InitFileMapperOnce(): %d",fd.inode.Ino())
	fd.CachedMappable.InitFileMapperOnce()
}
