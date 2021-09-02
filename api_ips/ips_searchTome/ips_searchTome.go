package main

import (
	"encoding/json"
	"fmt"

	//lib "static_db_wf/api_ips/lib_init"
	tools "static_db_wf/api_ips/tools"
)

func main() {

	//lib.Lib_init()

	var path string
	var num string
	url1 := "http://172.20.25.180:30080/engine/image-process/face_25000/v1/batch_detect_and_extract"
	url2 := "http://localhost:9090/search"
	var req req_wf_search
	for i := 7; i < 8; i++ {
		num = fmt.Sprint(i)
		path = `/Users/wufei1/Desktop/face/test` + num + `.jpeg`
		body := tools.Post_req(url1, tools.StructTojson_ips(path))
		//fmt.Print(string(body))
		var params tools.Res_ips
		//fmt.Print(1)
		err := json.Unmarshal(body, &params)
		//fmt.Print(2)
		if err != nil {
			fmt.Print(err)
			return
		}
		//fmt.Print(3)
		var blob = params.Responses[0].Feature.Blob
		feature, err := tools.GetRawFeatureFromBase64(blob)
		//fmt.Print(3)
		if err != nil {
			fmt.Print(err)
			return
		}
		//fmt.Print(4)
		b := tools.EncodeFloat32(feature.Raw)
		b = tools.Base64Encode(b)
		c := Feature{Blob: string(b)}
		req.Features = append(req.Features, c)
		//fmt.Println(string(reqjson))
		//jsonStu是[]byte类型，转化成string类型便于查看
	}
	req.SearchType = 1
	req.Topk = 3
	reqjson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return
	}
	body := tools.Post_req(url2, string(reqjson))
	//fmt.Print(string(body))
	var params res_wf_search
	//fmt.Print(3)
	// fmt.Print(string(body))
	err = json.Unmarshal(body, &params)
	//fmt.Print(4)
	if err != nil {
		fmt.Print(err)
		return
	}
	//lw.DB.NumOfAllNodes()

	fmt.Print(params.Distancetopk.Vectors)
	//fmt.Print(lib.Delete())

	//fmt.Print(lw.DB.NumOfAllNodes())
}

type req_wf_search struct {
	SearchType int32     `json:"search_type"`
	Features   []Feature `json:"features"`
	Topk       int32     `json:"topk"`
}
type res_wf_search struct {
	Distancetopk  Vectors    `json:"distancetopk"`
	Featuresgroup []Features `json:"featuresgroup"`
}
type Feature struct {
	Blob string `json:"blob"`
}
type Features struct {
	Features []Feature `json:"Features`
}
type Vectors struct {
	Vectors []Vector `json:vectors`
}
type Vector struct {
	Vector []float32 `json:vector`
}
