package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andrei-m/adstxt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./adstxtcmd <host>")
	}
	c := &http.Client{
		Timeout:       1 * time.Minute,
		CheckRedirect: adstxt.CheckRedirect,
	}
	adsTxt, err := adstxt.Resolve(c, os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(adsTxt)
}
