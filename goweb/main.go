package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type MyMux struct {
}
type Auth struct {
	Username string `json:"user"`
	Pwd      string `json:"password"`
}
type Resp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

var github = map[string]string{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/about" {
		about(w, r)
		return
	}
	if r.URL.Path == "/clear" {
		ClearMap(w, r)
		return
	}
	if r.URL.Path == "/login" {
		Login(w, r)
		return
	}
}
func about(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	body, _ := http.Get(url)
	a, _ := ioutil.ReadAll(body.Body)
	b := string(a)
	github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"] = b
	fmt.Println(github)
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
	return b
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	if r.Method == "POST" {
		s, _ := ioutil.ReadAll(r.Body)
		a := string(s)
		login(s)
		fmt.Print(a)
	}
}
func login(userInfo []byte) {
	var user Auth
	var result Resp
	json.Unmarshal(userInfo, &user)
	if user.Username == "admin" && user.Pwd == "123456" {
		result.Code = "200"
		result.Msg = "登录成功"
		fmt.Print(result.Msg)
	} else {
		result.Code = "401"
		result.Msg = "账户名或密码错误"
		fmt.Print(result.Msg)

	}

}
func Start() {
	mux := &MyMux{}
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func ClearMap(w http.ResponseWriter, r *http.Request) {

	delete(github, "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc")
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}
func main() {
	Start()
}
