package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/redis.v4"
)

const (
	SecretKey = "I have login"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type User struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {

	http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/about", about)
	http.HandleFunc("/clear", ClearMap)
	log.Println("Now listening...")
	http.ListenAndServe(":8080", nil)
}

func main() {
	StartServer()
}

var user UserCredentials

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	// var u *User = new(User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if user.Username == "admin" && user.Password == "123456" {
		fmt.Println("欢迎用户: " + user.Username + " 登录成功")
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)

		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		// claims["id"] = u.ID
		claims["userName"] = user.Username
		claims["password"] = user.Password
		token.Claims = claims
		fmt.Println(claims)
		tokenString, _ := token.SignedString([]byte(SecretKey))

		set(user.Username, tokenString, 6000)
		fmt.Println(tokenString)
		response := Token{tokenString}
		JsonResponse(response, w)

	} else {
		fmt.Fprintf(w, `{"code":401}`)
		fmt.Println("账号或密码错误")

	}

}
func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

var github = map[string]string{}

func about(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "*,Authorization")
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Request-Headers", "authorization")
	autoken := r.Header.Get("Authorization")
	fmt.Println(autoken)
	fmt.Println(user.Username)
	if autoken == get(user.Username) {
		url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
		body, _ := http.Get(url)
		a, _ := ioutil.ReadAll(body.Body)
		b := string(a)
		github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"] = b
		fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
	} else {
		fmt.Fprint(w, `{"code":400}`)
	}

	// url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	// body, _ := http.Get(url)
	// a, _ := ioutil.ReadAll(body.Body)
	// b := string(a)
	// github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"] = b
	// fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}
func ClearMap(w http.ResponseWriter, r *http.Request) {
	delete(github, "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc")
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:         "127.0.0.1:6379",
		PoolSize:     1000,
		ReadTimeout:  time.Millisecond * time.Duration(100),
		WriteTimeout: time.Millisecond * time.Duration(100),
		IdleTimeout:  time.Second * time.Duration(10),
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic("init redis error")
	} else {
		fmt.Println("init redis ok")
	}
}

// func get(key string) (string, bool) {
// 	// r, err := Client.Get(key).Result()
// 	// if err != nil {
// 	// 	return "", false
// 	// }
// 	// return r, true
// }
func get(key string) string {
	r, _ := Client.Get(key).Result()
	return r

}
func set(key string, val string, expTime int32) {
	Client.Set(key, val, time.Duration(expTime)*time.Second)
}
