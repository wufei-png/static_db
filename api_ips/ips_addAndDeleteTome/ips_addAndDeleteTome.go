package main

import (
	"encoding/json"
	"fmt"
	tools "static_db_wf/api_ips/tools"
)

func main() {
	var path string
	var num string
	url1 := "http://172.20.25.180:30080/engine/image-process/face_25000/v1/batch_detect_and_extract"
	url2 := "http://localhost:9090/add"
	req := Req_wf_add{}
	for i := 1; i < 7; i++ {
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
	reqjson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return
	}
	body := tools.Post_req(url2, string(reqjson))
	var params Res_wf_add
	//fmt.Print(3)
	// fmt.Print(string(body))
	err = json.Unmarshal(body, &params)
	//fmt.Print(4)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(params.Status)
}
func Delete() string {
	url := "http://localhost:9090/delete"
	var req Req_wf_delete
	req.Row = 2
	req.Col = 2
	reqjson, err := json.Marshal(req)
	if err != nil {
		fmt.Println("生成json字符串错误")
		return ""
	}
	body := tools.Post_req(url, string(reqjson))
	var params Res_wf_delete
	//fmt.Print(3)
	// fmt.Print(string(body))
	err = json.Unmarshal(body, &params)
	//fmt.Print(4)
	if err != nil {
		fmt.Print("unmarshal错误")
		return ""
	}
	return params.Status
	//fmt.Print(lw.DB.NumOfAllNodes())
}

// func test(){
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// index1 := gojsonq.New().FromString(string(body)).From("responses").First()
// interface2String(index1)
// fmt.Print(index1)
// index2 := gojsonq.New().FromString(index1).Find("feature")
// fmt.Print(index2)
// index3 := index2.Get().(map[string]interface{})
// feature, ok := params.(map[string]interface{}) //.(map[string]interface{})["feature"].(map[string]interface{})["blob"]
// if !ok {
// 	fmt.Println("1")
// 	return
// }
// feature1, ok := feature["responses"].(map[string]interface{})
// if !ok {
// 	fmt.Println("2")
// 	return
// }
// fmt.Print(feature1)
// strings.Split(index, ":")
// decoder := json.NewDecoder(res.Body)
// 用于存放参数key=value数据
// }

type Feature struct {
	Blob string `json:"blob"`
}
type Req_wf_add struct {
	Features []Feature `json:"features"`
}
type Res_wf_add struct {
	Status string `json:"status"`
}
type Req_wf_delete struct {
	Row int32 `json:"row"`
	Col int32 `json:"col"`
}
type Res_wf_delete struct {
	Status string `json:"status"`
}
