package main

// type Data struct {
// 	NAME string `json:"name"`
// }
// type Employee struct {
// 	DATA []Data `json:"items"`
// }

// func main() {
// 	url := "https://api.github.com/search/repositories?q=stars:%3E=500&sort=stars&order=desc"
// 	resp, _ := http.Get(url) // _下划线表示申明一个只写变量
// 	s := Employee{}
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	json.Unmarshal([]byte(body), &s)
// 	fmt.Println(fmt.Sprintf("%+v", s))
// }

// func main() {
// 	var a = [...]int{1, 2, 3}
// 	var b = &a
// 	fmt.Println(a[0], a[1])

// 	fmt.Println(b[0], b[1])

// 	for i, v := range b {
// 		fmt.Println(i, v)
// 	}
// }
