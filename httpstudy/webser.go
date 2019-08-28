package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MyMux struct {
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		sayHelloName(w, r)
		return
	}
	if r.URL.Path == "/about" {
		about(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path: ", r.URL.Path)
	fmt.Println("scheme: ", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key: ", k)
		fmt.Println("val: ", strings.Join(v, " "))
	}
	fmt.Fprintf(w, "hello chain!")
}

type Data struct {
	NAME string `json:"name"`
}
type Employee struct {
	DATA []Data `json:"items"`
}

func about(w http.ResponseWriter, r *http.Request) *Employee {

	// fmt.Fprintf(w, "i am chain, from shanghai")
	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	body, _ := http.Get(url)
	// a := Employee{}
	// b := new(Employee)
	a := new(Employee)
	s, _ := ioutil.ReadAll(body.Body)
	json.Unmarshal([]byte(s), &a)
	fmt.Println(fmt.Sprintf("%+v", a))
	fmt.Fprintf(w, fmt.Sprintf("%v", a))
	return a
}

func Start() {
	mux := &MyMux{}
	err := http.ListenAndServe(":9090", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main() {
	Start()

}
