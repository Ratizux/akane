package fakehostfs

import (
	"encoding/binary"
	"os"
	//"path"

	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/log"
)

type InodeMetadata struct {
	Mode uint16
	ReferenceCount uint16
	UID uint32
	GID uint32
	CTime int64
	MTime int64
}

func (metadata *InodeMetadata) Marshal() ([]byte, error) {
	buffer := make([]byte,28)
	_, err := binary.Encode(buffer,binary.NativeEndian,metadata)
	if err != nil {
		return buffer, linuxerr.EINVAL
	}
	return buffer, nil
}

func (metadata *InodeMetadata) Unmarshal(buffer []byte) error {
	_, err := binary.Decode(buffer,binary.NativeEndian,metadata)
	if err != nil {
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) SetInoMetadata(ino uint64, inodeMetadata InodeMetadata) (error) {
	_, metadataPath, err := nativeFS.GetInodePaths(ino)
	if err != nil {
		log.Debugf("Error getting inode path")
		return linuxerr.EINVAL
	}
	file, err := os.OpenFile(metadataPath, os.O_WRONLY|os.O_TRUNC, 0)
	defer file.Close()
	if err != nil {
		log.Debugf("Error opening metadata file")
		return linuxerr.EINVAL
	}
	buffer, err := inodeMetadata.Marshal()
	if err != nil {
		log.Debugf("Error serializing metadata")
		return linuxerr.EINVAL
	}
	_, err = file.Write(buffer)
	if err != nil {
		log.Debugf("Error writing metadata to file")
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) GetInoMetadata(ino uint64) (InodeMetadata, error) {
	_, metadataPath, err := nativeFS.GetInodePaths(ino)
	buffer, err := os.ReadFile(metadataPath)
	if err != nil {
		return InodeMetadata{}, linuxerr.EINVAL
	}
	inodeMetadata := InodeMetadata{}
	err = inodeMetadata.Unmarshal(buffer)
	if err != nil {
		return InodeMetadata{}, linuxerr.EINVAL
	}
	return inodeMetadata, nil
}
