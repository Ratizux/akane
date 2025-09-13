package graphics

func fourccCode(a, b, c, d uint32) uint32 {
	return (a)|(b<<8)|(c<<16)|(d<<24)
}

var DRM_FORMAT_BGRA8888 = fourccCode('B', 'A', '2', '4')
