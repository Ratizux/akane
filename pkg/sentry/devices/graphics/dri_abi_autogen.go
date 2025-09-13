// Automatically generated marshal implementation. See tools/go_marshal.

package graphics

import (
    "gvisor.dev/gvisor/pkg/gohacks"
    "gvisor.dev/gvisor/pkg/hostarch"
    "gvisor.dev/gvisor/pkg/marshal"
    "io"
    "reflect"
    "runtime"
    "unsafe"
)

// Marshallable types used by this file.
var _ marshal.Marshallable = (*DrmGetCap)(nil)
var _ marshal.Marshallable = (*DrmModeCRTC)(nil)
var _ marshal.Marshallable = (*DrmModeCRTC_PageFlip)(nil)
var _ marshal.Marshallable = (*DrmModeCardRes)(nil)
var _ marshal.Marshallable = (*DrmModeCreateDumb)(nil)
var _ marshal.Marshallable = (*DrmModeFbCmd)(nil)
var _ marshal.Marshallable = (*DrmModeFbCmd2)(nil)
var _ marshal.Marshallable = (*DrmModeGetConnector)(nil)
var _ marshal.Marshallable = (*DrmModeGetEncoder)(nil)
var _ marshal.Marshallable = (*DrmModeMapDumb)(nil)
var _ marshal.Marshallable = (*DrmModeModeinfo)(nil)
var _ marshal.Marshallable = (*DrmVersion)(nil)

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmGetCap) SizeBytes() int {
    return 16
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmGetCap) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.capability))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.value))
    dst = dst[8:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmGetCap) UnmarshalBytes(src []byte) []byte {
    d.capability = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.value = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmGetCap) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmGetCap) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmGetCap) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmGetCap) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmGetCap) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmGetCap) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmGetCap) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmGetCap) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeCRTC) SizeBytes() int {
    return 36 +
        (*DrmModeModeinfo)(nil).SizeBytes()
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeCRTC) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.set_connectors_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_connectors))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.crtc_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.fb_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.x))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.y))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.gamma_size))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.mode_valid))
    dst = dst[4:]
    dst = d.drm_mode_modeinfo.MarshalUnsafe(dst)
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeCRTC) UnmarshalBytes(src []byte) []byte {
    d.set_connectors_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.count_connectors = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.crtc_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.fb_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.x = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.y = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.gamma_size = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.mode_valid = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    src = d.drm_mode_modeinfo.UnmarshalUnsafe(src)
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeCRTC) Packed() bool {
    return d.drm_mode_modeinfo.Packed()
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeCRTC) MarshalUnsafe(dst []byte) []byte {
    if d.drm_mode_modeinfo.Packed() {
        size := d.SizeBytes()
        gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
        return dst[size:]
    }
    // Type DrmModeCRTC doesn't have a packed layout in memory, fallback to MarshalBytes.
    return d.MarshalBytes(dst)
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeCRTC) UnmarshalUnsafe(src []byte) []byte {
    if d.drm_mode_modeinfo.Packed() {
        size := d.SizeBytes()
        gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
        return src[size:]
    }
    // Type DrmModeCRTC doesn't have a packed layout in memory, fallback to UnmarshalBytes.
    return d.UnmarshalBytes(src)
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeCRTC) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    if !d.drm_mode_modeinfo.Packed() {
        // Type DrmModeCRTC doesn't have a packed layout in memory, fall back to MarshalBytes.
        buf := cc.CopyScratchBuffer(d.SizeBytes()) // escapes: okay.
        d.MarshalBytes(buf) // escapes: fallback.
        return cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    }

    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeCRTC) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeCRTC) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    if !d.drm_mode_modeinfo.Packed() {
        // Type DrmModeCRTC doesn't have a packed layout in memory, fall back to UnmarshalBytes.
        buf := cc.CopyScratchBuffer(d.SizeBytes()) // escapes: okay.
        length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
        // Unmarshal unconditionally. If we had a short copy-in, this results in a
        // partially unmarshalled struct.
        d.UnmarshalBytes(buf) // escapes: fallback.
        return length, err
    }

    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeCRTC) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeCRTC) WriteTo(writer io.Writer) (int64, error) {
    if !d.drm_mode_modeinfo.Packed() {
        // Type DrmModeCRTC doesn't have a packed layout in memory, fall back to MarshalBytes.
        buf := make([]byte, d.SizeBytes())
        d.MarshalBytes(buf)
        length, err := writer.Write(buf)
        return int64(length), err
    }

    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeCRTC_PageFlip) SizeBytes() int {
    return 24
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeCRTC_PageFlip) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.crtc_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.fb_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.flags))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.reserved))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.user_data))
    dst = dst[8:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeCRTC_PageFlip) UnmarshalBytes(src []byte) []byte {
    d.crtc_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.fb_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.flags = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.reserved = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.user_data = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeCRTC_PageFlip) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeCRTC_PageFlip) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeCRTC_PageFlip) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeCRTC_PageFlip) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeCRTC_PageFlip) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeCRTC_PageFlip) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeCRTC_PageFlip) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeCRTC_PageFlip) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeCardRes) SizeBytes() int {
    return 64
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeCardRes) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.fb_id_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.crtc_id_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.connector_id_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.encoder_id_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_fbs))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_crtcs))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_connectors))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_encoders))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.min_width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.max_width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.min_height))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.max_height))
    dst = dst[4:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeCardRes) UnmarshalBytes(src []byte) []byte {
    d.fb_id_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.crtc_id_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.connector_id_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.encoder_id_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.count_fbs = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.count_crtcs = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.count_connectors = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.count_encoders = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.min_width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.max_width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.min_height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.max_height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeCardRes) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeCardRes) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeCardRes) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeCardRes) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeCardRes) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeCardRes) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeCardRes) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeCardRes) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeCreateDumb) SizeBytes() int {
    return 32
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeCreateDumb) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.height))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.bpp))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.flags))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.handle))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pitch))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.size))
    dst = dst[8:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeCreateDumb) UnmarshalBytes(src []byte) []byte {
    d.height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.bpp = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.flags = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.handle = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.pitch = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.size = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeCreateDumb) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeCreateDumb) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeCreateDumb) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeCreateDumb) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeCreateDumb) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeCreateDumb) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeCreateDumb) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeCreateDumb) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeFbCmd) SizeBytes() int {
    return 28
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeFbCmd) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.fb_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.height))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pitch))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.bpp))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.depth))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.handle))
    dst = dst[4:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeFbCmd) UnmarshalBytes(src []byte) []byte {
    d.fb_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.pitch = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.bpp = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.depth = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.handle = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeFbCmd) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeFbCmd) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeFbCmd) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeFbCmd) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeFbCmd) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeFbCmd) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeFbCmd) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeFbCmd) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeFbCmd2) SizeBytes() int {
    return 20 +
        4*4 +
        4*4 +
        4*4 +
        8*4
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeFbCmd2) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.fb_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.height))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pixel_format))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.flags))
    dst = dst[4:]
    for idx := 0; idx < 4; idx++ {
        hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.handles[idx]))
        dst = dst[4:]
    }
    for idx := 0; idx < 4; idx++ {
        hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pitches[idx]))
        dst = dst[4:]
    }
    for idx := 0; idx < 4; idx++ {
        hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.offsets[idx]))
        dst = dst[4:]
    }
    for idx := 0; idx < 4; idx++ {
        hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.modifier[idx]))
        dst = dst[8:]
    }
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeFbCmd2) UnmarshalBytes(src []byte) []byte {
    d.fb_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.pixel_format = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.flags = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    for idx := 0; idx < 4; idx++ {
        d.handles[idx] = uint32(hostarch.ByteOrder.Uint32(src[:4]))
        src = src[4:]
    }
    for idx := 0; idx < 4; idx++ {
        d.pitches[idx] = uint32(hostarch.ByteOrder.Uint32(src[:4]))
        src = src[4:]
    }
    for idx := 0; idx < 4; idx++ {
        d.offsets[idx] = uint32(hostarch.ByteOrder.Uint32(src[:4]))
        src = src[4:]
    }
    for idx := 0; idx < 4; idx++ {
        d.modifier[idx] = uint64(hostarch.ByteOrder.Uint64(src[:8]))
        src = src[8:]
    }
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeFbCmd2) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeFbCmd2) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeFbCmd2) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeFbCmd2) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeFbCmd2) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeFbCmd2) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeFbCmd2) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeFbCmd2) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeGetConnector) SizeBytes() int {
    return 80
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeGetConnector) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.encoders_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.modes_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.props_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.prop_values_ptr))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_modes))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_props))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.count_encoders))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.encoder_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.connector_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.connector_type))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.connector_type_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.connection))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.mm_width))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.mm_height))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.subpixel))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pad))
    dst = dst[4:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeGetConnector) UnmarshalBytes(src []byte) []byte {
    d.encoders_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.modes_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.props_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.prop_values_ptr = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.count_modes = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.count_props = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.count_encoders = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.encoder_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.connector_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.connector_type = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.connector_type_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.connection = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.mm_width = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.mm_height = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.subpixel = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.pad = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeGetConnector) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeGetConnector) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeGetConnector) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeGetConnector) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeGetConnector) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeGetConnector) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeGetConnector) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeGetConnector) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeGetEncoder) SizeBytes() int {
    return 20
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeGetEncoder) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.encoder_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.encoder_type))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.crtc_id))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.possible_crtcs))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.possible_clones))
    dst = dst[4:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeGetEncoder) UnmarshalBytes(src []byte) []byte {
    d.encoder_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.encoder_type = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.crtc_id = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.possible_crtcs = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.possible_clones = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeGetEncoder) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeGetEncoder) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeGetEncoder) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeGetEncoder) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeGetEncoder) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeGetEncoder) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeGetEncoder) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeGetEncoder) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeMapDumb) SizeBytes() int {
    return 16
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeMapDumb) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.handle))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.pad))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.offset))
    dst = dst[8:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeMapDumb) UnmarshalBytes(src []byte) []byte {
    d.handle = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.pad = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.offset = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeMapDumb) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeMapDumb) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeMapDumb) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeMapDumb) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeMapDumb) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeMapDumb) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeMapDumb) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeMapDumb) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmModeModeinfo) SizeBytes() int {
    return 36 +
        1*DRM_DISPLAY_MODE_LEN
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmModeModeinfo) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.clock))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.hdisplay))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.hsync_start))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.hsync_end))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.htotal))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.hskew))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.vdisplay))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.vsync_start))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.vsync_end))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.vtotal))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint16(dst[:2], uint16(d.vscan))
    dst = dst[2:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.vrefresh))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.flags))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.mode_type))
    dst = dst[4:]
    for idx := 0; idx < DRM_DISPLAY_MODE_LEN; idx++ {
        dst[0] = byte(d.name[idx])
        dst = dst[1:]
    }
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmModeModeinfo) UnmarshalBytes(src []byte) []byte {
    d.clock = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.hdisplay = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.hsync_start = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.hsync_end = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.htotal = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.hskew = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vdisplay = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vsync_start = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vsync_end = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vtotal = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vscan = uint16(hostarch.ByteOrder.Uint16(src[:2]))
    src = src[2:]
    d.vrefresh = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.flags = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.mode_type = uint32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    for idx := 0; idx < DRM_DISPLAY_MODE_LEN; idx++ {
        d.name[idx] = src[0]
        src = src[1:]
    }
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmModeModeinfo) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmModeModeinfo) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmModeModeinfo) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmModeModeinfo) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmModeModeinfo) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmModeModeinfo) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmModeModeinfo) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmModeModeinfo) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

// SizeBytes implements marshal.Marshallable.SizeBytes.
func (d *DrmVersion) SizeBytes() int {
    return 60
}

// MarshalBytes implements marshal.Marshallable.MarshalBytes.
func (d *DrmVersion) MarshalBytes(dst []byte) []byte {
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.version_major))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.version_minor))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint32(dst[:4], uint32(d.version_patchlevel))
    dst = dst[4:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.name_len))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.name))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.date_len))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.date))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.desc_len))
    dst = dst[8:]
    hostarch.ByteOrder.PutUint64(dst[:8], uint64(d.desc))
    dst = dst[8:]
    return dst
}

// UnmarshalBytes implements marshal.Marshallable.UnmarshalBytes.
func (d *DrmVersion) UnmarshalBytes(src []byte) []byte {
    d.version_major = int32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.version_minor = int32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.version_patchlevel = int32(hostarch.ByteOrder.Uint32(src[:4]))
    src = src[4:]
    d.name_len = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.name = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.date_len = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.date = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.desc_len = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    d.desc = uint64(hostarch.ByteOrder.Uint64(src[:8]))
    src = src[8:]
    return src
}

// Packed implements marshal.Marshallable.Packed.
//go:nosplit
func (d *DrmVersion) Packed() bool {
    return true
}

// MarshalUnsafe implements marshal.Marshallable.MarshalUnsafe.
func (d *DrmVersion) MarshalUnsafe(dst []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(&dst[0]), unsafe.Pointer(d), uintptr(size))
    return dst[size:]
}

// UnmarshalUnsafe implements marshal.Marshallable.UnmarshalUnsafe.
func (d *DrmVersion) UnmarshalUnsafe(src []byte) []byte {
    size := d.SizeBytes()
    gohacks.Memmove(unsafe.Pointer(d), unsafe.Pointer(&src[0]), uintptr(size))
    return src[size:]
}

// CopyOutN implements marshal.Marshallable.CopyOutN.
func (d *DrmVersion) CopyOutN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyOutBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyOut implements marshal.Marshallable.CopyOut.
func (d *DrmVersion) CopyOut(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyOutN(cc, addr, d.SizeBytes())
}

// CopyInN implements marshal.Marshallable.CopyInN.
func (d *DrmVersion) CopyInN(cc marshal.CopyContext, addr hostarch.Addr, limit int) (int, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := cc.CopyInBytes(addr, buf[:limit]) // escapes: okay.
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return length, err
}

// CopyIn implements marshal.Marshallable.CopyIn.
func (d *DrmVersion) CopyIn(cc marshal.CopyContext, addr hostarch.Addr) (int, error) {
    return d.CopyInN(cc, addr, d.SizeBytes())
}

// WriteTo implements io.WriterTo.WriteTo.
func (d *DrmVersion) WriteTo(writer io.Writer) (int64, error) {
    // Construct a slice backed by dst's underlying memory.
    var buf []byte
    hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
    hdr.Data = uintptr(gohacks.Noescape(unsafe.Pointer(d)))
    hdr.Len = d.SizeBytes()
    hdr.Cap = d.SizeBytes()

    length, err := writer.Write(buf)
    // Since we bypassed the compiler's escape analysis, indicate that d
    // must live until the use above.
    runtime.KeepAlive(d) // escapes: replaced by intrinsic.
    return int64(length), err
}

