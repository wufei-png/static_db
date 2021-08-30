package test_db

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func test() {

	url := "http://localhost:9090/search"
	method := "POST"

	payload := strings.NewReader(`1 2 3 4 3 1`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
