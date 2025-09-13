package graphics

import (
	//"io"
	//"unsafe"
	//"reflect"
	//"runtime"

	//"gvisor.dev/gvisor/pkg/abi/linux"
	//"gvisor.dev/gvisor/pkg/sentry/vfs"
	//"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
	//"gvisor.dev/gvisor/pkg/gohacks"
	//"gvisor.dev/gvisor/pkg/usermem"
	//"gvisor.dev/gvisor/pkg/sentry/arch"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/errors/linuxerr"
	"gvisor.dev/gvisor/pkg/hostarch"
)

func DrmIoctlVersionHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_VERSION ioctl")
	ver := &DrmVersion{}
	bytesCopied, err := ver.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)
	ver.version_major = state.versionMajor
	ver.version_minor = 57
	ver.version_patchlevel = 114514
	ver.name_len = uint64(len(driverName))
	ver.date_len = uint64(len(driverDate))
	ver.desc_len = uint64(len(driverDesc))

	CopyString := func (t *kernel.Task, pointer uint64, srcString []byte) {
		if pointer == 0 {
			return
		}
		length, err := t.CopyOutBytes(hostarch.Addr(pointer), srcString)
		if err != nil {
			log.Debugf("Failed to copy to string")
		}
		log.Debugf("Copied %d bytes to string", length)
	}

	CopyString(t, ver.name, driverName)
	CopyString(t, ver.date, driverDate)
	CopyString(t, ver.desc, driverDesc)

	bytesCopied, err = ver.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)
	return 0, nil
}

func DrmIoctlModeGetResourcesHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_GETRESOURCES ioctl")
	res := &DrmModeCardRes{}
	bytesCopied, err := res.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	res.min_width = minWidth
	res.min_height = minHeight
	res.max_width = maxWidth
	res.max_height = maxHeight

	connectorIds := state.GetConnectorIds()
	SetUint32ArrayByPointer(t, res.connector_id_ptr, connectorIds[:res.count_connectors])
	res.count_connectors = uint32(len(connectorIds))

	encoderIds := state.GetEncoderIds()
	SetUint32ArrayByPointer(t, res.encoder_id_ptr, encoderIds[:res.count_encoders])
	res.count_encoders = uint32(len(encoderIds))

	crtcIds := state.GetCRTCids()
	SetUint32ArrayByPointer(t, res.crtc_id_ptr, crtcIds[:res.count_crtcs])
	res.count_crtcs = uint32(len(crtcIds))

	res.count_fbs = 0

	bytesCopied, err = res.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeGetConnectorHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_GETCONNECTOR ioctl")
	conn := &DrmModeGetConnector{}
	bytesCopied, err := conn.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	connector, ok := state.connectors[conn.connector_id]
	if !ok {
		// no such connector
		return 0, linuxerr.EINVAL
	}

	SetModeinfoArrayByPointer(t, conn.modes_ptr, connector.modes[:conn.count_modes])
	conn.count_modes = uint32(len(connector.modes))

	SetUint32ArrayByPointer(t, conn.encoders_ptr, connector.encoderIds[:conn.count_encoders])
	conn.count_encoders = uint32(len(connector.encoderIds))

	conn.encoder_id = 0 // TODO use current active encoder

	conn.connector_type = connector.connectorType
	conn.connector_type_id = 1 // represents connector id if there are multiple connectors of same type. don't care
	conn.connection = DRM_MODE_CONNECTED // always connected

	conn.subpixel = DRM_MODE_SUBPIXEL_HORIZONTAL_RGB

	conn.mm_width = 1280
	conn.mm_height = 720

	conn.count_props = 0 //TODO

	bytesCopied, err = conn.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmModeCreateDumbHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_CREATE_DUMB ioctl")
	target := &DrmModeCreateDumb{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	log.Debugf("%dx%d, bpp %d, flags %d", target.width, target.height, target.bpp, target.flags)

	target.handle = 114514
	target.pitch = 0
	target.size = 1280*720*4

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeAddFbHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_ADDFB ioctl")
	target := &DrmModeFbCmd{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	log.Debugf("id: %d, %dx%d, pitch: %d, depth: %d, handle: %d", target.fb_id, target.width, target.height, target.pitch, target.depth, target.handle)

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeAddFb2Handler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_ADDFB2 ioctl")
	target := &DrmModeFbCmd2{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	// TODO replace by real impl
	log.Debugf("id: %d, %dx%d, pitch: %d, pixel format: %d, flags: %d", target.fb_id, target.width, target.height, target.pixel_format, target.flags)

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeGetEncoderHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_GETENCODER ioctl")
	target := &DrmModeGetEncoder{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	encoder, ok := state.encoders[target.encoder_id]
	if !ok {
		// no such encoder
		return 0, linuxerr.EINVAL
	}
	target.encoder_type = encoder.encoderType

	target.crtc_id = 0
	target.possible_crtcs = 1 // 1<<0 = 1

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeMapDumbHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_MAP_DUMB ioctl")
	target := &DrmModeMapDumb{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	log.Debugf("Handle is %d", target.handle)
	target.offset = 0

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlSetMasterHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_SET_MASTER ioctl")
	if data != 0 {
		return 0, linuxerr.EINVAL
	}
	return 0, nil
}

func DrmIoctlModeGetCRTC_Handler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_GETCRTC ioctl")
	target := &DrmModeCRTC{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	crtc, ok := state.crtcs[target.crtc_id]
	if !ok {
		// no such CRTC
		return 0, linuxerr.EINVAL
	}

	SetUint32ArrayByPointer(t, target.set_connectors_ptr, crtc.connectorIds[:target.count_connectors])
	target.count_connectors = uint32(len(crtc.connectorIds))

	target.fb_id = 0
	target.x = 0
	target.y = 0
	target.gamma_size = 4096
	target.mode_valid = 1

	mode := state.connectors[0].modes[0]
	modeinfo := DrmModeModeinfo {
		hdisplay: mode.height,
		vdisplay: mode.width,
		vrefresh: mode.refreshRate,
	}
	target.drm_mode_modeinfo = modeinfo

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeSetCRTC_Handler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_SETCRTC ioctl")
	target := &DrmModeCRTC{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	// TODO replace by real impl
	_, ok := state.crtcs[target.crtc_id]
	if !ok {
		// no such CRTC
		return 0, linuxerr.EINVAL
	}

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlModeCRTC_PageFlipHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_MODE_PAGE_FLIP ioctl")
	target := &DrmModeCRTC_PageFlip{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	// TODO replace by real impl
	_, ok := state.crtcs[target.crtc_id]
	if !ok {
		// no such CRTC
		return 0, linuxerr.EINVAL
	}

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

func DrmIoctlGetCapHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got DRM_IOCTL_GET_CAP ioctl")
	target := &DrmGetCap{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	if target.capability == DRM_CAP_ASYNC_PAGE_FLIP {
		target.value = 1
	}

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}

/*
func XyzHandler(state *CardDeviceState, t *kernel.Task, data hostarch.Addr) (uintptr, error) {
	log.Debugf("Got XYZ ioctl")
	target := &Xyz{}
	bytesCopied, err := target.CopyIn(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	target.x = y
	target.z = w

	bytesCopied, err = target.CopyOut(t, data)
	if err != nil {
		log.Debugf("Failed to map guest struct")
		return 0, linuxerr.EINVAL
	}
	log.Debugf("%d bytes copied",bytesCopied)

	return 0, nil
}
*/
