package fakehostfs

import (
	"encoding/binary"
	"os"
	"path"

	"golang.org/x/sys/unix"
	"gvisor.dev/gvisor/pkg/log"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
)

// create new file entry and link to an inode (inode creation can be postponed)
// will not touch inode reference count
func (nativeFS *nativeFilesystem) RegisterNode(basePath string, parentName string, name string, ino uint64, root bool) error {
	log.Debugf("RegisterNode %s",name)
	newMetadataPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"i"+name)
	if root {
		newMetadataPath = path.Join(nativeFS.entriesPath,basePath,"i"+name)
	}
	file, err := os.OpenFile(newMetadataPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	defer file.Close()
	if err != nil {
		log.Debugf("Failed to open %s",newMetadataPath)
		return linuxerr.EINVAL
	}
	buffer := make([]byte,4)
	_, err = binary.Encode(buffer,binary.NativeEndian,uint32(ino))
	if err != nil {
		log.Debugf("Failed to encode: %s",err.Error())
		return linuxerr.EINVAL
	}
	_, err = file.Write(buffer)
	if err != nil {
		log.Debugf("Failed to write file: %s",err.Error())
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) RegisterFile(basePath string, parentName string, name string, ino uint64, root bool) error {
	return nativeFS.RegisterNode(basePath,parentName,name,ino,root)
}

func (nativeFS *nativeFilesystem) RegisterDirectory(basePath string, parentName string, name string, ino uint64, root bool) error {
	newEntryPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"x"+name)
	if root {
		newEntryPath = path.Join(nativeFS.entriesPath,basePath,"x"+name)
	}
	err := unix.Mkdir(newEntryPath, uint32(S_IFDIR|0o700))
	if err != nil {
		log.Debugf("Failed to create directory entry")
		return linuxerr.EINVAL
	}
	return nativeFS.RegisterNode(basePath,parentName,name,ino,root)
}

func (nativeFS *nativeFilesystem) DeleteNode(basePath string, parentName string, name string, root bool) error {
	log.Debugf("DeleteNode %s",name)
	newMetadataPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"i"+name)
	if root {
		newMetadataPath = path.Join(nativeFS.entriesPath,basePath,"i"+name)
	}
	err := os.Remove(newMetadataPath)
	if err != nil {
		log.Debugf("Failed to delete %s",newMetadataPath)
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) DeleteFile(basePath string, parentName string, name string, root bool) error {
	return nativeFS.DeleteNode(basePath,parentName,name,root)
}

func (nativeFS *nativeFilesystem) DeleteDirectory(basePath string, parentName string, name string, root bool) error {
	newEntryPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"x"+name)
	if root {
		newEntryPath = path.Join(nativeFS.entriesPath,basePath,"x"+name)
	}
	err := os.Remove(newEntryPath)
	if err != nil {
		log.Debugf("Failed to delete directory entry")
		return linuxerr.EINVAL
	}
	return nativeFS.DeleteNode(basePath,parentName,name,root)
}

func (nativeFS *nativeFilesystem) RenameNode(basePath string, parentName string, name string, root bool, dstBasePath string, dstParentName string, newName string, dstRoot bool) error {
	log.Debugf("RenameNode %s",name)
	newMetadataPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"i"+name)
	if root {
		newMetadataPath = path.Join(nativeFS.entriesPath,basePath,"i"+name)
	}
	dstMetadataPath := path.Join(nativeFS.entriesPath,dstBasePath,"x"+dstParentName,"i"+newName)
	if dstRoot {
		dstMetadataPath = path.Join(nativeFS.entriesPath,dstBasePath,"i"+newName)
	}
	err := unix.Rename(newMetadataPath, dstMetadataPath)
	if err != nil {
		log.Debugf("Failed to rename %s to %s",newMetadataPath, dstMetadataPath)
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) RenameFile(basePath string, parentName string, name string, root bool, dstBasePath string, dstParentName string, newName string, dstRoot bool) error {
	return nativeFS.RenameNode(basePath,parentName,name,root,dstBasePath,dstParentName,newName,dstRoot)
}

func (nativeFS *nativeFilesystem) RenameDirectory(basePath string, parentName string, name string, root bool, dstBasePath string, dstParentName string, newName string, dstRoot bool) error {
	log.Debugf("RenameNode %s",name)
	newEntryPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"x"+name)
	if root {
		newEntryPath = path.Join(nativeFS.entriesPath,basePath,"x"+name)
	}
	dstEntryPath := path.Join(nativeFS.entriesPath,dstBasePath,"x"+dstParentName,"x"+newName)
	if dstRoot {
		dstEntryPath = path.Join(nativeFS.entriesPath,dstBasePath,"x"+newName)
	}
	err := unix.Rename(newEntryPath, dstEntryPath)
	if err != nil {
		log.Debugf("Failed to rename %s to %s",newEntryPath, dstEntryPath)
		return linuxerr.EINVAL
	}
	return nativeFS.RenameNode(basePath,parentName,name,root,dstBasePath,dstParentName,newName,dstRoot)
}

func (nativeFS *nativeFilesystem) UpdatePathIno(logicalPath string, ino uint64) error {
	realPath := path.Join(nativeFS.entriesPath,logicalPath)
	file, err := os.OpenFile(realPath, os.O_WRONLY|os.O_TRUNC, 0)
	defer file.Close()
	if err != nil {
		log.Debugf("Failed to open %s",realPath)
		return linuxerr.EINVAL
	}
	buffer := make([]byte,4)
	_, err = binary.Encode(buffer,binary.NativeEndian,uint32(ino))
	if err != nil {
		log.Debugf("Failed to encode: %s",err.Error())
		return linuxerr.EINVAL
	}
	_, err = file.Write(buffer)
	if err != nil {
		log.Debugf("Failed to write file: %s",err.Error())
		return linuxerr.EINVAL
	}
	return nil
}

func (nativeFS *nativeFilesystem) GetIno(basePath string, parentName string, name string, root bool) (uint64, error) {
	log.Debugf("GetIno %s",name)
	newMetadataPath := path.Join(nativeFS.entriesPath,basePath,"x"+parentName,"i"+name)
	if root {
		newMetadataPath = path.Join(nativeFS.entriesPath,basePath,"i"+name)
	}
	buffer, err := os.ReadFile(newMetadataPath)
	if err != nil {
		// assume that no such file
		log.Debugf("Failed to read %s: %s",newMetadataPath,err.Error())
		return 0, linuxerr.ENOENT
	}
	var ino uint32
	_, err = binary.Decode(buffer,binary.NativeEndian,&ino)
	if err != nil {
		log.Debugf("Failed to decode: %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	return uint64(ino), nil
}

func (nativeFS *nativeFilesystem) GetInoFromPath(logicalPath string) (uint64, error) {
	realPath := path.Join(nativeFS.entriesPath,logicalPath)
	buffer, err := os.ReadFile(realPath)
	if err != nil {
		// assume that no such file
		log.Debugf("Failed to read %s: %s",realPath,err.Error())
		return 0, linuxerr.ENOENT
	}
	var ino uint32
	_, err = binary.Decode(buffer,binary.NativeEndian,&ino)
	if err != nil {
		log.Debugf("Failed to decode: %s",err.Error())
		return 0, linuxerr.EINVAL
	}
	return uint64(ino), nil
}
