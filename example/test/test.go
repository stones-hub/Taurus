package main

import (
	"fmt"
)

func demo(name string) func() {
	return func() {
		fmt.Println("demo name:", name)
		fmt.Println("do something")
		fmt.Println("demo name:", name)
	}
}

func demo2(name string) func(string) string {
	return func(do string) string {
		fmt.Println("demo2 name:", name)
		fmt.Println("do something:", do)
		fmt.Println("demo2 name:", name)
		return "demo2 name:" + name + " do something:" + do
	}
}

func createMid(name string) func(func()) func() {
	return func(next func()) func() {
		return func() {
			fmt.Println("createMid 开始 name1:", name)
			next()
			fmt.Println("createMid 结束 name2:", name)
		}
	}
}

func main() {

	d := demo("test")
	d()

	fmt.Println("--------------------------------")

	d2 := demo2("test")

	fmt.Println(d2("do worker"))

	fmt.Println("--------------------------------")

	handle := func() {
		fmt.Println("我是一个handler")
	}

	mid := createMid("日志")
	warppedHandle := mid(handle)
	warppedHandle()
}
