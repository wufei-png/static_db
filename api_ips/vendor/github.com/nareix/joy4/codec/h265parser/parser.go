package h265parser

import (
	"bytes"
	"fmt"

	"github.com/nareix/bits"
	"github.com/nareix/bits/pio"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/codec/h264parser"
	"github.com/nareix/joy4/log"
)

// https://chromium.googlesource.com/chromium/src/media/+/master/video/h265_parser.h
const (
	NALU_TRAIL_N        = 0
	NALU_TRAIL_R        = 1
	NALU_TSA_N          = 2
	NALU_TSA_R          = 3
	NALU_STSA_N         = 4
	NALU_STSA_R         = 5
	NALU_RADL_N         = 6
	NALU_RADL_R         = 7
	NALU_RASL_N         = 8
	NALU_RASL_R         = 9
	NALU_RSV_VCL_N10    = 10
	NALU_RSV_VCL_R11    = 11
	NALU_RSV_VCL_N12    = 12
	NALU_RSV_VCL_R13    = 13
	NALU_RSV_VCL_N14    = 14
	NALU_RSV_VCL_R15    = 15
	NALU_BLA_W_LP       = 16
	NALU_BLA_W_RADL     = 17
	NALU_BLA_N_LP       = 18
	NALU_IDR_W_RADL     = 19
	NALU_IDR_N_LP       = 20
	NALU_CRA_NUT        = 21
	NALU_RSV_IRAP_VCL22 = 22
	NALU_RSV_IRAP_VCL23 = 23
	NALU_RSV_VCL24      = 24
	NALU_RSV_VCL25      = 25
	NALU_RSV_VCL26      = 26
	NALU_RSV_VCL27      = 27
	NALU_RSV_VCL28      = 28
	NALU_RSV_VCL29      = 29
	NALU_RSV_VCL30      = 30
	NALU_RSV_VCL31      = 31
	NALU_VPS_NUT        = 32
	NALU_SPS_NUT        = 33
	NALU_PPS_NUT        = 34
	NALU_AUD_NUT        = 35
	NALU_EOS_NUT        = 36
	NALU_EOB_NUT        = 37
	NALU_FD_NUT         = 38
	NALU_PREFIX_SEI_NUT = 39
	NALU_SUFFIX_SEI_NUT = 40
	NALU_RSV_NVCL41     = 41
	NALU_RSV_NVCL42     = 42
	NALU_RSV_NVCL43     = 43
	NALU_RSV_NVCL44     = 44
	NALU_RSV_NVCL45     = 45
	NALU_RSV_NVCL46     = 46
	NALU_RSV_NVCL47     = 47
)

const (
	HEVC_SEI_TYPE_BUFFERING_PERIOD                     = 0
	HEVC_SEI_TYPE_PICTURE_TIMING                       = 1
	HEVC_SEI_TYPE_PAN_SCAN_RECT                        = 2
	HEVC_SEI_TYPE_FILLER_PAYLOAD                       = 3
	HEVC_SEI_TYPE_USER_DATA_REGISTERED_ITU_T_T35       = 4
	HEVC_SEI_TYPE_USER_DATA_UNREGISTERED               = 5
	HEVC_SEI_TYPE_RECOVERY_POINT                       = 6
	HEVC_SEI_TYPE_SCENE_INFO                           = 9
	HEVC_SEI_TYPE_FULL_FRAME_SNAPSHOT                  = 15
	HEVC_SEI_TYPE_PROGRESSIVE_REFINEMENT_SEGMENT_START = 16
	HEVC_SEI_TYPE_PROGRESSIVE_REFINEMENT_SEGMENT_END   = 17
	HEVC_SEI_TYPE_FILM_GRAIN_CHARACTERISTICS           = 19
	HEVC_SEI_TYPE_POST_FILTER_HINT                     = 22
	HEVC_SEI_TYPE_TONE_MAPPING_INFO                    = 23
	HEVC_SEI_TYPE_FRAME_PACKING                        = 45
	HEVC_SEI_TYPE_DISPLAY_ORIENTATION                  = 47
	HEVC_SEI_TYPE_SOP_DESCRIPTION                      = 128
	HEVC_SEI_TYPE_ACTIVE_PARAMETER_SETS                = 129
	HEVC_SEI_TYPE_DECODING_UNIT_INFO                   = 130
	HEVC_SEI_TYPE_TEMPORAL_LEVEL0_INDEX                = 131
	HEVC_SEI_TYPE_DECODED_PICTURE_HASH                 = 132
	HEVC_SEI_TYPE_SCALABLE_NESTING                     = 133
	HEVC_SEI_TYPE_REGION_REFRESH_INFO                  = 134
	HEVC_SEI_TYPE_MASTERING_DISPLAY_INFO               = 137
	HEVC_SEI_TYPE_CONTENT_LIGHT_LEVEL_INFO             = 144
	HEVC_SEI_TYPE_ALTERNATIVE_TRANSFER_CHARACTERISTICS = 147
)

type SliceType int

const (
	HEVC_SLICE_B SliceType = 0
	HEVC_SLICE_P           = 1
	HEVC_SLICE_I           = 2
)

/**
 * 7.4.2.1
 */
const (
	HEVC_MAX_SUB_LAYERS           = 7
	HEVC_MAX_VPS_COUNT            = 16
	HEVC_MAX_SPS_COUNT            = 32
	HEVC_MAX_PPS_COUNT            = 256
	HEVC_MAX_SHORT_TERM_RPS_COUNT = 64
	HEVC_MAX_CU_SIZE              = 128

	HEVC_MAX_REFS     = 16
	HEVC_MAX_DPB_SIZE = 16 // A.4.1

	HEVC_MAX_LOG2_CTB_SIZE = 6
)

const (
	FF_PROFILE_HEVC_MAIN               = 1
	FF_PROFILE_HEVC_MAIN_10            = 2
	FF_PROFILE_HEVC_MAIN_STILL_PICTURE = 3
	FF_PROFILE_HEVC_REXT               = 4
)

type HEVCDecoderConfigurationRecord struct {
	ConfigurationVersion              uint8
	GeneralTierFlag                   uint8
	GeneralProfileIdc                 uint8
	GeneralProfileCompatiabilityFlags uint32
	GeneralConstraintIndicatorFlags   uint64
	GeneralLevelIdc                   uint8
	MinSpatialSegmentationIdc         uint16
	ParallelismType                   uint8
	ChromaFormat                      uint8
	BitDepthLumaMinus8                uint8
	BitDepthChromaMinus8              uint8
	AvgFrameRate                      uint16
	ConstantFrameRate                 uint8
	NumTemporalLayers                 uint8
	TemporalIdNested                  uint8
	LengthSizeMinusOne                uint8
	NumOfArrays                       uint8

	VPS [][]byte
	SPS [][]byte
	PPS [][]byte
}

var ErrDecconfInvalid = fmt.Errorf("h265parser: HEVCDecoderConfRecord invalid")
var ErrDecconfUnsupport = fmt.Errorf("h265parser: HEVCDecoderConfRecord unsupported")

func (self *HEVCDecoderConfigurationRecord) Unmarshal(b []byte) (n int, err error) {
	if len(b) < 1 {
		err = ErrDecconfInvalid
		return
	}
	// XXX ff_isom_write_hvcc
	if b[0] != 1 {
		err = ErrDecconfUnsupport
		return
	}

	if len(b) < 23 {
		err = ErrDecconfInvalid
		return
	}

	var profileIdc, misc uint8
	var generalConstraintIndicatorFlagsHi uint32
	var generalConstraintIndicatorFlagsLo uint16

	pos := 0
	self.ConfigurationVersion = pio.U8(b[pos:])
	pos++
	profileIdc = pio.U8(b[pos:])
	pos++
	self.GeneralProfileCompatiabilityFlags = pio.U32BE(b[pos:])
	pos += 4
	generalConstraintIndicatorFlagsHi = pio.U32BE(b[pos:])
	pos += 4
	generalConstraintIndicatorFlagsLo = pio.U16BE(b[pos:])
	pos += 2
	self.GeneralLevelIdc = pio.U8(b[pos:])
	pos++
	self.MinSpatialSegmentationIdc = pio.U16BE(b[pos:])
	pos += 2
	self.ParallelismType = pio.U8(b[pos:])
	pos++
	self.ChromaFormat = pio.U8(b[pos:])
	pos++
	self.BitDepthLumaMinus8 = pio.U8(b[pos:])
	pos++
	self.BitDepthChromaMinus8 = pio.U8(b[pos:])
	pos++
	self.AvgFrameRate = pio.U16BE(b[pos:])
	pos += 2
	misc = pio.U8(b[pos:])
	pos++
	self.NumOfArrays = pio.U8(b[pos:])
	pos++

	self.GeneralProfileIdc = profileIdc >> 6
	self.GeneralTierFlag = (profileIdc >> 5) & 1
	self.GeneralProfileIdc = profileIdc & 0x1f
	self.GeneralConstraintIndicatorFlags = uint64(generalConstraintIndicatorFlagsHi)
	self.GeneralConstraintIndicatorFlags <<= 16
	self.GeneralConstraintIndicatorFlags |= uint64(generalConstraintIndicatorFlagsLo)
	self.MinSpatialSegmentationIdc &= 0xfff
	self.ParallelismType &= 3
	self.ChromaFormat &= 3
	self.BitDepthLumaMinus8 &= 7
	self.BitDepthChromaMinus8 &= 7
	self.ConstantFrameRate = misc >> 6
	self.NumTemporalLayers = (misc >> 3) & 7
	self.TemporalIdNested = (misc >> 2) & 1
	self.LengthSizeMinusOne = misc & 3

	numArrays := int(self.NumOfArrays)
	if pos != 23 {
		panic("should be 23")
	}
	// pos := 23
	var psPos [3]int
	for i := 0; i < numArrays; i++ {
		if pos+3 > len(b) {
			err = ErrDecconfInvalid
			return
		}
		naluType := b[pos] & 0x3f

		if naluType == NALU_VPS_NUT {
			psPos[0] = pos
		} else if naluType == NALU_SPS_NUT {
			psPos[1] = pos
		} else if naluType == NALU_PPS_NUT {
			psPos[2] = pos
		}
		numNalus := int(pio.U16BE(b[pos+1:]))
		pos += 3
		for j := 0; j < numNalus; j++ {
			if pos+2 > len(b) {
				err = ErrDecconfInvalid
				return
			}
			l := int(pio.U16BE(b[pos:]))
			pos += 2
			if pos+l > len(b) {
				err = ErrDecconfInvalid
				return
			}
			pos += l
		}
	}
	if psPos[0] == 0 || psPos[1] == 0 || psPos[2] == 0 {
		err = ErrDecconfInvalid
		return
	}

	for i := 0; i < 3; i++ {
		pos = psPos[i]
		numNalus := int(pio.U16BE(b[pos+1:]))
		pos += 3

		for j := 0; j < numNalus; j++ {
			l := int(pio.U16BE(b[pos:]))
			pos += 2
			buf := make([]byte, l)
			copy(buf, b[pos:pos+l])
			if i == 0 {
				self.VPS = append(self.VPS, buf)
			} else if i == 1 {
				self.SPS = append(self.SPS, buf)
			} else if i == 2 {
				self.PPS = append(self.PPS, buf)
			}
		}
	}
	return pos, nil
}

type Window struct {
	LeftOffset   uint
	RightOffset  uint
	TopOffset    uint
	BottomOffset uint
}

type TemporalLayer struct {
	MaxDecPicBuffering int
	NumReorderPics     int
	MaxLatencyIncrease int
}

type VUI struct {
}

type PTLCommon struct {
	ProfileSpace             uint8
	TierFlag                 uint8
	ProfileIdc               uint8
	ProfileCompatibilityFlag [32]uint8
	LevelIdc                 uint8
	ProgressiveSourceFlag    uint8
	InterlacedSourceFlag     uint8
	NonPackedConstraintFlag  uint8
	FrameOnlyConstraintFlag  uint8
}

type PTL struct {
	GeneralPTL  PTLCommon
	SubLayerPTL [HEVC_MAX_SUB_LAYERS]PTLCommon

	SubLayerProfilePresentFlag [HEVC_MAX_SUB_LAYERS]uint8
	SubLayerLevelPresentFlag   [HEVC_MAX_SUB_LAYERS]uint8
}

type ShortTermRPS struct {
	NumNegativePics    uint
	NumDeltaPocs       int
	RpsIdxNumDeltaPocs int
	DeltaPoc           [32]int32
	Used               [32]uint8
}

type ScalingList struct {
	/* This is a little wasteful, since sizeID 0 only needs 8 coeffs,
	 * and size ID 3 only has 2 arrays, not 6. */
	Sl   [4][6][64]uint8
	SlDc [2][6]uint8
}

type PCM struct {
	BitDepth              uint8
	BitDepthChroma        uint8
	Log2MinPcmCbSize      uint
	Log2MaxPcmCbSize      uint
	LoopFilterDisableFlag uint8
}

type VPSInfo struct {
	TemporalIDNestingFlag uint8
	MaxLayers             uint
	MaxSubLayers          uint

	PTL                             PTL
	SubLayerOrderingInfoPersentFlag uint8
	MaxDecPicBuffering              [HEVC_MAX_SUB_LAYERS]uint
	NumReorderPics                  [HEVC_MAX_SUB_LAYERS]uint
	MaxLatencyIncrease              [HEVC_MAX_SUB_LAYERS]uint
	MaxLayerID                      uint
	NumLayerSets                    uint
	TimingInfoPresentFlag           uint8
	NumUnitsInTick                  uint
	TimeScale                       uint
	PocProportionalToTimingFlag     uint8
	NumTicksPocDiffOne              uint
	NumHrdParameters                uint

	Data []byte
}

type SPSInfo struct {
	VpsID                   uint
	ChromaFormatIdc         int
	SeperateColourPlaneFlag uint8

	OutputWindow Window
	PicConfWin   Window

	BitDepth       uint
	BitDepthChroma uint
	PixelShift     int
	PixFmt         int

	Log2MaxPocLsb  uint
	PCMEnabledFlag int

	MaxSubLayers          uint
	TemporalLayer         []TemporalLayer
	TemporalIDNestingFlag uint8

	YUI VUI
	PTL PTL

	ScalingListEnableFlag uint8
	ScalingList           ScalingList

	NbStRps uint
	StRps   [HEVC_MAX_SHORT_TERM_RPS_COUNT]ShortTermRPS

	AmpEnabledFlag uint8
	SaoEnabled     uint8

	LongTermRefPicsPresentFlag uint8
	LtRefPicPocLsbSPS          [32]uint16
	UsedByCurrPicLtSPSFlag     [32]uint8
	NumLongTermRefPicsSPS      uint8

	PCM                               PCM
	SPSTemporalMvpEnabledFlag         uint8
	SPSStrongIntraSmoothingEnableFlag uint8

	Log2MinCbSize                 uint
	Log2DiffMaxMinCodingBlockSize uint
	Log2MinTbSize                 uint
	Log2MaxTrafoSize              uint
	Log2CtbSize                   uint
	Log2MinPuSize                 uint

	MaxTransformHierarchyDepthInter int
	MaxTransformHierarchyDepthIntra int

	TransformSkipRotationEnabledFlag    int
	TransformSkipContextEnabledFlag     int
	ImplicitRdpcmEnabledFlag            int
	ExplicitRdpcmEnabledFlag            int
	IntraSmoothingDisabledFlag          int
	PersistentRiceAdaptationEnabledFlag int

	///< coded frame dimension in various units
	Width       uint
	Height      uint
	CtbWidth    int
	CtbHeight   int
	CtbSize     int
	MinCbWidth  int
	MinCbHeight int
	MinTbWidth  int
	MinTbHeight int
	MinPuWidth  int
	TbMask      int
	MinPuHeight int

	HShift [3]int
	VShift [3]int

	QpBdOffset int
	Data       []byte
}

type CodecData struct {
	Record     []byte
	RecordInfo HEVCDecoderConfigurationRecord
	SPSInfo    SPSInfo
	VPSInfo    VPSInfo
}

func (self CodecData) Type() av.CodecType {
	return av.H265
}

func (self CodecData) SPS() []byte {
	return self.RecordInfo.SPS[0]
}

func (self CodecData) PPS() []byte {
	return self.RecordInfo.PPS[0]
}

func (self CodecData) VPS() []byte {
	return self.RecordInfo.VPS[0]
}

func (self CodecData) Width() int {
	return int(self.SPSInfo.Width)
}

func (self CodecData) Height() int {
	return int(self.SPSInfo.Height)
}

func decodeProfileTierLevel(r *bits.GolombBitReader) (ptl PTLCommon, err error) {
	var t uint
	if t, err = r.ReadBits(2); err != nil {
		return
	}
	ptl.ProfileSpace = uint8(t)
	if ptl.TierFlag, err = readBitToUint8(r); err != nil {
		return
	}
	if ptl.ProfileIdc, err = readNBitToUint8(r, 5); err != nil {
		return
	}
	switch ptl.ProfileIdc {
	case FF_PROFILE_HEVC_MAIN:
		log.Log(log.DEBUG, "Main profile bitstream")
	case FF_PROFILE_HEVC_MAIN_10:
		log.Log(log.DEBUG, "Main 10 profile bitstream")
	case FF_PROFILE_HEVC_MAIN_STILL_PICTURE:
		log.Log(log.DEBUG, "Main Still Picture profile bitstream")
	case FF_PROFILE_HEVC_REXT:
		log.Log(log.DEBUG, "Range Extension profile bitstream")
	default:
		log.Log(log.WARN, "Unknown HEVC profile: ", ptl.ProfileIdc)
	}

	for i := 0; i < 32; i++ {
		if ptl.ProfileCompatibilityFlag[i], err = readBitToUint8(r); err != nil {
			return
		}
		if ptl.ProfileIdc == 0 && i > 0 && ptl.ProfileCompatibilityFlag[i] != 0 {
			ptl.ProfileIdc = uint8(i)
		}
	}

	if ptl.ProgressiveSourceFlag, err = readBitToUint8(r); err != nil {
		return
	}
	if ptl.InterlacedSourceFlag, err = readBitToUint8(r); err != nil {
		return
	}
	if ptl.NonPackedConstraintFlag, err = readBitToUint8(r); err != nil {
		return
	}
	if ptl.FrameOnlyConstraintFlag, err = readBitToUint8(r); err != nil {
		return
	}

	// XXX_reserved_zero_44bits[0..15]
	if _, err = r.ReadBits(16); err != nil {
		return
	}
	// XXX_reserved_zero_44bits[16..31]
	if _, err = r.ReadBits(16); err != nil {
		return
	}
	// XXX_reserved_zero_44bits[32..43]
	if _, err = r.ReadBits(12); err != nil {
		return
	}

	return
}

func ParsePTL(r *bits.GolombBitReader, maxSubLayers int) (ptl PTL, err error) {
	if ptl.GeneralPTL, err = decodeProfileTierLevel(r); err != nil {
		return
	}

	if ptl.GeneralPTL.LevelIdc, err = readNBitToUint8(r, 8); err != nil {
		return
	}

	for i := 0; i < maxSubLayers-1; i++ {
		if ptl.SubLayerProfilePresentFlag[i], err = readBitToUint8(r); err != nil {
			return
		}
		if ptl.SubLayerLevelPresentFlag[i], err = readBitToUint8(r); err != nil {
			return
		}
	}

	if maxSubLayers-1 > 0 {
		for i := maxSubLayers - 1; i < 8; i++ {
			// reserved_zero_2bits[i]
			if _, err = r.ReadBits(2); err != nil {
				return
			}
		}
	}
	for i := 0; i < maxSubLayers-1; i++ {
		if ptl.SubLayerProfilePresentFlag[i] != 0 {
			if ptl.SubLayerPTL[i], err = decodeProfileTierLevel(r); err != nil {
				err = fmt.Errorf("PTL information for sublayer %d too short", i)
				return
			}
		}
		if ptl.SubLayerLevelPresentFlag[i] != 0 {
			if ptl.SubLayerPTL[i].LevelIdc, err = readNBitToUint8(r, 8); err != nil {
				err = fmt.Errorf("Not enough data for sublayer %d level_idc", i)
				return
			}
		}
	}

	return
}

func ParseVPS(data []byte) (self VPSInfo, err error) {
	log.Log(log.DEBUG, "Decode VPS")
	var t uint
	r := &bits.GolombBitReader{R: bytes.NewReader(data)}
	var vpsID uint
	if vpsID, err = r.ReadBits(4); err != nil {
		return
	}
	if vpsID >= HEVC_MAX_VPS_COUNT {
		err = fmt.Errorf("VPS id out of range: %d", vpsID)
		return
	}

	if t, err = r.ReadBits(2); err != nil {
		return
	}
	if t != 3 {
		err = fmt.Errorf("vps_reserved_three_2bits is not three")
		return
	}

	if self.MaxLayers, err = r.ReadBits(6); err != nil {
		return
	}
	self.MaxLayers++
	if self.MaxSubLayers, err = r.ReadBits(3); err != nil {
		return
	}
	self.MaxSubLayers++
	if self.TemporalIDNestingFlag, err = readBitToUint8(r); err != nil {
		return
	}

	if t, err = r.ReadBits(16); err != nil {
		return
	}
	if t != 0xffff {
		err = fmt.Errorf("vps_reserved_ffff_16bits is not 0xffff")
		return
	}

	if self.MaxSubLayers > HEVC_MAX_SUB_LAYERS {
		err = fmt.Errorf("vps_max_sub_layers out of range: %d", self.MaxSubLayers)
		return
	}

	if self.PTL, err = ParsePTL(r, int(self.MaxSubLayers)); err != nil {
		return
	}

	if self.SubLayerOrderingInfoPersentFlag, err = readBitToUint8(r); err != nil {
		return
	}
	i := self.MaxSubLayers - 1
	if self.SubLayerOrderingInfoPersentFlag != 0 {
		i = 0
	}
	for ; i < self.MaxSubLayers; i++ {
		if self.MaxDecPicBuffering[i], err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.MaxDecPicBuffering[i]++
		if self.NumReorderPics[i], err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		if self.MaxLatencyIncrease[i], err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.MaxLatencyIncrease[i] -= 1
		if self.MaxDecPicBuffering[i] > HEVC_MAX_DPB_SIZE || self.MaxDecPicBuffering[i] == 0 {
			err = fmt.Errorf("vps_max_dec_pic_buffering_minus1 out of range: %d", self.MaxDecPicBuffering[i]-1)
			return
		}
		if self.NumReorderPics[i] > self.MaxDecPicBuffering[i]-1 {
			log.Log(log.WARN, "vps_max_num_reorder_pics out of range: %d", self.NumReorderPics[i])
			// err_recpgntion
		}
	}

	if self.MaxLayerID, err = r.ReadBits(6); err != nil {
		return
	}
	if self.NumLayerSets, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	self.NumLayerSets++
	if self.NumLayerSets < 1 || self.NumLayerSets > 1024 {
		err = fmt.Errorf("too many layer_id_included_flags")
		return
	}

	for i := uint(0); i < self.NumLayerSets; i++ {
		for j := uint(0); j <= self.MaxLayerID; j++ {
			// layer_id_included_flag[i][j]
			if _, err = r.ReadBit(); err != nil {
				return
			}
		}
	}
	if self.TimingInfoPresentFlag, err = readBitToUint8(r); err != nil {
		return
	}

	if self.TimingInfoPresentFlag != 0 {
		if self.NumUnitsInTick, err = r.ReadBits(32); err != nil {
			return
		}
		if self.TimeScale, err = r.ReadBits(32); err != nil {
			return
		}
		if self.PocProportionalToTimingFlag, err = readBitToUint8(r); err != nil {
			return
		}
		if self.PocProportionalToTimingFlag != 0 {
			if self.NumTicksPocDiffOne, err = r.ReadExponentialGolombCode(); err != nil {
				return
			}
			self.NumTicksPocDiffOne++
			if self.NumHrdParameters, err = r.ReadExponentialGolombCode(); err != nil {
				return
			}
			if self.NumHrdParameters > self.NumLayerSets {
				err = fmt.Errorf("vps_num_hrd_parameters %d is invalid", self.NumHrdParameters)
				return
			}
			for i := uint(0); i < self.NumHrdParameters; i++ {
				// TODO
			}
		}
	}
	// TODO

	return
}

func ParseSPS(data []byte) (self SPSInfo, err error) {
	log.Log(log.DEBUG, "Decode SPS")
	var t uint
	r := &bits.GolombBitReader{R: bytes.NewReader(data)}
	if self.VpsID, err = r.ReadBits(4); err != nil {
		return
	}
	if self.VpsID >= HEVC_MAX_VPS_COUNT {
		err = fmt.Errorf("VPS id out of range: ", self.VpsID)
		return
	}
	// TODO check vps list

	if t, err = r.ReadBits(3); err != nil {
		return
	}
	self.MaxSubLayers = t + 1
	if self.MaxSubLayers > HEVC_MAX_SUB_LAYERS {
		err = fmt.Errorf("sps_max_sub_layers out of range: ", self.MaxSubLayers)
		return
	}

	if self.TemporalIDNestingFlag, err = readBitToUint8(r); err != nil {
		return
	}

	if self.PTL, err = ParsePTL(r, int(self.MaxSubLayers)); err != nil {
		return
	}

	var spsID uint
	if spsID, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	if spsID >= HEVC_MAX_SPS_COUNT {
		err = fmt.Errorf("SPS id out of range: %d", spsID)
		return
	}

	if t, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	self.ChromaFormatIdc = int(t)
	if self.ChromaFormatIdc > 3 {
		err = fmt.Errorf("chroma_format_idc %d is invalid", self.ChromaFormatIdc)
		return
	}

	if self.ChromaFormatIdc == 3 {
		if self.SeperateColourPlaneFlag, err = readBitToUint8(r); err != nil {
			return
		}
	}

	if self.SeperateColourPlaneFlag != 0 {
		self.ChromaFormatIdc = 0
	}

	if self.Width, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	if self.Height, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	// TODO check image size

	// pic_conformance_flag
	var picConformanceFlag uint
	if picConformanceFlag, err = r.ReadBit(); err != nil {
		return
	}
	if picConformanceFlag != 0 {
		vertMult := uint(1)
		if self.ChromaFormatIdc < 2 {
			vertMult++
		}
		horizMult := uint(1)
		if self.ChromaFormatIdc < 3 {
			horizMult++
		}
		if self.PicConfWin.LeftOffset, err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.PicConfWin.LeftOffset *= horizMult
		if self.PicConfWin.RightOffset, err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.PicConfWin.RightOffset *= horizMult
		if self.PicConfWin.TopOffset, err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.PicConfWin.TopOffset *= vertMult
		if self.PicConfWin.BottomOffset, err = r.ReadExponentialGolombCode(); err != nil {
			return
		}
		self.PicConfWin.BottomOffset *= vertMult
		// ignore crop
		self.OutputWindow = self.PicConfWin
	}

	if self.BitDepth, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	self.BitDepth += 8
	var bitDepthChroma uint
	if bitDepthChroma, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	bitDepthChroma += 8
	if self.ChromaFormatIdc != 0 && bitDepthChroma != self.BitDepth {
		log.Logf(log.ERROR, "Luma bit depth (%d) is different from chroma bit depth (%d), this is unsupported.", self.BitDepth, bitDepthChroma)
	}
	self.BitDepthChroma = bitDepthChroma
	// map pixel format

	if self.Log2MaxPocLsb, err = r.ReadExponentialGolombCode(); err != nil {
		return
	}
	self.Log2MaxPocLsb += 4
	if self.Log2MaxPocLsb > 16 {
		err = fmt.Errorf("log2_max_pic_order_cnt_lsb_minus4 out range: %d", self.Log2MaxPocLsb)
		return
	}

	// TODO
	return
}

func NewCodecDataFromPS(vps, sps, pps []byte) (self CodecData, err error) {
	vpsNal, _ := h264parser.ExtractRBSP(vps, true)
	vpsData := vpsNal.Rbsp
	// XXX ffplay VPS has NAL header
	if len(vpsData) >= 2 || vpsData[0]&0x3f == NALU_VPS_NUT {
		vpsData = vpsData[2:]
	}
	if self.VPSInfo, err = ParseVPS(vpsData); err != nil {
		return
	}
	spsNal, _ := h264parser.ExtractRBSP(sps, true)
	if self.SPSInfo, err = ParseSPS(spsNal.Rbsp); err != nil {
		return
	}
	record := HEVCDecoderConfigurationRecord{
		LengthSizeMinusOne: 3,
	}
	record.VPS = [][]byte{vps}
	record.SPS = [][]byte{sps}
	record.PPS = [][]byte{pps}
	self.RecordInfo = record
	// log.Log(log.DEBUG, "VPS: ", self.VPSInfo)
	// log.Log(log.DEBUG, "SPS: ", self.SPSInfo)
	return
}

func NewCodecDataFromHEVCDecoderConfRecord(record []byte) (self CodecData, err error) {
	r := HEVCDecoderConfigurationRecord{
		LengthSizeMinusOne: 3,
	}
	if _, err = r.Unmarshal(record); err != nil {
		return
	}
	if len(r.VPS) == 0 {
		err = fmt.Errorf("h264parser: no VPS found in AVCDecoderConfRecord")
		return
	}
	if len(r.SPS) == 0 {
		err = fmt.Errorf("h264parser: no SPS found in AVCDecoderConfRecord")
		return
	}
	if len(r.PPS) == 0 {
		err = fmt.Errorf("h264parser: no PPS found in AVCDecoderConfRecord")
		return
	}
	if self, err = NewCodecDataFromPS(r.VPS[0], r.SPS[0], r.PPS[0]); err != nil {
		return
	}
	self.RecordInfo = r
	self.Record = record
	return
}
