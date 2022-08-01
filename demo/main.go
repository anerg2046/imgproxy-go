package main

import (
	"fmt"
	imgproxygo "imgproxy-go"
)

func main() {
	imgproxy := imgproxygo.N(imgproxygo.Config{
		BaseUrl:       "http://localhost:8080",
		Key:           "736563726574",
		Salt:          "68656C6C6F",
		SignatureSize: 8,
		Encode:        true,
	})
	s := imgproxy.Builder().Width(300).Height(400).Quality(100).Gen("https://m.media-amazon.com/images/M/MV5BMmQ3ZmY4NzYtY2VmYi00ZDRmLTgyODAtZWYzZjhlNzk1NzU2XkEyXkFqcGdeQXVyNTc3MjUzNTI@.jpg")
	fmt.Println(s)
}
