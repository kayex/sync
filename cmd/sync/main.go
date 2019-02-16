package main

import (
	"fmt"
	"github.com/kayex/sync"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	if len(os.Args) < 4 {
		l.Fatal("usage: sync target destination source")
	}

	host := parseTargetURI("tcp://" + os.Args[1])
	dst := os.Args[2]
	src := os.Args[3]

	target := sync.NewSFTPTarget(host)
	l.Printf("dialing tcp://%v:%v", host.Location, host.Port)
	err := target.Connect()
	if err != nil {
		l.Fatal(err)
	}

	s := sync.NewSync(l, target)
	err = s.Sync(dst, src)
	if err != nil {
		l.Fatal(err)
	}
}

func parseTargetURI(uri string) sync.Host {
	u, err := url.Parse(uri)
	if err != nil {
		panic(fmt.Errorf("failed parsing target URI: %v", err))
	}

	host, p, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(fmt.Errorf("failed parsing target URI: %v", err))
	}

	port := 0
	if p != "" {
		port, err = strconv.Atoi(p)
		if err != nil {
			panic(fmt.Errorf("failed parsing target port: %v", err))
		}
	}

	user := u.User.Username()
	password, _ := u.User.Password()

	return sync.Host{
		Location: host,
		User: user,
		Password: password,
		Port: port,
	}
}
