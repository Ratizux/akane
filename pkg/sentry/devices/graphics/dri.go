package graphics

import (
	//"gvisor.dev/gvisor/pkg/abi/linux"
)


// got from include/uapi/drm/drm.h
var DRM_IOCTL_VERSION = DRM_IOWR(0x00, DrmVersion{})
var DRM_IOCTL_MODE_GETRESOURCES = DRM_IOWR(0xA0, DrmModeCardRes{})
var DRM_IOCTL_MODE_GETCONNECTOR = DRM_IOWR(0xA7, DrmModeGetConnector{})
var DRM_IOCTL_MODE_CREATE_DUMB = DRM_IOWR(0xB2, DrmModeCreateDumb{})
var DRM_IOCTL_MODE_ADDFB = DRM_IOWR(0xAE, DrmModeFbCmd{})
var DRM_IOCTL_MODE_ADDFB2 = DRM_IOWR(0xB8, DrmModeFbCmd2{})
var DRM_IOCTL_MODE_GETENCODER = DRM_IOWR(0xA6, DrmModeGetEncoder{})
var DRM_IOCTL_MODE_MAP_DUMB = DRM_IOWR(0xB3, DrmModeMapDumb{})
var DRM_IOCTL_SET_MASTER = DRM_IO(0x1e)
var DRM_IOCTL_MODE_SETCRTC = DRM_IOWR(0xA2, DrmModeCRTC{})
var DRM_IOCTL_MODE_GETCRTC = DRM_IOWR(0xA1, DrmModeCRTC{})
var DRM_IOCTL_MODE_PAGE_FLIP = DRM_IOWR(0xB0, DrmModeCRTC_PageFlip{})
var DRM_IOCTL_GET_CAP = DRM_IOWR(0x0c, DrmGetCap{})

var driverName = []byte("blk")
var driverDate = []byte("19260817")
var driverDesc = []byte("Blockcity Virtual Display Adapter")

var minWidth uint32 = 0
var minHeight uint32 = 0
var maxWidth uint32 = 4096
var maxHeight uint32 = 4096
