package tool

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"

	"strings"

	"gitlab.sz.sensetime.com/viper/gosdkwrapper/feature"
)

type Request struct {
	Image                Image  `json:"image"`
	Face_selection       string `json:"face_selection"`
	Auto_rotation_thresh string `json:"auto_rotation_thresh"`
}
type Image struct {
	Format string `json:"format"`
	Data   string `json:"data"`
}
type Req struct {
	Request     []Request `json:"requests"`
	Detect_mode string    `json:"detect_mode"`
	Face_type   string    `json:"face_type"`
}

func StructTojson_ips(path string) string { //ips request json
	//实例化一个数据结构，用于生成json字符串
	image := Image{
		Format: "IMAGE_UNKNOWN",
		Data:   ChangeTobase(path),
	}

	request := Request{
		Image:                image,
		Face_selection:       "LargestFace",
		Auto_rotation_thresh: "1",
	}
	var requests []Request
	requests = append(requests, request)
	req := Req{
		Request:     requests,
		Detect_mode: "Default",
		Face_type:   "Large",
	}
	//fmt.Print(req)
	//Marshal失败时err!=nil
	reqjson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}
	//fmt.Println(string(reqjson))
	//jsonStu是[]byte类型，转化成string类型便于查看
	return string(reqjson)
}
func Post_req(url string, reader string) []byte {

	url1 := url
	method := "POST"

	payload := strings.NewReader(reader)
	//fmt.Print(payload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url1, payload)

	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(string(body))
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte("")
	}
	return body
}
func Interface2String(inter interface{}) {

	switch inter.(type) {

	case string:
		fmt.Println("string", inter.(string))

		break
	case int:
		fmt.Println("int", inter.(int))
		break
	case float64:
		fmt.Println("float64", inter.(float64))
		break
	}

}

// fmt.Println(string(body))
// var wf string
// wf = "U09GMahhAAAQAgAAAAEAAAAAAAAFAAAAAAAAAAAAAAASh1go3N8Aoh5rqugx5Ex5AcHPdST65xOB+B+vPIMnv/ZWitN2YytqSNf2NZsuV0wH1KMa2Sxo0opSZDi9imNSir5Q32YSM+wV3/SS52nVj0TZlb78ozaj37z6pr12Cc7xTRjE0ImvB1rSSTOkpLnZgF3ESwgm4DcZU4Apa4/yvsnn4I0hqqt9HGliY17qI5FKmsklYrx/3NRBKpzBILbYacF9toTelJrhmv6gzT0F08A9Tl4hB4CHOnPDKh7lmg8ccxSLCrhtPGNaNXcTnZvEpqOn+A67QywWFNKJjd/VjPT/O+tAtAkyc/zqxyLZ6kcsqCh4ncDgaLkzi6pjZKRJuamVksAFRjRvD55u4dlNjgjlnXPXOJQWE/of4EJ2DH1dfq88xsm39oVGorr+08uUTSHvZuMj3n5x73/62VWzUorLqdaETc9lpeVxa8DFADSDz2rukTlr0nlbz0ua0nVzgAjUQws7UxDb8q0u19s4auIi/W9IbkGPVzvwHG14nuRz57Kj6zeFQX77J4JOyDpzcFMSZ5STzfhq7vA4uTn9hL0WmoSbqIMv43u/Sqp4CcAX48Nc31UCvSIofVAkmjeZ/pTTvJp/h7q5k2R6fuE8E3MUdrTwMpjWzHWdfHJQwf194LU8wuG03WTAYfgs42UVw4xz2FroLqgyAjPRo8+K3uJVbyj1kyicb1WItOCxPM8="
// data, err := getRawFeatureFromBase64(wf)
// if err != nil {
// 	fmt.Errorf("出问题")
// } else {
// 	fmt.Print(data)

func GetRawFeatureFromBase64(feat string) (*feature.RawFeature, error) {
	blob, err := base64.StdEncoding.DecodeString(feat)
	if err != nil {
		return nil, err
	}
	pf, err := feature.NewPersistedFeatureFromBytes(blob, false)
	if err != nil {
		return nil, err
	}
	fdec, err := feature.NewDecoder([]byte("1234567890123456"))
	if err != nil {
		return nil, err
	}
	raw, err := fdec.Decode(pf)
	if err != nil {
		return nil, err
	}

	return &raw, nil
}

type Res_ips struct {
	Results []struct {
		Code   int32  `json:"code"`
		Error  string `json:"error"`
		Status string `json:"status"`
	} `json:"results"`
	Responses []struct {
		FaceInfo struct {
			Type string `json:"type"`
			Face struct {
				Quality   float32 `json:"quality"`
				Rectangle struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				TrackID string `json:"track_id"`
				Angle   struct {
					Yaw   json.Number `json:"yaw"`
					Pitch json.Number `json:"pitch"`
					Roll  json.Number `json:"roll"`
				} `json:"angle"`
				Landmarks []struct {
					X int32 `json:"x"`
					Y int32 `json:"y"`
				} `json:"landmarks"`
				Attributes struct {
				} `json:"attributes"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
				FaceScore json.Number `json:"face_score"`
			} `json:"face"`
			Pedestrian struct {
				Quality   float32 `json:"quality"`
				Rectangle struct {
					Vertices []struct {
						X string `json:"x"`
						Y string `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				TrackID             string `json:"track_id"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
				PedestrianScore float32 `json:"pedestrian_score"`
			} `json:"pedestrian"`
			Automobile struct {
				Quality   float32 `json:"quality"`
				Rectangle struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				TrackID             string `json:"track_id"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
			} `json:"automobile"`
			HumanPoweredVehicle struct {
				Quality   float32 `json:"quality"`
				Rectangle struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				TrackID             string `json:"track_id"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
			} `json:"human_powered_vehicle"`
			Cyclist struct {
				Quality   float32 `json:"quality"`
				Rectangle struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				TrackID             string `json:"track_id"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
			} `json:"cyclist"`
			Crowd struct {
				Quantity string `json:"quantity"`
				Incident []struct {
					ID         string `json:"id"`
					Type       string `json:"type"`
					Status     string `json:"status"`
					StartTime  string `json:"start_time"`
					StopTime   string `json:"stop_time"`
					UpdateTime string `json:"update_time"`
					UUID       string `json:"uuid"`
				} `json:"incident"`
				DensitySize struct {
					Width  int32 `json:"width"`
					Height int32 `json:"height"`
				} `json:"density_size"`
				DensityMap string `json:"density_map"`
				StrandMap  struct {
					Format string `json:"format"`
					Data   string `json:"data"`
					URL    string `json:"url"`
				} `json:"strand_map"`
				FullHeadTargets struct {
					HeadTargets []struct {
						Coordinate struct {
							X float32 `json:"x"`
							Y float32 `json:"y"`
						} `json:"coordinate"`
						Rectangle struct {
							Vertices []struct {
								X int32 `json:"x"`
								Y int32 `json:"y"`
							} `json:"vertices"`
						} `json:"rectangle"`
					} `json:"head_targets"`
				} `json:"full_head_targets"`
			} `json:"crowd"`
			Event struct {
				EventID string `json:"event_id"`
				Rule    struct {
					Type   string `json:"type"`
					RuleID string `json:"rule_id"`
					Roi    struct {
						Vertices []struct {
							X float32 `json:"x"`
							Y float32 `json:"y"`
						} `json:"vertices"`
					} `json:"roi"`
					Duration  string `json:"duration"`
					Direction struct {
						X float32 `json:"x"`
						Y float32 `json:"y"`
					} `json:"direction"`
				} `json:"rule"`
				EventStatus string `json:"event_status"`
				EventType   string `json:"event_type"`
				Rectangle   struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
				AttributesWithScore struct {
				} `json:"attributes_with_score"`
			} `json:"event"`
			PortraitImageLocation struct {
				PanoramicImageSize struct {
					Width  int32 `json:"width"`
					Height int32 `json:"height"`
				} `json:"panoramic_image_size"`
				PortraitImageInPanoramic struct {
					Vertices []struct {
						X string `json:"x"`
						Y string `json:"y"`
					} `json:"vertices"`
				} `json:"portrait_image_in_panoramic"`
				PortraitInPanoramic struct {
					Vertices []struct {
						X string `json:"x"`
						Y string `json:"y"`
					} `json:"vertices"`
				} `json:"portrait_in_panoramic"`
			} `json:"portrait_image_location"`
			ObjectID     string `json:"object_id"`
			Associations []struct {
				Type            string `json:"type"`
				ObjectID        string `json:"object_id"`
				AssociationType string `json:"association_type"`
			} `json:"associations"`
			Algo struct {
				AppName       string `json:"app_name"`
				AppVersion    int32  `json:"app_version"`
				ObjectType    string `json:"object_type"`
				ObjectVersion int32  `json:"object_version"`
				Data          struct {
					TypeURL string `json:"type_url"`
					Value   string `json:"value"`
				} `json:"data"`
				Rectangle struct {
					Vertices []struct {
						X int32 `json:"x"`
						Y int32 `json:"y"`
					} `json:"vertices"`
				} `json:"rectangle"`
			} `json:"algo"`
			Diagnosis struct {
				TypeFrame         string `json:"type_frame"`
				DiagnoseSummaries []struct {
					DiagnosisItem string  `json:"diagnosis_item"`
					TypeItem      string  `json:"type_item"`
					Score         float32 `json:"score"`
				} `json:"diagnose_summaries"`
			} `json:"diagnosis"`
		} `json:"face_info"`
		Feature struct {
			Type    string `json:"type"`
			Version int32  `json:"version"`
			Blob    string `json:"blob"`
		} `json:"feature"`
		ImagesOrientation string `json:"images_orientation"`
	} `json:"responses"`
}

const Float32Bytes = 4

func Base64Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, err
	} else {
		return dst[:n], nil
	}
}

func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func EncodeFloat32(src []float32) []byte {
	dst := make([]byte, len(src)*Float32Bytes)
	offset := 0
	for _, f := range src {
		binary.LittleEndian.PutUint32(dst[offset:], math.Float32bits(f))
		offset += Float32Bytes
	}
	return dst
}

func DecodeFloat32(src []byte) []float32 {
	dst := make([]float32, len(src)/Float32Bytes)
	offset := 0
	for i := range dst {
		b := binary.LittleEndian.Uint32(src[offset:])
		offset += Float32Bytes
		dst[i] = math.Float32frombits(b)
	}
	return dst
}
func ReadLines(r io.Reader, delim byte) chan []byte {
	br := bufio.NewReader(r)
	ch := make(chan []byte, 8)

	go func() {
		defer close(ch)
		for {
			b, err := br.ReadBytes(delim)
			if err != nil {
				if err == io.EOF {
					return
				} else {
					fmt.Print("failed to read: %s", err)
				}
			}
			//log.Debugf("read %d bytes from file", len(b))
			ch <- b
		}
	}()
	return ch
}

// fp, err := os.Open(file)
func ChangeTobase(path string) string {
	srcByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Print(err)
	}
	res := base64.StdEncoding.EncodeToString(srcByte)
	return res
}
