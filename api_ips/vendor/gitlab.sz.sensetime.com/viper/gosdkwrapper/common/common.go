package common

import (
	"math"
	"unsafe"
)

type Rect struct {
	Left   float32
	Top    float32
	Right  float32
	Bottom float32
}

func (r Rect) Width() float32 {
	return r.Right - r.Left
}

func (r Rect) Height() float32 {
	return r.Bottom - r.Top
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Size struct {
	Width  int32 `json:"width"`
	Height int32 `json:"height"`
}

type Vector2F Point

func (v *Vector2F) Normalize() Vector2F {
	len := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	return Vector2F{
		X: v.X / len,
		Y: v.Y / len,
	}
}

func (v *Vector2F) NormalizeInplace() {
	len := float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
	v.X = v.X / len
	v.Y = v.Y / len
}

type DeviceType int

const (
	DeviceCPU     DeviceType = 0
	DeviceGPU     DeviceType = 1
	DeviceUnknown DeviceType = 2
)

func (d DeviceType) String() string {
	switch d {
	case DeviceCPU:
		return "cpu"
	case DeviceGPU:
		return "gpu"
	default:
		return "unknown"
	}
}

type PixelFormat int

// nolint: golint
const (
	PIXEL_FORMAT_NONE    PixelFormat = iota
	PIXEL_FORMAT_GRAY                = iota
	PIXEL_FORMAT_BGR                 = iota
	PIXEL_FORMAT_BGRA                = iota
	PIXEL_FORMAT_RGB                 = iota
	PIXEL_FORMAT_ARGB                = iota
	PIXEL_FORMAT_YUV420P             = iota
	PIXEL_FORMAT_NV12                = iota
	PIXEL_FORMAT_NV21                = iota
	PIXEL_FORMAT_GRAY32              = iota
)

type OriginInfo struct {
	Rate    float64
	Width   int
	Height  int
	Resized bool
}

type Image interface {
	Data() unsafe.Pointer

	UserData() unsafe.Pointer

	ExtraInfoBuffer() unsafe.Pointer

	// 图片解码之后如果被resize, OriginInfo中将会存储原始图片的信息
	OriginInfo() *OriginInfo

	// Crop 函数返回的Image语义上是独立的
	Crop(x, y, w, h int) (Image, error)
	DeviceType() DeviceType
	Width() int
	Height() int
	Stride() int
	CommonPixelFormat() PixelFormat
	Release()
}

// type FrameType int

/*
const (
	UNKNOWN_FRAME           = C.KESTREL_UNKNOWN_FRAME
	FLUSH_FRAME   FrameType = C.KESTREL_FLUSH_FRAME
	VIDEO_FRAME             = C.KESTREL_VIDEO_FRAME
	AUDIO_FRAME             = C.KESTREL_AUDIO_FRAME
)
*/

type VideoFrame interface {
	Image
	Pts() int64
	Dts() int64
}
