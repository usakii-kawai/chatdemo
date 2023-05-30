package main

import "chatdemo/server"

const version = "v1"

func main() {
	opts := &server.ServerStartOptions{
		Id:     "abc",
		Listen: ":8080",
	}
	server.RunServerStart(opts, version)
}
