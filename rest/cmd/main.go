package main

import (
	"fmt"
	"github.com/hosseintrz/torob/rest/conf"
	"github.com/hosseintrz/torob/rest/internal/server"
)

func main() {
	authConf := &conf.Config{
		Host: "localhost",
		Port: "9090",
	}
	restConf := &conf.Config{
		Host: "localhost",
		Port: "8080",
	}
	prodConf := &conf.Config{
		Host: "localhost",
		Port: "8282",
	}
	suppConf := &conf.Config{
		Host: "localhost",
		Port: "9191",
	}

	s := server.New(restConf, authConf, prodConf, suppConf)
	errChan := s.Serve()
	select {
	case err := <-errChan:
		fmt.Printf("error : %s ", err.Error())
	}
}
