package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type MyMux struct {
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
