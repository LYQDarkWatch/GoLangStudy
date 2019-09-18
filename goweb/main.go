package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtCustomClaims struct {
	jwt.StandardClaims

	// 追加自己需要的信息
	Uid   uint `json:"uid"`
	Admin bool `json:"admin"`
}

type Token struct {
	Token string `json:"token"`
}
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

func CreateToken(SecretKey []byte, issuer string, Uid uint, isAdmin bool) (tokenString string, err error) {
	claims := &jwtCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: int64(time.Now().Add(time.Hour * 72).Unix()),
			Issuer:    issuer,
		},
		Uid,
		isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

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
func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	body, _ := http.Get(url)
	a, _ := ioutil.ReadAll(body.Body)
	b := string(a)
	github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"] = b
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	if r.Method == "POST" {
		var user Auth
		s, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal([]byte(s), &user)
		if user.Username == "admin" && user.Pwd == "123456" {
			keyInfo := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()"
			var dataStr = string(s)
			data := jwt.StandardClaims{Subject: dataStr, ExpiresAt: time.Now().Unix() - 1000}
			tokenInfo := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
			tokenStr, _ := tokenInfo.SignedString([]byte(keyInfo))
			fmt.Println("myToken is: ", tokenStr)
			tokenInfo, _ = jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
				return keyInfo, nil
			})

			fmt.Fprintf(w, `{"code":200}`)
			fmt.Println("欢迎用户: " + user.Username + " 登录成功")
		} else {
			fmt.Fprintf(w, `{"code":401}`)
			fmt.Println("账号或密码错误")
		}

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
