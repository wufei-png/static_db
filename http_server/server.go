/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package http_server

import (
	"encoding/json"
	"fmt"
	"net/http"
	lw "static_db_wf/lib_worker"

	"strconv"
	"strings"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

// type;  add delete
// search -》search_type 01
// vectors:1 1 1 1
// topk 1
// index : i,j
func strToInt(str string) []float32 {
	str_arr := strings.Split(str, " ")
	var int_arr []float32
	for _, i := range str_arr {
		j, err := (strconv.ParseFloat(i, 32))
		if err != nil {
			panic(err)
		}
		int_arr = append(int_arr, float32(j))
	}
	return int_arr
}

// func handle_req(w http.ResponseWriter, r *http.Request) {
// 	var req lw.Request
// 	var params map[string]string = handlePostJson(w, r)
// 	if params["type"] == "" {
// 		fmt.Fprint(w, "未给type")
// 		return
// 	}
// 	if params["type"] == "search" {
// 		var int_arr = strToInt(params["vectors"])
// 		length := len(int_arr)
// 		if length == 0 || length%4 != 0 {
// 			fmt.Fprint(w, "向量长度不对")
// 			return
// 		}
// 		for i := 0; i < (length / 4); i++ {
// 			req.Vectors.Vector = append(req.Vectors.Vector, int_arr[4*i:4*i+4]) //可以设置不同维度的
// 		}

// 		var topk = strToInt(params["topk"])
// 		if len(topk) != 1 {
// 			fmt.Fprint(w, "topk长度不对")
// 			return
// 		}
// 		req.Topk = int(topk[0])
// 		var request_type = strToInt(params["search_type"])
// 		if len(request_type) != 1 {
// 			fmt.Fprint(w, "request_type长度不对")
// 			return
// 		}
// 		req.Request_type = int(request_type[0])
// 		var result lw.Result = lw.Search(req)
// 		fmt.Fprintln(w, result)
// 		return
// 	} else if params["type"] == "add" {
// 		var int_arr = strToInt(params["vectors"])
// 		length := len(int_arr)

// 		if length%4 != 0 {
// 			fmt.Fprint(w, "向量长度不对")
// 			return
// 		}
// 		for i := 0; i < (length / 4); i++ {
// 			req.Vectors.Vector = append(req.Vectors.Vector, int_arr[4*i:4*i+4]) //可以设置不同维度的
// 		}
// 		fmt.Fprintln(w, lw.Add(req))
// 		return

// 	} else if params["type"] == "delete" {
// 		var int_arr = strToInt(params["index"])
// 		if len(int_arr) != 2 {
// 			fmt.Fprint(w, "index长度不对")
// 			return
// 		}
// 		req.I_delete = int_arr[0]
// 		req.J_delete = int_arr[1]
// 		fmt.Fprintln(w, lw.Delete(req))
// 		return
// 	} else {
// 		fmt.Fprint(w, "type类型不对")
// 	}
// }
func handle_add(w http.ResponseWriter, r *http.Request) {
	var req lw.Request
	var params map[string]string = handlePostJson(w, r)
	var int_arr = strToInt(params["vectors"])
	length := len(int_arr)
	if length%4 != 0 {
		fmt.Fprint(w, "向量长度不对")
		return
	}
	for i := 0; i < (length / 4); i++ {
		req.Vectors.Vector = append(req.Vectors.Vector, int_arr[4*i:4*i+4]) //可以设置不同维度的
	}
	fmt.Fprintln(w, lw.Add(req))
	return
}
func handle_delete(w http.ResponseWriter, r *http.Request) {
	var req lw.Request
	var params map[string]string = handlePostJson(w, r)
	var int_arr = strToInt(params["index"])
	if len(int_arr) != 2 {
		fmt.Fprint(w, "index长度不对")
		return
	}
	req.I_delete = int_arr[0]
	req.J_delete = int_arr[1]
	fmt.Fprintln(w, lw.Delete(req))
	return
}
func handle_search(w http.ResponseWriter, r *http.Request) {
	var req lw.Request
	var params map[string]string = handlePostJson(w, r)
	var int_arr = strToInt(params["vectors"])
	length := len(int_arr)
	if length == 0 || length%4 != 0 {
		fmt.Fprint(w, "向量长度不对")
		return
	}
	for i := 0; i < (length / 4); i++ {
		req.Vectors.Vector = append(req.Vectors.Vector, int_arr[4*i:4*i+4]) //可以设置不同维度的
	}

	var topk = strToInt(params["topk"])
	if len(topk) != 1 {
		fmt.Fprint(w, "topk长度不对")
		return
	}
	req.Topk = int(topk[0])
	var request_type = strToInt(params["search_type"])
	if len(request_type) != 1 {
		fmt.Fprint(w, "request_type长度不对")
		return
	}
	req.Request_type = int(request_type[0])
	var result lw.Result = lw.Search(req)
	fmt.Fprintln(w, result)
	return

}
func handlePostJson(writer http.ResponseWriter, request *http.Request) map[string]string {
	// 根据请求body创建一个json解析器实例
	decoder := json.NewDecoder(request.Body)

	// 用于存放参数key=value数据
	var params map[string]string

	// 解析参数 存入map
	decoder.Decode(&params)

	return params
}

// func index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, index2())
// }
// func Server_init() {
// 	// Set up a connection to the server.
// 	//lw.DB = lw.Lib_worker_DBinit(10000, 100, 4)

// 	//lw.Display_allfirst()
// 	//lw.Display_DBbydistence()
// 	//lw.Dbinit_train()
// 	server_gw()

// 	server.Server_grpc()
// }
