package main

import "fmt"

func main() {
	dao := InitDao("localhost", 6379)

	err := dao.Set("hello", "123")
	if err != nil {
		fmt.Println(dao.Get("hello"))
	}
}
