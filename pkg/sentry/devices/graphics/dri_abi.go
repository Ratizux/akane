package graphics

// please call go_marshal command to generate corresponding CopyIn/CopyOut methods of structs

// +marshal
type DrmVersion struct {
	version_major int32
	version_minor int32
	version_patchlevel int32
	name_len uint64
	name uint64
	date_len uint64
	date uint64
	desc_len uint64
	desc uint64
}

// +marshal
type DrmModeCardRes struct {
	fb_id_ptr uint64
	crtc_id_ptr uint64
	connector_id_ptr uint64
	encoder_id_ptr uint64
	count_fbs uint32
	count_crtcs uint32
	count_connectors uint32
	count_encoders uint32
	min_width uint32
	max_width uint32
	min_height uint32
	max_height uint32
}

// +marshal
type DrmModeGetConnector struct {
	encoders_ptr uint64
	modes_ptr uint64
	props_ptr uint64
	prop_values_ptr uint64
	count_modes uint32 // mode: choice of resolution+refresh rate
	count_props uint32
	count_encoders uint32
	encoder_id uint32
	connector_id uint32
	connector_type uint32
	connector_type_id uint32
	connection uint32
	mm_width uint32
	mm_height uint32
	subpixel uint32
	pad uint32
}

const (
	DRM_DISPLAY_MODE_LEN = 32
)

// +marshal
type DrmModeModeinfo struct {
	clock uint32
	hdisplay uint16 // height in pixel
	hsync_start uint16
	hsync_end uint16
	htotal uint16
	hskew uint16
	vdisplay uint16 // width in pixel
	vsync_start uint16
	vsync_end uint16
	vtotal uint16
	vscan uint16
	vrefresh uint32 // refresh rate
	flags uint32
	mode_type uint32 // mode_ prefix to avoid conflict with go keyword
	name [DRM_DISPLAY_MODE_LEN]byte
};

// +marshal
type DrmModeCreateDumb struct {
	height uint32
	width uint32
	bpp uint32
	flags uint32
	handle uint32 // will be used in DrmModeFbCmd
	pitch uint32
	size uint64
}

// +marshal
type DrmModeFbCmd struct {
	fb_id uint32
	width uint32
	height uint32
	pitch uint32
	bpp uint32
	depth uint32
	handle uint32
}

// +marshal
type DrmModeGetEncoder struct {
	encoder_id uint32
	encoder_type uint32
	crtc_id uint32
	possible_crtcs uint32
	possible_clones uint32
}

// +marshal
type DrmModeMapDumb struct {
	handle uint32 // to the dumb buffer
	pad uint32 // useless
	offset uint64 // offset for mmap()
}

// +marshal
type DrmModeCRTC struct {
	set_connectors_ptr uint64 // array of connector ids in uint32
	count_connectors uint32
	crtc_id uint32
	fb_id uint32
	x uint32
	y uint32
	gamma_size uint32
	mode_valid uint32
	drm_mode_modeinfo DrmModeModeinfo
}

// +marshal
type DrmModeFbCmd2 struct {
	 fb_id uint32
	 width uint32
	 height uint32
	 pixel_format uint32
	 flags uint32
	 handles [4]uint32
	 pitches [4]uint32
	 offsets [4]uint32
	 modifier [4]uint64
}

// +marshal
type DrmModeCRTC_PageFlip struct {
	crtc_id uint32
	fb_id uint32
	flags uint32
	reserved uint32
	user_data uint64
}

// +marshal
type DrmGetCap struct {
	capability uint64
	value uint64
}
