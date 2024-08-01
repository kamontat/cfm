package main

import (
	"fmt"
	"log"
	"net/url"
)

func BuildCloudflareCDN(hostname string) (target *url.URL) {
	var cdn = fmt.Sprintf("https://%s.cdn.cloudflare.net", hostname)
	var err error
	target, err = url.Parse(cdn)
	if err != nil {
		log.Fatalf("Invalid cloudflare CDN URL: %v", err)
	}
	return
}
