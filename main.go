package main

import (
	"flag"
	"log"
	"os"

	_ "net/http/pprof"

	"github.com/muziling/freedns-go/freedns"
	"github.com/muziling/freedns-go/chinaip"
)

func main() {
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	*/

	var (
		fastDNS  string
		cleanDNS string
		listen   string
		chnRouteFile string
	)

	flag.StringVar(&fastDNS, "f", "114.114.114.114:53", "The fast/local DNS upstream.")
	flag.StringVar(&cleanDNS, "c", "8.8.8.8:53", "The clean/remote DNS upstream.")
	flag.StringVar(&listen, "l", "0.0.0.0:53", "Listening address.")
	flag.StringVar(&chnRouteFile, "r", "/root/chnroute.txt", "China ip list file.")

	flag.Parse()

	b, err := chinaip.LoadChnRoute(chnRouteFile)
	if !b {
		log.Fatalln(err)
		os.Exit(-1)
	}

	s, err := freedns.NewServer(freedns.Config{
		FastDNS:   fastDNS,
		CleanDNS:  cleanDNS,
		Listen:    listen,
		CacheSize: 1024 * 5,
	})
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}

	log.Fatalln(s.Run())
	os.Exit(-1)
}
