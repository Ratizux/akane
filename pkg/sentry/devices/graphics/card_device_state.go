package graphics

type CardDeviceState struct {
	// mappingMu sync.Mutex
	mf *MemmapFile

	connectors map[uint32]Connector
	encoders map[uint32]Encoder
	crtcs map[uint32]CRTC

	versionMajor int32
}

func (state *CardDeviceState) GetConnectorIds() []uint32 {
	ids := []uint32{}
	for key, _ := range state.connectors {
		ids = append(ids, key)
	}
	return ids
}

func (state *CardDeviceState) GetEncoderIds() []uint32 {
	ids := []uint32{}
	for key, _ := range state.encoders {
		ids = append(ids, key)
	}
	return ids
}

func (state *CardDeviceState) GetCRTCids() []uint32 {
	ids := []uint32{}
	for key, _ := range state.crtcs {
		ids = append(ids, key)
	}
	return ids
}

type Connector struct {
	modes []Mode
	encoderIds []uint32
	connectorType uint32
}

type Mode struct {
	width uint16
	height uint16
	refreshRate uint32
	name string
}

type Encoder struct {
	encoderType uint32
}

type CRTC struct {
	connectorIds []uint32
}
