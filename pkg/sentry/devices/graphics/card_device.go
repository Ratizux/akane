package graphics

import (
	"gvisor.dev/gvisor/pkg/sentry/vfs"
	"gvisor.dev/gvisor/pkg/context"
	"gvisor.dev/gvisor/pkg/log"
)

type CardDevice struct {
	id int32
	state *CardDeviceState
}

func (dev *CardDevice) Init() {
	log.Debugf("CardDevice: ---> Init()")
	defer log.Debugf("CardDevice: <--- Init()")
	state := &CardDeviceState{
		mf: &MemmapFile{},
		versionMajor: 3,
	}
	mode1280x720 := Mode{
		width: 1280,
		height: 720,
		refreshRate: 60,
		name: "1280x720",
	}
	state.connectors=map[uint32]Connector{}
	state.connectors[0] = Connector{
		modes: []Mode{
			mode1280x720,
		},
		encoderIds: []uint32{
			0,
		},
		connectorType: DRM_MODE_CONNECTOR_VIRTUAL,
	}
	state.encoders=map[uint32]Encoder{}
	state.encoders[0] = Encoder {
		encoderType: DRM_MODE_ENCODER_VIRTUAL,
	}
	state.crtcs=map[uint32]CRTC{}
	state.crtcs[0] = CRTC {
		connectorIds: []uint32{
			0,
		},
	}
	state.mf.Init(1280*720*4)

	dev.state = state
}

// Open implements vfs.Device.Open.
func (dev *CardDevice) Open(ctx context.Context, mnt *vfs.Mount, vfsd *vfs.Dentry, opts vfs.OpenOptions) (*vfs.FileDescription, error) {
	log.Debugf("%d",dev.id)
	if dev.state == nil {
		log.Debugf("state is nil")
	}
	return NewCardFD(ctx, mnt, vfsd, opts, dev.state)
}
