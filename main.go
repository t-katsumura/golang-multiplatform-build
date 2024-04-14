package main

import "runtime"

func main() {
	println("GOOS :", runtime.GOOS)
	println("GOARCG :", runtime.GOARCH)
}
