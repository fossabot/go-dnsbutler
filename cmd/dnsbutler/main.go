package main

import (
	"flag"

	"github.com/stahlstift/go-dnsbutler/pkg/dnsbutler"
)

var configPath = "dnsbutler.json"

func init() {
	flag.StringVar(&configPath, "-c", configPath, "config file path")
	flag.Parse()
}

func main() {
	dnsbutler.Start(configPath)
}
