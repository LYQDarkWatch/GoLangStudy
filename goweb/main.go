package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	// var u *User = new(User)
	var user UserCredentials
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
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
	body, _ := http.Get(url)
	a, _ := ioutil.ReadAll(body.Body)
	b := string(a)
	github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"] = b
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}
func ClearMap(w http.ResponseWriter, r *http.Request) {
	delete(github, "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc")
	fmt.Fprint(w, github["https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"])
}

const (
	RedisURL            = "redis://*****:6379"
	redisMaxIdle        = 3   //最大空闲连接数
	redisIdleTimeoutSec = 240 //最大空闲连接时间
	RedisPassword       = "*****"
)

// NewRedisPool 返回redis连接池
func NewRedisPool(redisURL string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     redisMaxIdle,
		IdleTimeout: redisIdleTimeoutSec * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
				return nil, fmt.Errorf("redis auth password error: %s", authErr)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}
