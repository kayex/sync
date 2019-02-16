package main

import (
	"github.com/kayex/sync"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stdout, "", 0)

	if len(os.Args) < 4 {
		l.Fatal("usage: sync target destination source")
	}

	host := sync.ParseTargetURI("tcp://" + os.Args[1])
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

