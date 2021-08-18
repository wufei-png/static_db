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
package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	lw "static_db_wf/lib_worker"

	"strconv"
	"strings"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func handle_req(w http.ResponseWriter, r *http.Request) {
	var req lw.Request
	p, _ := ioutil.ReadAll(r.Body)
	//wufei := []string{string(p)}
	res := strings.TrimSpace(string(p))
	//把字符串以空格分割成字符串数组
	str_arr := strings.Split(res, " ")
	var int_arr []float64
	for _, i := range str_arr {
		j, err := strconv.ParseFloat(i, 64)
		if err != nil {
			panic(err)
		}
		int_arr = append(int_arr, j)
	}
	var length = len(int_arr)

	if int_arr[length-1] == 3 { //删除
		req.Request_type = 3
		req.I_delete = int_arr[0]
		req.J_delete = int_arr[1]
		if length != 3 {
			fmt.Fprintln(w, "删除参数长度错误") //比上面那个多一个topk
			return
		}
		fmt.Fprintln(w, lw.Delete(req))
		return
	} else if int_arr[length-1] == 4 { //添加
		if (length-1)%4 != 0 {
			fmt.Fprintln(w, "添加参数长度错误") //比上面那个多一个topk
			return
		}
		req.Request_type = 4
		for i := 0; i < ((len(int_arr) - 1) / 4); i++ {
			req.Xlzu.Xl = append(req.Xlzu.Xl, int_arr[4*i:4*i+4]) //可以设置不同维度的
		}
		fmt.Fprintln(w, lw.Add(req))
		return
	} else if int_arr[length-1] == 1 || int_arr[length-1] == 2 {
		if (length-2)%4 != 0 {
			fmt.Fprintln(w, "参数长度错误") //比上面那个多一个topk
			return
		}
		for i := 0; i < ((len(int_arr) - 2) / 4); i++ {
			req.Xlzu.Xl = append(req.Xlzu.Xl, int_arr[4*i:4*i+4]) //可以设置不同维度的
		}
		req.Topk = int(int_arr[(len(int_arr) - 2)])
		req.Request_type = int(int_arr[(length - 1)])
		var result lw.Result = lw.Search(req)
		fmt.Fprintln(w, result)
		return
	}

	//fmt.Println(req.Xlzu.Xl)

}

// func index(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, index2())
// }
func Server_init() {
	// Set up a connection to the server.
	lw.DB = lw.Lib_worker_DBinit(100, 5, 4)
	//fmt.Print(lw.DB.Head.Data)
	lw.Dbinit_train()

	http.HandleFunc("/search", handle_req)
	// 启动web服务，监听9090端口
	http.ListenAndServe(":9090", nil)
}
