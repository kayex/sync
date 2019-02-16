package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	l := log.New(os.Stderr, "", log.LstdFlags)

	host := Host{}

	flag.StringVar(&host.Location, "host", "localhost", "target host")
	flag.StringVar(&host.User, "user", "root", "user")
	flag.StringVar(&host.Password, "password", "", "password")
	flag.IntVar(&host.Port, "port", 22, "port")
	flag.Parse()

	target := NewSFTPTarget(host)
	l.Printf("dialing tcp:%v:%v@%v:%d", host.User, host.Password, host.Location, host.Port)
	err := target.Connect()
	if err != nil {
		l.Fatal(err)
	}
	l.Printf("connected to %v on port %d", host.Location, host.Port)
	err = target.Close()
	if err != nil {
		l.Fatal(err)
	}
	l.Print("connection closed.")
}
