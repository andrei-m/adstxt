package main

import (
	"log"
	"os"

	"github.com/andrei-m/adstxt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./adstxtcmd <host>")
	}
	adsTxt, err := adstxt.DefaultResolve(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(adsTxt)
}
