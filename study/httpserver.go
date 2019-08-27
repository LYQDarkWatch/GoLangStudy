package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	NAME string `json:"name"`
}

type Employee struct {
	DATA []Data `json:"items"`
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm() //解析参数，默认是不会解析的
	// fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }
	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	resp, _ := http.Get(url) // _下划线表示申明一个只写变量
	s := Employee{}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal([]byte(body), &s)
	fmt.Println(fmt.Sprintf("%+v", s))
	fmt.Fprintf(w, "Hello Wrold!,LYQQQQQQ") //这个写入到w的是输出到客户端的
}
func main() {
	http.HandleFunc("/", sayhelloName)                //设置访问的路由
	err := http.ListenAndServe("127.0.0.1:9090", nil) //设置要监听的Ip和端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
