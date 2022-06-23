package main

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/internal/gateway"
	server2 "github.com/hosseintrz/torob/rest/internal/server"
)

func main() {
	fmt.Println("here1")
	authConf := &server2.Config{
		Host: "localhost",
		Port: "9090",
	}
	restConf := &server2.Config{
		Host: "localhost",
		Port: "8080",
	}
	authClient := gateway.InitAuthClient(fmt.Sprintf("%s:%s", authConf.Host, authConf.Port))
	s := server2.New(restConf, authClient)
	errChan := s.Serve()
	fmt.Println("here")
	select {
	case err := <-errChan:
		fmt.Printf("error : %s ", err.Error())
	}
}
