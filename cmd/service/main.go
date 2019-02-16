package main

import (
	"fmt"
	"github.com/kayex/sync"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "", log.LstdFlags)

	//g := sync.NewGit("https://github.com/kayex/driverstats-client.git")
	//err := g.Clone("src")
	//if err != nil {
	//	l.Fatal(err)
	//}

	host := sync.ParseTargetURI("tcp://" + req("TARGET"))
	target := sync.NewSFTPTarget(host)
	l.Printf("dialing tcp://%v:%v", host.Location, host.Port)
	err := target.Connect()
	if err != nil {
		l.Fatal(err)
	}

	s := sync.NewSync(l, target)

	err = s.Sync("/www", "src/dist/driverstats-client")
	if err != nil {
		l.Fatal(err)
	}
}

func req(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Errorf("missing required environment variable %v", key))
	}
	return v
}
