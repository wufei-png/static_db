package cuda

/*
#cgo LDFLAGS: -ldl

#include "decoder.h"

CUresult wrapCuvidMapVideoFrame64(CUvideodecoder hDecoder, int nPicIdx, CUdeviceptr *pDevPtr, unsigned int *pPitch, CUVIDPROCPARAMS *pVPP) {
	return cuvidMapVideoFrame64(hDecoder, nPicIdx, pDevPtr, pPitch, pVPP);
}

CUresult wrapCuvidUnmapVideoFrame64(CUvideodecoder hDecoder, CUdeviceptr DevPtr) {
	return cuvidUnmapVideoFrame64(hDecoder, DevPtr);
}
*/
import "C"
import (
	"encoding/binary"
	"unsafe"

	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/codec"
	"github.com/nareix/joy4/codec/h264parser"
	"github.com/nareix/joy4/codec/h265parser"
	"github.com/nareix/joy4/codec/mjpegparser"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/hal/video/decode"
)

func CuvidInit(flags uint) error {
	ce := C.cuvidInit(C.uint(flags))
	return newCudaError(ce)
}

type VideoParser struct {
	ctx *C.decoder_context_t

	params    C.CUVIDPARSERPARAMS
	extParams C.CUVIDEOFORMATEX

	codec av.CodecType

	decodeBuffer []byte
	lengthSize   int
}

var (
	// There is a synchronization marker to define the boundaries of the NAL's
	// units. Each synchronization marker holds a value of 0x00 0x00 0x01
	// except to the very first one which is 0x00 0x00 0x00 0x01. If we run the
	// hexdump on the generated h264 bitstream, we can identify at least three
	// NALs in the beginning of the file.
	h264SequenceStart = []byte{0, 0, 0, 1}
)

func (v *VideoParser) initH264(stream h264parser.CodecData) {
	v.lengthSize = int(stream.RecordInfo.LengthSizeMinusOne) + 1
	// packet data already AnnexB
	if stream.RecordInfo.AnnexBPacket {
		v.lengthSize = 0
	}
	v.params.CodecType = C.cudaVideoCodec_H264
	var buf []byte
	// annexb
	buf = append(buf, h264SequenceStart...)
	buf = append(buf, stream.SPS()...)
	buf = append(buf, h264SequenceStart...)
	buf = append(buf, stream.PPS()...)
	// #nosec
	slice := (*[1024]byte)(unsafe.Pointer(&v.extParams.raw_seqhdr_data[0]))[:]
	copy(slice, buf)
	v.extParams.format.seqhdr_data_length = C.uint(len(buf))

}

func (v *VideoParser) initH265(stream h265parser.CodecData) {
	v.lengthSize = 4
	v.params.CodecType = C.cudaVideoCodec_HEVC
	var buf []byte
	// annexb
	buf = append(buf, h264SequenceStart...)
	buf = append(buf, stream.VPS()...)
	buf = append(buf, h264SequenceStart...)
	buf = append(buf, stream.SPS()...)
	buf = append(buf, h264SequenceStart...)
	buf = append(buf, stream.PPS()...)
	slice := (*[1024]byte)(unsafe.Pointer(&v.extParams.raw_seqhdr_data[0]))[:]
	copy(slice, buf)
	v.extParams.format.seqhdr_data_length = C.uint(len(buf))
}

func (v *VideoParser) initMJPEG(stream mjpegparser.CodecData) {
	v.params.CodecType = C.cudaVideoCodec_JPEG
}

func toCUVIDCodec(ty av.CodecType) C.cudaVideoCodec {
	switch ty {
	case av.MPEG1:
		return C.cudaVideoCodec_MPEG1
	case av.MPEG2:
		return C.cudaVideoCodec_MPEG2
	case av.MPEG4:
		return C.cudaVideoCodec_MPEG4
	case av.MJPEG:
		return C.cudaVideoCodec_JPEG
	case av.VP8:
		return C.cudaVideoCodec_VP8
	case av.VP9:
		return C.cudaVideoCodec_VP9
	default:
		return C.cudaVideoCodec_NumCodecs
	}
}

func (v *VideoParser) initFFMPEG(stream codec.FFMPEGVideoCodecData) error {
	c := toCUVIDCodec(stream.Type())
	if c >= C.cudaVideoCodec_NumCodecs {
		return ErrUnknownFormat
	}
	v.params.CodecType = c

	if stream.ExtraData != nil {
		slice := (*[1024]byte)(unsafe.Pointer(&v.extParams.raw_seqhdr_data[0]))[:]
		copy(slice, stream.ExtraData)
		v.extParams.format.seqhdr_data_length = C.uint(len(stream.ExtraData))
	}
	return nil
}

func NewVideoParser(stream av.CodecData, maxNumSurfaces int) (*VideoParser, error) {
	if maxNumSurfaces <= 0 {
		panic("maxNumSurfaces should be larger than 0")
	}
	vp := &VideoParser{}
	// vp.params.pExtVideoInfo = &vp.extParams
	vp.params.ulMaxNumDecodeSurfaces = C.uint(maxNumSurfaces)
	vp.params.ulMaxDisplayDelay = 1

	vp.codec = stream.Type()
	switch v := stream.(type) {
	case h264parser.CodecData:
		vp.initH264(v)
	case h265parser.CodecData:
		vp.initH265(v)
	case mjpegparser.CodecData:
		vp.initMJPEG(v)
	case codec.FFMPEGVideoCodecData:
		err := vp.initFFMPEG(v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnknownFormat
	}

	ce := C.init_decoder_context(&vp.ctx, &vp.params, &vp.extParams)
	if ce != 0 {
		return nil, newCudaError(ce)
	}

	return vp, nil
}

func (v *VideoParser) Push(pkt decode.CommonPacket, useTimeStamp bool) error {
	// log.Print(hex.Dump(pkt.Data[:16]))
	if len(pkt.GetData()) == 0 {
		// signal EOF
		ce := C.decoder_push_packet(v.ctx, nil, 0, 0, 0)
		return newCudaError(ce)
	}
	// TODO timestamp
	var flags C.uint
	var pts C.longlong
	if useTimeStamp {
		flags = flags | C.CUVID_PKT_TIMESTAMP
		pts = C.longlong(pkt.GetTime())
	}

	switch v.codec {
	case av.H264, av.H265:
		if v.lengthSize == 0 {
			ce := C.decoder_push_packet(v.ctx, unsafe.Pointer(&pkt.GetData()[0]), C.ulong(len(pkt.GetData())), pts, flags)
			return newCudaError(ce)
		}
		if cap(v.decodeBuffer) < len(pkt.GetData()) {
			v.decodeBuffer = make([]byte, 0, len(pkt.GetData())*2)
		}
		v.decodeBuffer = v.decodeBuffer[:0]
		if len(pkt.GetData()) < v.lengthSize {
			return ErrInvalidPacket
		}
		offset := 0
		for {
			if offset+v.lengthSize > len(pkt.GetData()) {
				break
			}
			l := 0
			if v.lengthSize == 1 {
				l = int(pkt.GetData()[offset])
			} else if v.lengthSize == 2 {
				l = int(binary.BigEndian.Uint16(pkt.GetData()[offset:]))
			} else if v.lengthSize == 4 {
				l = int(binary.BigEndian.Uint32(pkt.GetData()[offset:]))
			}
			offset += v.lengthSize
			if l == 0 || offset+l > len(pkt.GetData()) {
				return ErrInvalidPacket
			}
			v.decodeBuffer = append(v.decodeBuffer, h264SequenceStart...)
			v.decodeBuffer = append(v.decodeBuffer, pkt.GetData()[offset:offset+l]...)
			offset += l
		}
		// #nosec
		ce := C.decoder_push_packet(v.ctx, unsafe.Pointer(&v.decodeBuffer[0]), C.ulong(len(v.decodeBuffer)), pts, flags)
		return newCudaError(ce)
	default:
		// #nosec
		ce := C.decoder_push_packet(v.ctx, unsafe.Pointer(&pkt.GetData()[0]), C.ulong(len(pkt.GetData())), pts, flags)
		return newCudaError(ce)
	}
}

func (v *VideoParser) PushEOF() error {
	ce := C.decoder_push_packet(v.ctx, nil, 0, 0, 0)
	return newCudaError(ce)
}

func (v *VideoParser) pullFrame() (C.CuvidParsedFrame, bool, error) {
	var more C.int
	var frame C.CuvidParsedFrame
	ce := C.decoder_pull_frame(v.ctx, &frame, &more)
	if ce != 0 {
		return frame, false, newCudaError(ce)
	}
	return frame, more != 0, nil
}

type VideoMappedFrame struct {
	Data   DevicePointer
	Width  int
	Height int
	Pitch  uint
	Pts    int64
}

func (v *VideoParser) PullMappedFrame() (VideoMappedFrame, bool, error) {
	var params C.CUVIDPROCPARAMS
	var frame VideoMappedFrame
	parsedFrame, ok, err := v.pullFrame()
	if err != nil {
		return frame, false, err
	}
	if !ok {
		return frame, ok, nil
	}
	params.progressive_frame = parsedFrame.dispinfo.progressive_frame
	params.second_field = parsedFrame.second_field
	params.top_field_first = parsedFrame.dispinfo.top_field_first

	var pitch C.uint
	ce := C.wrapCuvidMapVideoFrame64(C.decoder_get_internal(v.ctx), parsedFrame.dispinfo.picture_index, (*C.CUdeviceptr)(&frame.Data), &pitch, &params)
	if ce != 0 {
		return frame, true, newCudaError(ce)
	}
	frame.Pitch = uint(pitch)
	frame.Width = int(parsedFrame.desc.width)
	frame.Height = int(parsedFrame.desc.height)
	frame.Pts = int64(parsedFrame.dispinfo.timestamp)
	return frame, true, nil
}

func (v *VideoParser) UnmapFrame(frame VideoMappedFrame) error {
	ce := C.wrapCuvidUnmapVideoFrame64(C.decoder_get_internal(v.ctx), C.CUdeviceptr(frame.Data))
	if ce != 0 {
		return newCudaError(ce)
	}
	return nil
}

func (v *VideoParser) Close() error {
	ce := C.close_decoder_context(v.ctx)
	return newCudaError(ce)
}
