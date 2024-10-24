package main

import "runtime"

func main() {
	println("GOOS :", runtime.GOOS)
	println("GOARCH :", runtime.GOARCH)
}
